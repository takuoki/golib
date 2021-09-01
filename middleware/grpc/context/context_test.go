package grpc_context_test

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/applog"
	grpc_context "github.com/takuoki/golib/middleware/grpc/context"
)

const (
	requestIDKey     = "request-id"
	requestID        = "request-id-001"
	authorizationKey = "authorization"
	authorization    = "authorization-001"
)

type assertingPingService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func (s *assertingPingService) Ping(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.PingResponse, error) {
	reqID := appctx.RequestID(ctx)
	switch ping.Value {
	case "exist-request-id":
		assert.Equal(s.T, requestID, reqID, "request ID doesn't match")
	case "not-exist-request-id":
		_, err := uuid.Parse(reqID)
		assert.Nil(s.T, err, "request ID is not a UUID format")
	}
	auth := appctx.Authorization(ctx)
	switch ping.Value {
	case "exist-authorization":
		assert.Equal(s.T, authorization, auth, "authorization doesn't match")
	}
	return s.TestServiceServer.Ping(ctx, ping)
}

func TestRequestLogTestSuite(t *testing.T) {
	s := &RequestLogTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.UnaryInterceptor(grpc_context.UnaryServerInterceptor(
					applog.NewBasicLogger(io.Discard),
					grpc_context.RequestIDKey(requestIDKey),
					grpc_context.AuthorizationKey(authorizationKey),
				)),
			},
		},
	}
	suite.Run(t, s)
}

type RequestLogTestSuite struct {
	*grpc_testing.InterceptorTestSuite
}

func (s *RequestLogTestSuite) TestUnary_ExistRequestID() {
	md := metadata.New(map[string]string{requestIDKey: requestID})
	_, _ = s.Client.Ping(
		metadata.NewOutgoingContext(s.SimpleCtx(), md),
		&pb_testproto.PingRequest{Value: "exist-request-id", SleepTimeMs: 9999},
		grpc.Header(&md),
	)
}

func (s *RequestLogTestSuite) TestUnary_NotExistRequestID() {
	_, _ = s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "not-exist-request-id", SleepTimeMs: 9999})
}

func (s *RequestLogTestSuite) TestUnary_ExistAuthorization() {
	md := metadata.New(map[string]string{authorizationKey: authorization})
	_, _ = s.Client.Ping(
		metadata.NewOutgoingContext(s.SimpleCtx(), md),
		&pb_testproto.PingRequest{Value: "exist-authorization", SleepTimeMs: 9999},
		grpc.Header(&md),
	)
}
