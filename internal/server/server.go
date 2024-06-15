// Package server implements gRPC server
package server

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SerjRamone/vaultme/internal/config"
	"github.com/SerjRamone/vaultme/internal/models"
	pb "github.com/SerjRamone/vaultme/pkg/vaultme_v1"
)

// Server represents gRPC server
type Server struct {
	gRPCServer *grpc.Server
	log        *zap.Logger
	UserServer *UserServer
	address    string
}

// NewServer creates Server instance, initializes gRPC server
func NewServer(userStor models.UserStorage, log *zap.Logger, cfg *config.Server) (*Server, error) {
	s := &Server{
		address:    cfg.Address,
		log:        log,
		UserServer: NewUserServer(log, userStor),
	}

	// TODO: secure grpc server
	creds, err := getCredentials()
	if err != nil {
		return nil, err
	}

	options := grpc.ChainUnaryInterceptor(s.requestLogger())

	s.gRPCServer = grpc.NewServer(grpc.Creds(creds), options)

	return s, nil
}

// getCredentials - returns gRPC server credentials
// TODO: secure
func getCredentials() (credentials.TransportCredentials, error) {
	creds := insecure.NewCredentials()
	return creds, nil
}

// Serve - starts gRPC server
func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	pb.RegisterUsersServer(s.gRPCServer, s.UserServer)

	if err := s.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("starting gRPC server error: %w", err)
	}
	return nil
}

// Stop - stops gRPC server gracefully
func (s *Server) Stop() {
	s.gRPCServer.GracefulStop()
}
