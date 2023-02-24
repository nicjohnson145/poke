package internal

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
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
	}
}

type Runner struct {
	log          zerolog.Logger
	httpExecutor Executor
	parser       Parser
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
	for idx, call := range seq.Calls {
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
	}

	return nil
}

//nolint:ireturn
func (r *Runner) getClient(typ RequestType) (Executor, error) {
	switch typ {
	case RequestTypeHttp:
		return r.httpExecutor, nil
	default:
		return nil, fmt.Errorf("unhandled client type of '%v'", typ)
	}
}
