package main

import (
	"GRPCExample/example"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"net"
	"time"
)

type UserServer struct {

}

func (u * UserServer)Login(ctx context.Context, user *example.User) (*example.LoginResponse, error) {
	fmt.Println(user)
	return &example.LoginResponse{Success: true}, nil
}

func main() {
	user := &example.User{
		Username:             "Jason",
		Password:             "123456",
	}

	// Started server
	go Server()

	time.Sleep(1*time.Second)

	// Create our client
	client := Client()

	response, err := client.Login(context.Background(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}


// Creating GRPC Client
func Client() example.UserServiceClient {
	cc, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return example.NewUserServiceClient(cc)
}

// Setting up server
func Server() {
	lis, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}

	selfDefinedServer := &UserServer{}
	grpcServer := grpc.NewServer()
	example.RegisterUserServiceServer(grpcServer, selfDefinedServer)

	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}

func Hash(message proto.Message) [64]byte {
	messageMarshalled, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	hash := sha3.Sum512(messageMarshalled)

	return hash
}