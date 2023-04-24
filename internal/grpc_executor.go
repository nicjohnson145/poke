package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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

func (g *GRPCExecutor) fetchDescriptors(service string, host string, dialInsecure bool) (grpcurl.DescriptorSource, error) {
	ds, ok := g.descriptors[service]
	if ok {
		g.log.Debug().Str("service", service).Msg("descriptor already fetched, using cache")
		return ds, nil
	}

	conn, err := g.connection(host, dialInsecure)
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

func (g *GRPCExecutor) connection(host string, dialInsecure bool) (*grpc.ClientConn, error) {
	conn, ok := g.connections[host]
	if ok {
		g.log.Debug().Str("host", host).Msg("reusing existing connection")
		return conn, nil
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		g.log.Err(err).Msg("error getting SSL pool")
		return nil, err
	}
	creds := credentials.NewTLS(&tls.Config{RootCAs: certPool})
	if dialInsecure {
		creds = insecure.NewCredentials()
	}

	g.log.Debug().Str("host", host).Msg("acquiring connection")
	// TODO: header support
	ctx := metadata.NewOutgoingContext(context.Background(), grpcurl.MetadataFromHeaders([]string{}))
	conn, err = grpcurl.BlockingDial(ctx, "tcp", host, creds)
	if err != nil {
		g.log.Err(err).Msg("error dialing service")
		return nil, err
	}

	g.connections[host] = conn
	return conn, nil
}

func (g *GRPCExecutor) executeRPC(call Call) (map[string]any, codes.Code, error) {
	g.log.Debug().Msg("fetching descriptors")
	descriptor, err := g.fetchDescriptors(g.callToServiceName(call), call.ServiceHost, call.SkipVerify)
	if err != nil {
		g.log.Err(err).Msg("error fetching descriptor")
		return nil, 0, err
	}

	var input io.Reader = strings.NewReader("")
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
	conn, err := g.connection(call.ServiceHost, call.SkipVerify)
	if err != nil {
		g.log.Err(err).Msg("error dialing service")
		return nil, 0, err
	}

	headers := g.makeRequestHeaderList(call)
	err = grpcurl.InvokeRPC(ctx, descriptor, conn, call.Url, headers, handler, reqParser.Next)
	if err != nil {
		g.log.Err(err).Msg("error invoking RPC")
		return nil, status.Convert(err).Code(), err
	}

	if handler.Status.Code() != codes.OK {
		return nil, handler.Status.Code(), handler.Status.Err()
	}

	g.log.Debug().Msg("decoding body")
	var outBody map[string]any
	err = json.NewDecoder(outBytes).Decode(&outBody)
	if err != nil {
		if err != io.EOF {
			g.log.Err(err).Msg("error unmarshalling body")
			return nil, 0, err
		}
		// No body, move along
	}

	return outBody, codes.OK, nil
}

func (g *GRPCExecutor) makeRequestHeaderList(call Call) []string {
	headers := []string{}
	for key, val := range call.Headers {
		headers = append(headers, fmt.Sprintf("%v: %v", key, val))
	}
	return headers
}

func (g *GRPCExecutor) Execute(call Call) (*ExecuteResult, error) {
	body, code, err := g.executeRPC(call)

	return &ExecuteResult{
		StatusCode: int(code),
		Body: body,
		Error: err,
	}, nil
}
