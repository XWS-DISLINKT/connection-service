package api

import (
	"connection-service/application"
	"connection-service/domain"
	events "github.com/XWS-DISLINKT/dislinkt/common/saga/create_user"
	saga "github.com/XWS-DISLINKT/dislinkt/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		connectionService: connectionService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	user := &domain.User{Id: command.User.Id, IsPrivate: command.User.IsPrivate}
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.CreateUserConnection:
		success, err := handler.connectionService.InsertUser(user)
		if err != nil || !success {
			reply.Type = events.UserConnectionNotCreated
		} else {
			reply.Type = events.UserConnectionCreated
		}
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
