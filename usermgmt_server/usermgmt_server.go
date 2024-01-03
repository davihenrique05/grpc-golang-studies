package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/davihenrique05/grpc-studies/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func getOrCreateUsersFile() *pb.UserList {
	userList := &pb.UserList{}
	readBytes, err := os.ReadFile("users.json")
	var data []byte

	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			if err := os.WriteFile("users.json", data, 0664); err != nil {
				log.Fatalf("Failed while writting to file: %v", err)
			}
			return userList
		}
	}

	if err := protojson.Unmarshal(readBytes, userList); err != nil {
		log.Fatalf("Failed to parse user list: %v", err)
	}
	return userList
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, input *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", input.GetName())
	userList := getOrCreateUsersFile()

	userId := int32(rand.Intn(1000))
	createdUser := &pb.User{Id: userId, Name: input.GetName(), Age: input.GetAge()}

	userList.Users = append(userList.Users, createdUser)

	jsonBytes, err := protojson.Marshal(userList)

	if err != nil {
		log.Fatalf("JSON Marshalling failed: %v", err)
	}

	if err := os.WriteFile("users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("Failed while writting to file: %v", err)
	}

	return createdUser, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, input *pb.GetUsersParams) (*pb.UserList, error) {
	usersList := getOrCreateUsersFile()
	return usersList, nil
}

func main() {
	var userMgmtServer *UserManagementServer = NewUserManagementServer()

	if err := userMgmtServer.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
