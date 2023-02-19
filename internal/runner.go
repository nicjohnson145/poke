package internal

import "github.com/rs/zerolog"

type RunnerOpts struct {
	Logger zerolog.Logger
}

func NewRunner(opts RunnerOpts) *Runner {
	return &Runner{
		log: opts.Logger,
	}
}

type Runner struct {
	log zerolog.Logger
}

func (r *Runner) RunSequence(seq Sequence) error {
	return nil
}
