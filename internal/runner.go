package internal

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"text/template"

	"github.com/itchyny/gojq"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"github.com/google/go-cmp/cmp"
)

type RunnerOpts struct {
	Logger       zerolog.Logger
	HttpExecutor Executor
	Parser       Parser
}

func NewRunner(opts RunnerOpts) *Runner {
	return &Runner{
		log:          opts.Logger,
		httpExecutor: opts.HttpExecutor,
		parser:       opts.Parser,
		ctxVariables: make(map[string]string),
	}
}

type Runner struct {
	log          zerolog.Logger
	httpExecutor Executor
	parser       Parser
	ctxVariables map[string]string
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
		}
	}

	return errors.Join(errs...)
}

func (r *Runner) runSingleSequence(seq Sequence) error {
	for idx, c := range seq.Calls {
		call, err := r.evaluateTemplate(c)
		if err != nil {
			return err
		}

		exec, err := r.getClient(call.Type)
		if err != nil {
			return fmt.Errorf("error creating request client: %w", err)
		}

		result, err := exec.Execute(call)
		if err != nil {
			name := call.Name
			if name == "" {
				name = fmt.Sprintf("call_%v", idx)
			}
			return fmt.Errorf("error executing call %v: %w", name, err)
		}

		wantStatus := call.WantStatus
		if wantStatus == 0 {
			wantStatus = http.StatusOK
		}

		if result.StatusCode != wantStatus {
			r.log.Error().Interface("body", result.Body).Msg("body")
			return fmt.Errorf("got incorrect status: want (%v) got (%v)", wantStatus, result.StatusCode)
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
	default:
		return nil, fmt.Errorf("unhandled client type of '%v'", typ)
	}
}

func (r *Runner) evaluateTemplate(call Call) (Call, error) {
	callBytes, err := yaml.Marshal(call)
	if err != nil {
		r.log.Err(err).Msg("error marshalling call")
		return Call{}, err
	}

	t, err := template.New("").Parse(string(callBytes))
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
