package internal

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/itchyny/gojq"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type RunnerOpts struct {
	Logger       zerolog.Logger
	HttpExecutor Executor
	GrpcExecutor Executor
	Parser       Parser
	Output       io.Writer
	FailFast     bool
}

func NewRunner(opts RunnerOpts) *Runner {
	return &Runner{
		log:          opts.Logger,
		httpExecutor: opts.HttpExecutor,
		grpcExecutor: opts.GrpcExecutor,
		parser:       opts.Parser,
		ctxVariables: make(map[string]any),
		output:       opts.Output,
		failFast:     opts.FailFast,
	}
}

type Runner struct {
	log          zerolog.Logger
	httpExecutor Executor
	grpcExecutor Executor
	parser       Parser
	ctxVariables map[string]any
	output       io.Writer
	failFast     bool
}

func (r *Runner) Run(path string) error {
	sequences, err := r.parser.Parse(path)
	if err != nil {
		return fmt.Errorf("error parsing: %w", err)
	}

	return r.runSequences(sequences)
}

func (r *Runner) runSequences(seqs map[string]Sequence) error {
	var errs []error

	for name, seq := range seqs {
		r.log.Info().Str("sequence", name).Msg("executing sequence")
		if err := r.runSingleSequence(seq); err != nil {
			r.log.Err(err).Msg("encountered error during execution")
			errs = append(errs, fmt.Errorf("error during sequence %v: %w", name, err))
			if r.failFast {
				break
			}
		}
		// Reset variables at the end of a sequence so we don't bleed over
		r.ctxVariables = map[string]any{}
	}

	return errors.Join(errs...)
}

func (r *Runner) runSingleSequence(seq Sequence) error {
	// Set any predefined global vars
	if seq.Vars != nil {
		for k, v := range seq.Vars {
			r.ctxVariables[k] = v
		}
	}
	if err := r.resolveImports(&seq); err != nil {
		return fmt.Errorf("error resolving imports: %w", err)
	}

	for idx, c := range seq.Calls {
		call := c
		if c.FromImport != nil {
			impSeq, ok := seq.importedCalls[c.FromImport.Name]
			if !ok {
				return fmt.Errorf("unable to find import %v", c.FromImport.Name)
			}
			impCall, ok := impSeq[c.FromImport.Call]
			if !ok {
				return fmt.Errorf("unable to find call %v in imported sequence %v", c.FromImport.Call, c.FromImport.Name)
			}
			call = impCall
		}
		name := call.Name
		if name == "" {
			name = fmt.Sprintf("call_%v", idx)
		}
		r.log.Info().Str("call", name).Msg("executing call")
		call, err := r.evaluateTemplate(call, seq.path)
		if err != nil {
			return err
		}

		exec, err := r.getClient(call.Type)
		if err != nil {
			return fmt.Errorf("error creating request client: %w", err)
		}

		result, err := exec.Execute(call)
		if err != nil {
			return fmt.Errorf("error executing call %v: %w", name, err)
		}

		wantStatus := call.WantStatus
		if wantStatus == 0 && call.GetType() == RequestTypeHttp {
			wantStatus = http.StatusOK
		}

		if result.StatusCode != wantStatus {
			r.log.Error().Interface("body", result.Body).Msg("body")
			r.log.Err(result.Error).Msg("error msg")
			return fmt.Errorf("got incorrect status: want (%v) got (%v)", wantStatus, result.StatusCode)
		}

		if call.Print {
			bodyBytes, err := json.MarshalIndent(result.Body, "", "   ")
			if err != nil {
				r.log.Err(err).Msg("error marshalling body for output")
				return err
			}
			fmt.Fprint(r.output, string(bodyBytes))
		}

		for _, exp := range call.Exports {
			value, err := r.executeJQString(result.Body, exp.JQ)
			if err != nil {
				return err
			}
			r.ctxVariables[exp.As] = value
		}

		for _, ass := range call.Asserts {
			value, err := r.executeJQ(result.Body, ass.JQ)
			if err != nil {
				return err
			}
			if diff := cmp.Diff(ass.Expected, value); diff != "" {
				r.log.Error().Msg("failed assertion")
				fmt.Println(diff)
				return fmt.Errorf("failed assert")
			}
		}
	}

	return nil
}

//nolint:ireturn
func (r *Runner) getClient(typ RequestType) (Executor, error) {
	switch typ {
	case RequestTypeHttp:
		return r.httpExecutor, nil
	case "":
		return r.httpExecutor, nil
	case RequestTypeGrpc:
		return r.grpcExecutor, nil
	default:
		return nil, fmt.Errorf("unhandled client type of '%v'", typ)
	}
}

func (r *Runner) resolvePath(seqPath string, path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(seqPath, path)
	}
}

func (r *Runner) readFile(seqPath string, path string) ([]byte, error) {
	fileBytes, err := os.ReadFile(r.resolvePath(seqPath, path))
	if err != nil {
		return nil, fmt.Errorf("error reading %v: %w", path, err)
	}
	return fileBytes, nil
}

func (r *Runner) genFuncs(seqPath string) template.FuncMap {
	return template.FuncMap{
		"env": func(key string) (string, error) {
			val, ok := os.LookupEnv(key)
			if !ok {
				return "", fmt.Errorf("env var %v not set", key)
			}
			return val, nil
		},
		"readfileb64": func(path string) (string, error) {
			fBytes, err := r.readFile(seqPath, path)
			if err != nil {
				return "", err
			}
			return base64.StdEncoding.EncodeToString(fBytes), nil
		},
		"readfile": func(path string) (string, error) {
			fBytes, err := r.readFile(seqPath, path)
			if err != nil {
				return "", err
			}
			return string(fBytes), nil
		},
	}
}

func (r *Runner) resolveImports(seq *Sequence) error {
	r.log.Debug().Msg("resolving imports")
	for name, path := range seq.Imports {
		r.log.Debug().Str("import", name).Msg("parsing imported sequence")
		impSeq, err := r.parser.ParseSingleSequence(r.resolvePath(seq.path, path))
		if err != nil {
			return fmt.Errorf("error parsing imported sequence '%v': %w", path, err)
		}

		if impSeq.Imports != nil {
			return fmt.Errorf("import %v contains imports. Only one level of nesting supported", name)
		}

		seq.importedCalls[name] = make(map[string]Call)
		for _, call := range impSeq.Calls {
			seq.importedCalls[name][call.Name] = call
		}
	}
	return nil
}

func (r *Runner) evaluateTemplate(call Call, seqPath string) (Call, error) {
	callBytes, err := yaml.Marshal(call)
	if err != nil {
		r.log.Err(err).Msg("error marshalling call")
		return Call{}, err
	}

	funcs := r.genFuncs(seqPath)

	t, err := template.New("").Funcs(funcs).Parse(string(callBytes))
	if err != nil {
		r.log.Err(err).Msg("error parsing call as template")
		return Call{}, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, r.ctxVariables); err != nil {
		r.log.Err(err).Msg("error performing substitutions")
		return Call{}, err
	}

	var newCall Call
	if err := yaml.Unmarshal(buf.Bytes(), &newCall); err != nil {
		r.log.Err(err).Msg("marshalling back to yaml")
		return Call{}, err
	}

	return newCall, nil
}

func (r *Runner) executeJQ(body any, jq string) (any, error) {
	query, err := gojq.Parse(jq)
	if err != nil {
		r.log.Err(err).Str("jq", jq).Msg("error parsing jq query")
		return "", err
	}

	var outVal any
	iterCount := 0
	iter := query.Run(body)
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		iterCount += 1
		if iterCount > 1 {
			return "", fmt.Errorf("jq resulted in more than 1 value")
		}

		if err, ok := val.(error); ok {
			r.log.Err(err).Msg("error executing JQ")
			return "", err
		}

		outVal = val
	}

	return outVal, nil
}

func (r *Runner) executeJQString(body any, jq string) (string, error) {
	val, err := r.executeJQ(body, jq)
	if err != nil {
		return "", err
	}
	return val.(string), err
}
