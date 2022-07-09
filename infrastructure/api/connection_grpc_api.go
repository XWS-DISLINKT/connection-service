package api

import (
	"connection-service/application"
	"connection-service/domain"
	"connection-service/infrastructure/services"
	"connection-service/startup/config"
	"context"
	"fmt"
	pb "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
	profile "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
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

func (handler *ConnectionHandler) UpdateUser(ctx context.Context, request *pb.User) (*pb.Status, error) {
	user := domain.User{Id: request.UserId, IsPrivate: request.IsPrivate}
	success, err := handler.service.UpdateUser(&user)
	response := &pb.Status{Success: success}
	return response, err
}

func (handler *ConnectionHandler) MakeConnectionWithPublicProfile(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	requestSenderId := request.GetRequestSenderId()
	requestReceiverId := request.GetRequestReceiverId()
	success, err := handler.service.MakeConnectionWithPublicProfile(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}

	//kreiranje notifikacije
	cfg := config.NewConfig()
	profileAddress := fmt.Sprintf(cfg.ProfileServiceHost + ":" + cfg.ProfileServicePort)
	profileResponse, _ := services.ProfilesClient(profileAddress).SendNotification(context.TODO(),
		&profile.NewNotificationRequest{SenderId: requestSenderId, ReceiverId: requestReceiverId, NotificationType: "connection"})

	fmt.Printf("\ncreated notification {%s}", profileResponse.Id)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (handler *ConnectionHandler) MakeConnectionRequest(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	requestSenderId := request.GetRequestSenderId()
	requestReceiverId := request.GetRequestReceiverId()
	success, err := handler.service.MakeConnectionRequest(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}
	if err != nil {
		return response, err
	}

	//kreiranje notifikacije
	cfg := config.NewConfig()
	profileAddress := fmt.Sprintf(cfg.ProfileServiceHost + ":" + cfg.ProfileServicePort)
	profileResponse, _ := services.ProfilesClient(profileAddress).SendNotification(context.TODO(),
		&profile.NewNotificationRequest{SenderId: requestSenderId, ReceiverId: requestReceiverId, NotificationType: "request"})

	fmt.Printf("\ncreated notification {%s}", profileResponse.Id)

	return response, nil
}

func (handler *ConnectionHandler) ApproveConnectionRequest(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	requestSenderId := request.GetRequestSenderId()
	requestReceiverId := request.GetRequestReceiverId()
	success, err := handler.service.ApproveConnectionRequest(requestSenderId, requestReceiverId)
	response := &pb.ConnectionResponse{
		Success: success,
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (handler *ConnectionHandler) BlockConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	requestSenderId := request.GetRequestSenderId()
	blockedUserId := request.GetRequestReceiverId()
	success, err := handler.service.BlockUser(requestSenderId, blockedUserId)
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

func (handler *ConnectionHandler) GetBlockedConnectionsUsernames(ctx context.Context, request *pb.GetConnectionsUsernamesRequest) (*pb.GetConnectionsUsernamesResponse, error) {
	userId := request.GetId()
	usernames, err := handler.service.GetBlockedConnectionsUsernames(userId)
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

func (handler *ConnectionHandler) GetSuggestionIdsFor(ctx context.Context, request *pb.GetSuggestionIdsRequest) (*pb.GetSuggestionIdsResponse, error) {
	userId := request.GetId()
	usernames, err := handler.service.GetSuggestionIdsFor(userId)
	if err != nil {
		return nil, err
	}
	response := &pb.GetSuggestionIdsResponse{
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
	user1 := domain.User{Id: "623b0cc3a34d25d8567f9f82", IsPrivate: false}
	user2 := domain.User{Id: "623b0cc3a34d25d8567f9f83", IsPrivate: false}
	user3 := domain.User{Id: "623b0cc3a34d25d8567f9f84", IsPrivate: false}
	user4 := domain.User{Id: "623b0cc3a34d25d8567f9f87", IsPrivate: true}
	user5 := domain.User{Id: "623b0cc3a34d25d8567f9f88", IsPrivate: true}
	handler.service.DeleteEverything()
	handler.service.InsertUser(&user1)
	handler.service.InsertUser(&user2)
	handler.service.InsertUser(&user3)
	handler.service.InsertUser(&user4)
	handler.service.InsertUser(&user5)
	handler.service.MakeConnectionWithPublicProfile(user1.Id, user2.Id)
	handler.service.MakeConnectionRequest(user3.Id, user4.Id)
	handler.service.ApproveConnectionRequest(user3.Id, user4.Id)
	handler.service.MakeConnectionWithPublicProfile(user3.Id, user2.Id)
	handler.service.MakeConnectionRequest(user1.Id, user4.Id)
	connections, _ := handler.service.GetConnectionsUsernamesFor(user1.Id)
	requests, _ := handler.service.GetRequestsUsernamesFor(user1.Id)
	fmt.Println(connections)
	fmt.Println(requests)
}
