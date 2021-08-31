package grpc_requestlog_test

import (
	"bytes"
	"testing"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/takuoki/golib/applog"
	grpc_requestlog "github.com/takuoki/golib/middleware/grpc/requestlog"
)

type assertingPingService struct {
	pb_testproto.TestServiceServer
}

func TestRequestLogTestSuite(t *testing.T) {
	buf := &bytes.Buffer{}
	s := &RequestLogTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}},
			ServerOpts: []grpc.ServerOption{
				grpc.UnaryInterceptor(grpc_requestlog.UnaryServerInterceptor(
					applog.NewBasicLogger(buf, applog.TimeFormatOption("15:04:05")),
				)),
			},
		},
		buf: buf,
	}
	suite.Run(t, s)
}

type RequestLogTestSuite struct {
	*grpc_testing.InterceptorTestSuite
	buf *bytes.Buffer
}

func (s *RequestLogTestSuite) TestUnary_RequestLog() {
	s.buf.Reset()
	s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "something", SleepTimeMs: 9999})
	assert.Regexp(
		s.T(),
		`^{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"request log","labels":{"content_type":"application/grpc","ip_address":"127.0.0.1:[0-9]+","service_method":"/mwitkow.testproto.TestService/Ping","user_agent":"grpc-go/.+"}}`+"\n$",
		s.buf.String(),
	)
}
