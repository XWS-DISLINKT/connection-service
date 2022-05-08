package api

import (
	"connection-service/application"
	"connection-service/domain"
	"context"
	"fmt"
	pb "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
)

type ConnectionHandler struct {
	pb.UnsafeConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}
func (handler *ConnectionHandler) InsertUser(ctx context.Context, request *pb.User) (*pb.Status, error) {
	user := domain.User{Id: request.UserId, IsPrivate: request.IsPrivate}
	success, err := handler.service.InsertUser(&user)
	response := &pb.Status{Success: success}
	return response, err
}

func (handler *ConnectionHandler) MakeConnectionWithPublicProfile(ctx context.Context, request *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	requestSenderId := request.ConnectionBody.GetRequestSenderId()
	requestReceiverId := request.ConnectionBody.GetRequestReceiverId()
	success, err := handler.service.MakeConnectionWithPublicProfile(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (handler *ConnectionHandler) MakeConnectionRequest(ctx context.Context, request *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	requestSenderId := request.ConnectionBody.GetRequestSenderId()
	requestReceiverId := request.ConnectionBody.GetRequestReceiverId()
	success, err := handler.service.MakeConnectionRequest(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (handler *ConnectionHandler) ApproveConnectionRequest(ctx context.Context, request *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	requestSenderId := request.ConnectionBody.GetRequestSenderId()
	requestReceiverId := request.ConnectionBody.GetRequestReceiverId()
	success, err := handler.service.ApproveConnectionRequest(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (handler *ConnectionHandler) GetConnectionsUsernamesFor(ctx context.Context, request *pb.GetConnectionsUsernamesRequest) (*pb.GetConnectionsUsernamesResponse, error) {
	userId := request.GetId()
	usernames, err := handler.service.GetConnectionsUsernamesFor(userId)
	if err != nil {
		return nil, err
	}
	response := &pb.GetConnectionsUsernamesResponse{
		Usernames: []string{},
	}

	for _, username := range usernames {
		response.Usernames = append(response.Usernames, username)
	}

	return response, nil
}

func (handler *ConnectionHandler) GetRequestsUsernamesFor(ctx context.Context, request *pb.GetConnectionsUsernamesRequest) (*pb.GetConnectionsUsernamesResponse, error) {
	userId := request.GetId()
	usernames, err := handler.service.GetRequestsUsernamesFor(userId)
	if err != nil {
		return nil, err
	}
	response := &pb.GetConnectionsUsernamesResponse{
		Usernames: []string{},
	}

	for _, username := range usernames {
		response.Usernames = append(response.Usernames, username)
	}

	return response, nil
}

func (handler *ConnectionHandler) Demo() {
	user1 := domain.User{Id: "1", IsPrivate: false}
	user2 := domain.User{Id: "2", IsPrivate: false}
	user3 := domain.User{Id: "3", IsPrivate: true}
	user4 := domain.User{Id: "4", IsPrivate: true}
	user5 := domain.User{Id: "5", IsPrivate: true}
	handler.service.DeleteEverything()
	handler.service.InsertUser(&user1)
	handler.service.InsertUser(&user2)
	handler.service.InsertUser(&user3)
	handler.service.InsertUser(&user4)
	handler.service.InsertUser(&user5)
	handler.service.MakeConnectionWithPublicProfile(user1.Id, user2.Id)
	handler.service.MakeConnectionRequest(user3.Id, user4.Id)
	handler.service.ApproveConnectionRequest(user3.Id, user4.Id)
	handler.service.MakeConnectionRequest(user1.Id, user4.Id)
	connections, _ := handler.service.GetConnectionsUsernamesFor(user1.Id)
	requests, _ := handler.service.GetRequestsUsernamesFor(user1.Id)
	fmt.Println(connections)
	fmt.Println(requests)
}
