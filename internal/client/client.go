package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SerjRamone/vaultme/internal/config"
	"github.com/SerjRamone/vaultme/internal/models"
	pb "github.com/SerjRamone/vaultme/pkg/vaultme_v1"
)

// Client - is a VaultMe gRPC client
type Client struct {
	conn    *grpc.ClientConn
	address string
	log     *zap.Logger
}

// NewClient - returns new Client instance
func NewClient(log *zap.Logger, cfg *config.Client) (*Client, error) {
	c := &Client{
		address: cfg.ServerAddress,
		log:     log,
	}

	if err := c.initConnection(); err != nil {
		return nil, fmt.Errorf("init connection error: %w", err)
	}

	return c, nil
}

// getCredentials - returns gRPC credentials
// TODO: secure conn
func getCredentials() (credentials.TransportCredentials, error) {
	creds := insecure.NewCredentials()
	return creds, nil
}

func (c *Client) initConnection() error {
	var dopts []grpc.DialOption

	creds, err := getCredentials()
	if err != nil {
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	dopts = append(dopts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.NewClient(c.address, dopts...)
	if err != nil {
		return fmt.Errorf("gRPC client creation error: %w", err)
	}

	c.conn = conn

	return nil
}

// CreateUser - registration user request
func (c *Client) CreateUser(ctx context.Context, us *models.UserDTO) (*models.User, error) {
	resp, err := pb.NewUsersClient(c.conn).Register(ctx, &pb.RegisterRequest{
		Login:    us.Login,
		Password: us.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("user registration error: %w", err)
	}

	return &models.User{
		ID:           resp.GetUser().GetId(),
		Login:        us.Login,
		PasswordHash: us.Password,
	}, nil
}

// GetUser - login user request
func (c *Client) GetUser(ctx context.Context, u *models.UserDTO) (*models.User, error) {
	resp, err := pb.NewUsersClient(c.conn).Login(ctx, &pb.LoginRequest{
		Login:    u.Login,
		Password: u.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("user login error: %w", err)
	}

	return &models.User{
		ID:           resp.GetUser().GetId(),
		Login:        u.Login,
		PasswordHash: u.Password,
	}, nil
}
