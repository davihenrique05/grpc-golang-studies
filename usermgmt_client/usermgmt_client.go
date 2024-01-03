package main

import (
	"context"
	"log"
	"time"

	pb "github.com/davihenrique05/grpc-studies/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var newUsers = make(map[string]int32)
	newUsers["Alice"] = 43
	newUsers["Bob"] = 30

	for name, age := range newUsers {
		response, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})

		if err != nil {
			log.Fatalf("Could not create user: %v", err)
		}

		log.Printf("User Details: \n NAME: %s\n AGE: %d \n ID: %d", response.GetName(), response.GetAge(), response.GetId())
	}
	params := &pb.GetUsersParams{}
	response, err := c.GetUsers(ctx, params)

	if err != nil {
		log.Fatalf("Could not retrieve users %v", err)
	}

	log.Printf("\n USERLIST \n%v", response.GetUsers())
}
