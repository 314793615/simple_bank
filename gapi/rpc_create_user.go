package gapi

import (
	"time"

	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"github.com/314793615/simplebank/val"
	"github.com/314793615/simplebank/worker"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validCreateUserParams(req)
	if violations !=  nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %w", err)
	}
	createUserParam := db.CreatUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
	}

	user, err := server.store.CreatUser(ctx, createUserParam)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %w", err)
	}

	sendEmailPayload := worker.SendEmailPayload{
		Username: user.Username,
	}
	options := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, sendEmailPayload, options...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute email payload: %w", err)
	}
	return &pb.CreateUserResponse{
		User: transferUser(user),
	}, nil
}


func validCreateUserParams(req *pb.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	err := val.ValidUserName(req.Username)
	if err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	err = val.ValidFullName(req.FullName)
	if err != nil {
		violations = append(violations, FieldViolation("fullname", err))
	}

	err = val.ValidPassword(req.Password)
	if err != nil {
		violations = append(violations, FieldViolation("password", err))
	}

	err  = val.ValidEmail(req.Email)
	if err != nil {
		violations = append(violations, FieldViolation("email", err))
	}
	return violations

}