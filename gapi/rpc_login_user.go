package gapi

import (
	"context"
	"database/sql"

	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	username := req.GetUsername()

	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "the user is not exited: %w", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get the user: %w", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "the password is not correct: %w", err)
	}

	accessToken, payload, err := server.tokenMaker.CreateToken(username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create token %w", err)
	}
	// CREATE TABLE "sessions" (
	// 	"id" bigserial PRIMARY KEY 
	// 	 "username" varchar NOT NULL,
	// 	 "refresh_token" varchar NOT NULL,
	// 	 "user_agent" varchar NOT NULL,
	// 	 "client_ip" varchar NOT NULL,
	// 	 "is_blocked" boolean NOT NULL DEFAULT false,
	// 	 "expires_at" timestamptz NOT NULL,
	// 	 "created_at" timestamptz NOT NULL DEFAULT (now())
	// );
	arg := db.CreateSessionParams{
		
	}
	session, err := server.store.CreateSession()

	rsp := &pb.LoginUserResponse{
		User: transferUser(user),
		AccessToken: accessToken,
		AccessTokenExpiresAt: timestamppb.New(payload.ExpiredAt),
	
	}
	return rsp, nil

}