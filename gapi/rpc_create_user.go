package gapi

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %w", err)
	}
	createUserParam := db.CreatUserParams{
		Username: req.GetUsername(),
		HashedPassword: hashedPassword,
		Email: req.GetEmail(),
		FullName: req.GetFullName(),
	}

	user, err := server.store.CreatUser(ctx, createUserParam)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %w", err)
	}

	return &pb.CreateUserResponse{
		User: transferUser(user),
	}, nil

}

