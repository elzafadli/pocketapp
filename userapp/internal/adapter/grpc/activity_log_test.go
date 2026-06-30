package grpc

import (
	"net"
	"testing"

	"userapp/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	go_grpc "google.golang.org/grpc"
)

func TestActivityLogAdapter_Lifecycle(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	defer listener.Close()

	grpcServer := go_grpc.NewServer()
	go func() {
		_ = grpcServer.Serve(listener)
	}()
	defer grpcServer.Stop()

	addr := listener.Addr().String()
	adapter := &ActivityLog{
		Conf: &config.Config{
			Grpc: config.GrpcConfig{
				ActivityLogAddr: addr,
			},
		},
	}

	err = adapter.Startup()
	assert.NoError(t, err)
	assert.NotNil(t, adapter.GetConn())

	err = adapter.Shutdown()
	assert.NoError(t, err)
}
