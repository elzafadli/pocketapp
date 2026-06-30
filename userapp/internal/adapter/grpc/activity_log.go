package grpc

import (
	"fmt"

	"github.com/elzafadli/bookrpc/pb"
	"userapp/config"

	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ActivityLogConnection interface {
	GetConn() *go_grpc.ClientConn
	GetClient() pb.ActivityLogServiceClient
}

type ActivityLog struct {
	Conf   *config.Config `inject:"config"`
	conn   *go_grpc.ClientConn
	client pb.ActivityLogServiceClient
}

func (a *ActivityLog) Startup() error {
	addr := a.Conf.Grpc.ActivityLogAddr

	conn, err := go_grpc.NewClient(addr, go_grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create activity log client at %s: %w", addr, err)
	}

	fmt.Println(addr, "up grpc")

	a.conn = conn
	a.client = pb.NewActivityLogServiceClient(conn)
	return nil
}

func (a *ActivityLog) Shutdown() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *ActivityLog) GetConn() *go_grpc.ClientConn {
	return a.conn
}

func (a *ActivityLog) GetClient() pb.ActivityLogServiceClient {
	return a.client
}
