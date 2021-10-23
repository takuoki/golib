package grpc_error_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
	grpc_error "github.com/takuoki/golib/middleware/grpc/error"
	"github.com/takuoki/golib/notice"
)

const (
	domain                  = "dummy.domain"
	internalServerErrorCode = "SERVER_ERROR"

	apperrClientStatus  = codes.InvalidArgument
	apperrClientCode    = "APPERROR_CLIENT"
	apperrClientMessage = "this is apperr client"

	apperrServerStatus  = codes.Internal
	apperrServerCode    = "APPERROR_SERVER"
	apperrServerMessage = "this is apperr server"
	apperrServerLog     = "this is apperr server log"

	generalErrorMessage      = "this is general error"
	notificationErrorMessage = "return-error"
)

type assertingPingService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func (s *assertingPingService) Ping(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.PingResponse, error) {
	switch ping.Value {
	case "apperr-client":
		return nil, apperr.NewClientError(codes.InvalidArgument, apperrClientCode, apperrClientMessage)
	case "apperr-server":
		return nil, apperr.NewServerError(codes.Internal, apperrServerCode, apperrServerMessage, apperrServerLog)
	case "general-error":
		return nil, errors.New(generalErrorMessage)
	case "notification-error":
		return nil, errors.New(notificationErrorMessage)
	}
	return s.TestServiceServer.Ping(ctx, ping)
}

func TestErrorHandlerTestSuite(t *testing.T) {
	lBuf := &bytes.Buffer{}
	nBuf := &bytes.Buffer{}
	s := &ErrorHandlerTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.UnaryInterceptor(
					grpc_error.UnaryServerInterceptor(
						domain,
						internalServerErrorCode,
						applog.NewBasicLogger(lBuf, applog.TimeFormatOption("15:04:05")),
						newTestNotifier(nBuf),
					),
				),
			},
		},
		lBuf: lBuf,
		nBuf: nBuf,
	}
	suite.Run(t, s)
}

type ErrorHandlerTestSuite struct {
	*grpc_testing.InterceptorTestSuite
	lBuf *bytes.Buffer
	nBuf *bytes.Buffer
}

func (s *ErrorHandlerTestSuite) TestUnary_Success() {
	s.lBuf.Reset()
	s.nBuf.Reset()
	_, err := s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "success", SleepTimeMs: 9999})
	assert.Nil(s.T(), err)
	assert.Empty(s.T(), s.lBuf.String(), "log must be empty")
	assert.Empty(s.T(), s.nBuf.String(), "notification must be empty")
}

func (s *ErrorHandlerTestSuite) TestUnary_ApperrClient() {
	s.lBuf.Reset()
	s.nBuf.Reset()
	_, err := s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "apperr-client", SleepTimeMs: 9999})
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), apperrClientStatus, st.Code(), "status doesn't match")
		assert.Equal(s.T(), apperrClientMessage, st.Message(), "message doesn't match")
		if assert.Len(s.T(), st.Details(), 1, "length of error details") {
			if ed, ok := st.Details()[0].(*errdetails.ErrorInfo); ok {
				assert.Equal(s.T(), apperrClientCode, ed.Reason, "code doesn't match")
			} else {
				s.T().Error("error detail must be cast to errdetails.ErrorInfo")
			}
		}
	} else {
		s.T().Error("status.Status must be retrievable from error")
	}
	assert.Empty(s.T(), s.lBuf.String(), "log must be empty")
	assert.Empty(s.T(), s.nBuf.String(), "notification must be empty")
}

func (s *ErrorHandlerTestSuite) TestUnary_ApperrServer() {
	s.lBuf.Reset()
	s.nBuf.Reset()
	_, err := s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "apperr-server", SleepTimeMs: 9999})
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), apperrServerStatus, st.Code(), "status doesn't match")
		assert.Equal(s.T(), apperrServerMessage, st.Message(), "message doesn't match")
		if assert.Len(s.T(), st.Details(), 1, "length of error details") {
			if ed, ok := st.Details()[0].(*errdetails.ErrorInfo); ok {
				assert.Equal(s.T(), apperrServerCode, ed.Reason, "code doesn't match")
			} else {
				s.T().Error("error detail must be cast to errdetails.ErrorInfo")
			}
		}
	} else {
		s.T().Error("status.Status must be retrievable from error")
	}
	assert.Regexp(s.T(), `^{"time":"\d{2}:\d{2}:\d{2}","level":"ERROR","message":"this is apperr server log"}`+"\n$", s.lBuf.String(), "log message doesn't match")
	assert.Equal(s.T(), "ERROR: this is apperr server", s.nBuf.String(), "notification doesn't match")
}

func (s *ErrorHandlerTestSuite) TestUnary_GeneralErr() {
	s.lBuf.Reset()
	s.nBuf.Reset()
	_, err := s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "general-error", SleepTimeMs: 9999})
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), codes.Internal, st.Code(), "status doesn't match")
		assert.Equal(s.T(), "internal server error", st.Message(), "message doesn't match")
		if assert.Len(s.T(), st.Details(), 1, "length of error details") {
			if ed, ok := st.Details()[0].(*errdetails.ErrorInfo); ok {
				assert.Equal(s.T(), internalServerErrorCode, ed.Reason, "code doesn't match")
			} else {
				s.T().Error("error detail must be cast to errdetails.ErrorInfo")
			}
		}
	} else {
		s.T().Error("status.Status must be retrievable from error")
	}
	assert.Regexp(s.T(), `^{"time":"\d{2}:\d{2}:\d{2}","level":"ERROR","message":"this is general error"}`+"\n$", s.lBuf.String(), "log message doesn't match")
	assert.Equal(s.T(), "ERROR: this is general error", s.nBuf.String(), "notification doesn't match")
}

func (s *ErrorHandlerTestSuite) TestUnary_NotificationError() {
	s.lBuf.Reset()
	s.nBuf.Reset()
	_, err := s.Client.Ping(s.SimpleCtx(), &pb_testproto.PingRequest{Value: "notification-error", SleepTimeMs: 9999})
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), codes.Internal, st.Code(), "status doesn't match")
		assert.Equal(s.T(), "internal server error", st.Message(), "message doesn't match")
		if assert.Len(s.T(), st.Details(), 1, "length of error details") {
			if ed, ok := st.Details()[0].(*errdetails.ErrorInfo); ok {
				assert.Equal(s.T(), internalServerErrorCode, ed.Reason, "code doesn't match")
			} else {
				s.T().Error("error detail must be cast to errdetails.ErrorInfo")
			}
		}
	} else {
		s.T().Error("status.Status must be retrievable from error")
	}
	assert.Regexp(s.T(), `{"time":"\d{2}:\d{2}:\d{2}","level":"ERROR","message":"failed to send nortification: error"}`+"\n", s.lBuf.String(), "log message doesn't match")
	assert.Empty(s.T(), s.nBuf.String(), "notification must be empty")
}

type testNotifier struct {
	mu  sync.Mutex
	out io.Writer
}

func newTestNotifier(w io.Writer) notice.Notifier {
	return &testNotifier{out: w}
}

func (n *testNotifier) Error(err error) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if err.Error() == notificationErrorMessage {
		return errors.New("error")
	}
	fmt.Fprintf(n.out, "ERROR: %v", err)
	return nil
}

func (n *testNotifier) Critical(err error) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Fprintf(n.out, "CRITICAL: %v", err)
	return nil
}
