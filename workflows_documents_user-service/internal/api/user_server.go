package api

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/api/proto/userpb"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const servTimeout = 2 * time.Second

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
	done := make(chan error, 1)
	go func() {
		select {
		case done <- s.UserService.DeactivateUserByID(activateUserRequest.Id):
		case <-ctx.Done():
		}
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
		select {
		case done <- s.UserService.Update(userDTO):
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		close(done)
		return nil, ctx.Err()
	case r := <-done:
		if r != nil {
			return nil, r
		}
		return &userpb.UpdateUserResponse{}, nil
	}
}

func (s *UserGRPCServer) UpdateUserPassword(ctx context.Context, req *userpb.UpdateUserPasswordRequest) (*userpb.UpdateUserPasswordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	doneChan := make(chan error, 1)
	go func() {
		dto := &dto.ChangePassword{
			Password:    req.GetOldPassword(),
			NewPassword: req.GetNewPassword(),
		}
		err := s.UserService.UpdatePassword(req.GetId(), dto)
		select {
		case <-ctx.Done():
		case doneChan <- err:
		}
	}()
	select {
	case <-ctx.Done():
		close(doneChan)
		return nil, ctx.Err()
	case err := <-doneChan:
		if err != nil {
			return nil, err
		}
		return &userpb.UpdateUserPasswordResponse{}, nil
	}
}

func (s *UserGRPCServer) FindUserByID(ctx context.Context, req *userpb.FindUserByIDRequest) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	doneChan := make(chan struct {
		user *dto.UserDTO
		err  error
	}, 1)
	go func() {
		user, err := s.UserService.FindByID(req.GetId())
		select {
		case doneChan <- struct {
			user *dto.UserDTO
			err  error
		}{user, err}:
		case <-ctx.Done():
		}
	}()
	select {
	case resp := <-doneChan:
		if resp.err != nil {
			return nil, resp.err
		}
		return &userpb.User{
			Id:       resp.user.ID,
			Name:     resp.user.Name,
			LastName: resp.user.LastName,
			Email:    resp.user.Email,
			Role:     resp.user.Role,
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *UserGRPCServer) FindUserByEmail(ctx context.Context, req *userpb.FindUserByEmailRequest) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	doneChan := make(chan struct {
		user *dto.UserDTO
		err  error
	}, 1)
	go func() {
		user, err := s.UserService.FindByEmail(req.Email)
		select {
		case doneChan <- struct {
			user *dto.UserDTO
			err  error
		}{user, err}:
		case <-ctx.Done():
		}
	}()
	select {
	case resp := <-doneChan:
		if resp.err != nil {
			return nil, resp.err
		}
		return &userpb.User{
			Id:       resp.user.ID,
			Name:     resp.user.Name,
			LastName: resp.user.LastName,
			Email:    resp.user.Email,
			Role:     resp.user.Role,
		}, nil
	case <-ctx.Done():
		close(doneChan)
		return nil, ctx.Err()
	}
}

func (s *UserGRPCServer) FindUserByName(ctx context.Context, req *userpb.FindUserByNameRequest) (*userpb.FindUserByNameResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, servTimeout)
	defer cancel()
	doneChan := make(chan struct {
		users []*dto.UserDTO
		err   error
	}, 1)
	go func() {
		users, err := s.UserService.FindByName(req.Name)
		select {
		case <-ctx.Done():
		case doneChan <- struct {
			users []*dto.UserDTO
			err   error
		}{users, err}:
		}
	}()
	select {
	case <-ctx.Done():
		close(doneChan)
		return nil, ctx.Err()
	case resp := <-doneChan:
		if resp.err != nil {
			return nil, resp.err
		}
		users := make([]*userpb.User, len(resp.users))
		for i := range len(resp.users) {
			users[i] = &userpb.User{
				Id:       resp.users[i].ID,
				Name:     resp.users[i].Name,
				LastName: resp.users[i].LastName,
				Email:    resp.users[i].Email,
				Role:     resp.users[i].Role,
			}
		}
		return &userpb.FindUserByNameResponse{Users: users}, nil
	}
}

func (s *UserGRPCServer) Run(userService *service.UserService) {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	userpb.RegisterUserServiceServer(grpcServer, &UserGRPCServer{UserService: userService})
	log.Println("server started on port: 8089")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
