package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/314793615/simplebank/db/sqlc"

	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	username := req.GetUsername()

	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "the user is not exited: %w", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get the user: %w", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "the password is not correct: %w", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token %w", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token %w", err)
	}
	mtdt := server.extractMetadata(ctx)

	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.userAgent,
		ClientIp:     mtdt.clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session %w", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  transferUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		SessionId:             session.ID.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil

}
