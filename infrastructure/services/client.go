package services

import (
	profile "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ProfilesClient(address string) profile.ProfileServiceClient {
	prof, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to profile service: %v", err)
	}
	return profile.NewProfileServiceClient(prof)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
