package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

type GRPCExecutorOpts struct {
	Logger zerolog.Logger
}

func NewGRPCExecutor(opts GRPCExecutorOpts) *GRPCExecutor {
	return &GRPCExecutor{
		log:         opts.Logger,
		descriptors: make(map[string]grpcurl.DescriptorSource),
		connections: make(map[string]*grpc.ClientConn),
	}
}

var _ Executor = (*GRPCExecutor)(nil)

type GRPCExecutor struct {
	log         zerolog.Logger
	descriptors map[string]grpcurl.DescriptorSource
	connections map[string]*grpc.ClientConn
}

func (g *GRPCExecutor) fetchDescriptors(service string, host string) (grpcurl.DescriptorSource, error) {
	ds, ok := g.descriptors[service]
	if ok {
		g.log.Debug().Str("service", service).Msg("descriptor already fetched, using cache")
		return ds, nil
	}

	conn, err := g.connection(host)
	if err != nil {
		return nil, err
	}
	g.log.Debug().Msg("fetching descriptors using reflection")
	ctx := context.Background()
	client := grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(conn))
	source := grpcurl.DescriptorSourceFromServer(ctx, client)

	g.descriptors[service] = source
	return source, nil
}

func (g *GRPCExecutor) callToServiceName(call Call) string {
	return strings.Split(call.Url, "/")[0]
}

func (g *GRPCExecutor) connection(host string) (*grpc.ClientConn, error) {
	conn, ok := g.connections[host]
	if ok {
		g.log.Debug().Str("host", host).Msg("reusing existing connection")
		return conn, nil
	}

	g.log.Debug().Str("host", host).Msg("acquiring connection")
	// TODO: header support
	ctx := metadata.NewOutgoingContext(context.Background(), grpcurl.MetadataFromHeaders([]string{}))
	conn, err := grpcurl.BlockingDial(ctx, "tcp", host, nil)
	if err != nil {
		g.log.Err(err).Msg("error dialing service")
		return nil, err
	}

	g.connections[host] = conn
	return conn, nil
}

func (g *GRPCExecutor) executeRPC(call Call) (map[string]any, codes.Code, error) {
	g.log.Debug().Msg("fetching descriptors")
	descriptor, err := g.fetchDescriptors(g.callToServiceName(call), call.ServiceHost)
	if err != nil {
		g.log.Err(err).Msg("error fetching descriptor")
		return nil, 0, err
	}

	var input io.Reader
	if call.Body != nil {
		bodyBytes, err := json.Marshal(call.Body)
		if err != nil {
			g.log.Err(err).Msg("error marshalling request body")
			return nil, 0, err
		}
		input = bytes.NewReader(bodyBytes)
	}

	g.log.Debug().Msg("building parser and formatter")
	reqParser, reqFormatter, err := grpcurl.RequestParserAndFormatter(
		grpcurl.FormatJSON,
		descriptor,
		input,
		grpcurl.FormatOptions{
			EmitJSONDefaultFields: false,
			IncludeTextSeparator:  false,
			AllowUnknownFields:    false,
		},
	)
	if err != nil {
		g.log.Err(err).Msg("error getting request formatter & parser")
		return nil, 0, err
	}

	outBytes := &bytes.Buffer{}

	handler := &grpcurl.DefaultEventHandler{
		Out:            outBytes,
		Formatter:      reqFormatter,
		VerbosityLevel: 0,
	}

	ctx := context.Background()
	conn, err := g.connection(call.ServiceHost)
	if err != nil {
		g.log.Err(err).Msg("error dialing service")
		return nil, 0, err
	}

	err = grpcurl.InvokeRPC(ctx, descriptor, conn, call.Url, []string{}, handler, reqParser.Next)
	if err != nil {
		g.log.Err(err).Msg("error invoking RPC")
		return nil, status.Convert(err).Code(), err
	}

	g.log.Debug().Msg("decoding body")
	var outBody map[string]any
	err = json.Unmarshal(outBytes.Bytes(), &outBody)
	if err != nil {
		g.log.Err(err).Msg("error unmarshalling body")
		return nil, 0, err
	}

	return outBody, codes.OK, nil
}

func (g *GRPCExecutor) Execute(call Call) (*ExecuteResult, error) {
	body, code, err := g.executeRPC(call)
	if err != nil {
		g.log.Err(err).Msg("error executing RPC")
		return nil, err
	}

	return &ExecuteResult{
		StatusCode: int(code),
		Body: body,
	}, nil
}
