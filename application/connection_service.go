package application

import (
	"connection-service/domain"
	"connection-service/infrastructure/persistence"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ConnectionService struct {
	databaseDriver neo4j.Driver
}

func NewConnectionService() *ConnectionService {
	databaseDriver, err := persistence.GetDriver()
	if err != nil {
		panic(err)
	}
	err = databaseDriver.VerifyConnectivity()
	if err != nil {
		panic(err)
	}
	return &ConnectionService{
		databaseDriver: databaseDriver,
	}
}

func (service *ConnectionService) DeleteEverything() error {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"match (n) detach delete n",
			map[string]interface{}{})
		return nil, err
	})
	return err
}

func (service *ConnectionService) InsertUser(user *domain.User) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	exists, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match (user:User {id: $id}) return user",
			map[string]interface{}{"id": user.Id})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return false, err
	}

	if exists == nil {
		return false, errors.New("Error.")
	}

	if exists.(bool) {
		return false, errors.New("User exists.")
	}

	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"create (usr:User {id: $id, isPrivate: $isPrivate}) return usr is not null",
			map[string]interface{}{"id": user.Id, "isPrivate": user.IsPrivate})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return successful.(bool), err
}

func (service *ConnectionService) UpdateUser(user *domain.User) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match (user:User {id: $id}) set user.isPrivate = $isPrivate return user is not null",
			map[string]interface{}{"id": user.Id, "isPrivate": user.IsPrivate})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if successful == nil {
		return false, errors.New("Error.")
	}

	return successful.(bool), err
}

func (service *ConnectionService) GetSuggestionIdsFor(userId string) ([]string, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	connections, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(
			"match(u1:User{id: $userId})-[c1:CONNECTS]->(u2) match (u2)-[c2:CONNECTS]->(u3)  where not u3.id = $userId and not (u1)-[:CONNECTS]->(u3) return distinct u3.id",
			map[string]interface{}{"userId": userId})

		users := Users{users: []string{}}
		if records == nil {
			return Users{users: []string{}}, nil
		}

		for records.Next() {
			record := records.Record()
			userId, _ := record.Get("u3.id")
			users.users = append(users.users, userId.(string))
		}
		return users, err
	})

	return connections.(Users).users, err
}

func (service *ConnectionService) MakeConnectionWithPublicProfile(requestSenderId string, requestReceiverId string) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	canMakeConnection, err := service.canMakeConnection(requestSenderId, requestReceiverId)

	if !canMakeConnection {
		return false, err
	}

	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match(u1:User{id: $requestSenderId}) match(u2:User{id: $requestReceiverId})"+
				" create (u1)-[c1:CONNECTS{isApproved: true}]->(u2)"+
				" create (u1)<-[c2:CONNECTS{isApproved: true}]-(u2)"+
				" create (u1)-[c3:CONNECTS{isBlocked: false}]->(u2)"+
				" create (u1)<-[c4:CONNECTS{isBlocked: false}]-(u2)"+
				" return c1 is not null",
			map[string]interface{}{"requestSenderId": requestSenderId, "requestReceiverId": requestReceiverId})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return successful.(bool), err
}

func (service *ConnectionService) MakeConnectionRequest(requestSenderId string, requestReceiverId string) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	canMakeConnection, err := service.canMakeConnection(requestSenderId, requestReceiverId)

	if !canMakeConnection {
		return false, err
	}

	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match(u1:User{id: $requestSenderId}) match(u2:User{id: $requestReceiverId})"+
				" create (u1)-[c1:CONNECTS{isApproved: false}]->(u2)"+
				" create (u1)<-[c2:CONNECTS{isApproved: false}]-(u2)"+
				" create (u1)-[c3:CONNECTS{isBlocked: false}]->(u2)"+
				" create (u1)<-[c4:CONNECTS{isBlocked: false}]-(u2)"+
				" return c1 is not null",
			map[string]interface{}{"requestSenderId": requestSenderId, "requestReceiverId": requestReceiverId})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return successful.(bool), err
}

func (service *ConnectionService) ApproveConnectionRequest(requestSenderId string, requestReceiverId string) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match(u1:User) "+
				"match(u2:User)"+
				"match((u1)-[c1:CONNECTS]->(u2)) "+
				"match((u2)-[c2:CONNECTS]->(u1)) "+
				"where u1.id = $requestSenderId and u2.id = $requestReceiverId "+
				"set c1.isApproved = true "+
				"set c2.isApproved = true "+
				"return c1 is not null",
			map[string]interface{}{"requestSenderId": requestSenderId, "requestReceiverId": requestReceiverId})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return false, result.Err()
	})
	return successful.(bool), err
}

type Users struct {
	users []string
}

func (service *ConnectionService) GetConnectionsUsernamesFor(userId string) ([]string, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	connections, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(
			"match(u1:User) "+
				"match(u2:User) "+
				"match((u1)-[c1:CONNECTS]->(u2)) "+
				"match((u1)-[c2:CONNECTS]->(u2)) "+
				"where u1.id = $userId and c1.isApproved and not c2.isBlocked "+
				"return distinct u2.id",
			map[string]interface{}{"userId": userId})

		users := Users{users: []string{}}
		if records == nil {
			return Users{users: []string{}}, nil
		}

		for records.Next() {
			record := records.Record()
			userId, _ := record.Get("u2.id")
			users.users = append(users.users, userId.(string))
		}
		return users, err
	})

	return connections.(Users).users, err
}

func (service *ConnectionService) GetRequestsUsernamesFor(userId string) ([]string, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	connections, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(
			"match(u1:User) "+
				"match(u2:User) "+
				"match((u1)-[c1:CONNECTS]->(u2)) "+
				"match((u1)-[c2:CONNECTS]->(u2)) "+
				"where u1.id = $userId and not c1.isApproved and not c2.isBlocked "+
				"return u2.id",
			map[string]interface{}{"userId": userId})

		users := Users{users: []string{}}
		if records == nil {
			return users, nil
		}

		for records.Next() {
			record := records.Record()
			userId, _ := record.Get("u2.id")
			users.users = append(users.users, userId.(string))
		}
		return users, err
	})

	return connections.(Users).users, err
}

func (service *ConnectionService) canMakeConnection(requestSenderId string, requestReceiverId string) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	success, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"optional match(u1:User{id: $requestSenderId}) optional match(u2:User{id: $requestReceiverId}) optional match((u1)-[c1:CONNECTS]->(u2)) return u1 is not null and u2 is not null and c1 is null",
			map[string]interface{}{"requestSenderId": requestSenderId, "requestReceiverId": requestReceiverId})
		if err != nil {
			return false, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return false, result.Err()
	})
	return success.(bool), err
}

func (service *ConnectionService) BlockUser(requestSenderId string, blockedUserId string) (bool, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	successful, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match(u1:User) "+
				"match(u2:User)"+
				"match((u1)-[c1:CONNECTS]->(u2)) "+
				"match((u2)-[c2:CONNECTS]->(u1)) "+
				"match((u1)-[c3:CONNECTS]->(u2)) "+
				"match((u2)-[c4:CONNECTS]->(u1)) "+
				"where u1.id = $requestSenderId and u2.id = $blockedUserId "+
				"set c1.isApproved = false "+
				"set c2.isApproved = false "+
				"set c3.isBlocked = true "+
				"set c4.isBlocked = true "+
				"return c1 is not null",
			map[string]interface{}{"requestSenderId": requestSenderId, "blockedUserId": blockedUserId})
		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0], nil
		}

		return false, result.Err()
	})
	return successful.(bool), err
}

func (service *ConnectionService) GetBlockedConnectionsUsernames(userId string) ([]string, error) {
	session := service.databaseDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	connections, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(
			"match(u1:User) "+
				"match(u2:User) "+
				"match((u1)-[c1:CONNECTS]->(u2)) "+
				"where u1.id = $userId and c1.isBlocked "+
				"return distinct u2.id",
			map[string]interface{}{"userId": userId})

		users := Users{users: []string{}}
		if records == nil {
			return Users{users: []string{}}, nil
		}

		for records.Next() {
			record := records.Record()
			userId, _ := record.Get("u2.id")
			users.users = append(users.users, userId.(string))
		}
		return users, err
	})

	return connections.(Users).users, err
}
