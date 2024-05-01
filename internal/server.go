package internal

import (
	"context"
	"net"

	"github.com/osousa/resticky/internal/protobuf/resticky"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server interface {
	LockAll(context.Context, *resticky.RestickyRequest) (*resticky.RestickyResponse, error)
	UnlockAll(context.Context, *resticky.RestickyRequest) (*resticky.RestickyResponse, error)
	Start()
	Stop()
}

type srvr struct {
	resticky.UnimplementedRestickyServiceServer // Embed UnimplementedServer
	grpcServer                                  *grpc.Server
	connManager                                 ConnManager
	lis                                         net.Listener
	Log                                         *zap.Logger
	Mode                                        string
}

type ServerConfig *ServerCfg

type ServerCfg struct {
	Port string
	Host string
	Mode string
}

func NewServerConfig() ServerConfig {
	return &ServerCfg{}
}

func ProvideServer(cfg ServerConfig, log *zap.Logger, connMan ConnManager) Server {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Error("failed to listen", zap.Error(err))
		return nil
	}
	grpcServer := grpc.NewServer()
	s := &srvr{
		lis:         lis,
		Log:         log,
		Mode:        cfg.Mode,
		connManager: connMan,
		grpcServer:  grpcServer,
	}

	resticky.RegisterRestickyServiceServer(grpcServer, s)

	return s
}

func (s *srvr) LockAll(ctx context.Context, r *resticky.RestickyRequest) (*resticky.RestickyResponse, error) {
	// Implement your server logic here
	err := s.connManager.LockAll()
	if err != nil {
		return &resticky.RestickyResponse{
			Success: false,
		}, err
	}

	// Return a response
	return &resticky.RestickyResponse{
		Success: true,
	}, nil

}

func (s *srvr) UnlockAll(ctx context.Context, r *resticky.RestickyRequest) (*resticky.RestickyResponse, error) {
	// Implement your server logic here
	err := s.connManager.UnlockAll()
	if err != nil {
		return &resticky.RestickyResponse{
			Success: false,
		}, err
	}

	// Return a response
	return &resticky.RestickyResponse{
		Success: true,
	}, nil
}

func (s *srvr) Start() {
	go func() {
		if err := s.grpcServer.Serve(s.lis); err != nil {
			s.Log.Error("failed to serve", zap.Error(err))
			panic(err)
		}
	}()
}

func (s *srvr) Stop() {
	ctx := context.Background()
	s.Log.Info("Stopping server")
	s.UnlockAll(ctx, nil)
	s.grpcServer.GracefulStop()
}
