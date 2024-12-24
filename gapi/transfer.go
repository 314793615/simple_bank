package gapi

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func transferUser (user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}