package internal

import "github.com/rs/zerolog"

type GRPCExecutorOpts struct {
	Logger zerolog.Logger
}

func NewGRPCExecutor(opts GRPCExecutorOpts) *GRPCExecutor {
	return &GRPCExecutor{
		log: opts.Logger,
	}
}

var _ Executor = (*GRPCExecutor)(nil)

type GRPCExecutor struct {
	log zerolog.Logger
}

func (g *GRPCExecutor) Execute(call Call) (*ExecuteResult, error) {
	return nil, nil
}
