package server

import (
	context "context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/SerjRamone/vaultme/internal/models"
	pb "github.com/SerjRamone/vaultme/pkg/vaultme_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer - gRPC server for handling user requests
type UserServer struct {
	pb.UnimplementedUsersServer
	log  *zap.Logger
	stor models.UserStorage
}

// NewUserServer - returns new UserServer instance
func NewUserServer(log *zap.Logger, storage models.UserStorage) *UserServer {
	return &UserServer{
		log:  log,
		stor: storage,
	}
}

// Register is a user registration handler
func (s *UserServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var u models.UserDTO
	u.Login = request.GetLogin()
	u.Password = request.GetPassword()

	hashed, err := getHash(u.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hashing password error: %v", err)
	}
	u.Password = hashed

	user, err := u.CreateUser(ctx, s.stor)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "not unique user: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "registration error: %v", err)
	}

	var resp pb.RegisterResponse
	resp.User = &pb.User{Id: user.ID}
	return &resp, nil
}

// Login is a user login handler
func (s *UserServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var u models.UserDTO
	u.Login = request.GetLogin()
	u.Password = request.GetPassword()

	user, err := u.GetUser(ctx, s.stor)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotExists) {
			return nil, status.Errorf(codes.NotFound, "getting user error: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "getting user error: %v", err)
	}

	if !validatePassword(user.PasswordHash, u.Password) {
		return nil, status.Errorf(codes.NotFound, "invalid credentials: %v", err)
	}

	var resp pb.LoginResponse
	resp.User = &pb.User{Id: user.ID}
	return &resp, nil
}

// getHash returns hash from password
func getHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("getting hash from password error: %w", err)
	}
	return string(bytes), nil
}

// validatePassword compares password hash with password
func validatePassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
