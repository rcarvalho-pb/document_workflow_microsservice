package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/api/proto/userpb"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/service"
	"google.golang.org/grpc"
)

const servTimeout = 10 * time.Second

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	UserService *service.UserService
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, createUserRequest *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()

	type resp struct {
		resp *userpb.CreateUserResponse
		err  error
	}

	respChan := make(chan *resp)

	go func() {
		userDTO := &dto.UserDTO{
			Name:     createUserRequest.Name,
			LastName: createUserRequest.LastName,
			Email:    createUserRequest.Email,
			Password: createUserRequest.Password,
			Role:     createUserRequest.Role,
		}
		id, err := s.UserService.Save(userDTO)
		if err != nil {
			respChan <- &resp{nil, err}
			return
		}
		response := &resp{
			resp: &userpb.CreateUserResponse{Id: id},
			err:  nil,
		}
		respChan <- response
	}()

	select {
	case <-ctx.Done():
		close(respChan)
		return nil, ctx.Err()
	case resp := <-respChan:
		return resp.resp, resp.err
	}
}

func (s *UserGRPCServer) DeactivateUser(ctx context.Context, deactivateUserRequest *userpb.DeactivateUserRequest) (*userpb.DeactivateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	done := make(chan error)
	go func() {
		err := s.UserService.DeactivateUserByID(deactivateUserRequest.GetId())
		done <- err
	}()
	select {
	case <-ctx.Done():
		close(done)
		return nil, ctx.Err()
	case r := <-done:
		return &userpb.DeactivateUserResponse{}, r
	}
}

func (s *UserGRPCServer) ActivateUser(ctx context.Context, activateUserRequest *userpb.ActivateUserRequest) (*userpb.ActivateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	done := make(chan error)
	go func() {
		err := s.UserService.DeactivateUserByID(activateUserRequest.GetId())
		done <- err
	}()
	select {
	case <-ctx.Done():
		close(done)
		return nil, ctx.Err()
	case r := <-done:
		return &userpb.ActivateUserResponse{}, r
	}
}

func (s *UserGRPCServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		userDTO := &dto.UserDTO{
			ID:       req.GetId(),
			Name:     req.GetName(),
			LastName: req.GetLastName(),
			Email:    req.GetEmail(),
			Role:     req.GetRole(),
		}
		err := s.UserService.Update(userDTO)
		select {
		case d one <- err:
		case <-ct x.Done():
		}
	}()
	select {
	case <-ctx.Done():
		close(done)
		return nil, ctx.Err()
	case r := <-done:
		return &userpb.UpdateUserResponse{}, r
	}
}

func (s *UserGRPCServer) UpdateUserPassword(ctx context.Context, req *userpb.UpdateUserPasswordRequest) (*userpb.UpdateUserPasswordResponse, error) {
	dto := &dto.ChangePassword{
		Password:    req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}
	err := s.UserService.UpdatePassword(req.GetId(), dto)
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateUserPasswordResponse{}, nil
}

func (s *UserGRPCServer) FindUserByID(ctx context.Context, req *userpb.FindUserByIDRequest) (*userpb.User, error) {
	user, err := s.UserService.FindByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return &userpb.User{
		Id:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserGRPCServer) FindUserByEmail(ctx context.Context, req *userpb.FindUserByEmailRequest) (*userpb.User, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserGRPCServer) FindUserByName(ctx context.Context, req *userpb.FindUserByNameRequest) (*userpb.FindUserByNameResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserGRPCServer) Run() {
	lis, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer()
}
