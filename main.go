package main

import (
	"GRPCExample/example"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

type UserServer struct {
	mu sync.Mutex
	registry map[string]string
	loggedIn map[string]bool
	counter map[int32]int32
}

func (U *UserServer)Register(ctx context.Context, user *example.RegisterRequest) (*example.RegisterResponse, error){
	U.mu.Lock()
	defer U.mu.Unlock()
	U.registry[user.Username] = user.Password
	fmt.Println(user," has been registered")
	return &example.RegisterResponse{Success:true}, nil
}

func (U *UserServer)Login(ctx context.Context, user *example.LoginRequest) (*example.LoginResponse, error){
	U.mu.Lock()
	defer U.mu.Unlock()
	if _, ok := U.registry[user.Username]; ok{
		U.loggedIn[user.Username] = true
	}

	return &example.LoginResponse{Success:U.loggedIn[user.Username]},nil
}

func (U *UserServer) DoAction(ctx context.Context, in *example.DoActionRequest) (*example.DoActionResponse, error) {
	U.mu.Lock()
	defer U.mu.Unlock()
	if U.loggedIn[in.Username] {
		U.counter[in.Counter] += in.Number
	}
	return &example.DoActionResponse{}, nil
}



func main() {
	// Maps and slices should be initialized for it to be used
	//var m map[int]int
	//m = make(map[int]int)
	userLogin := &example.LoginRequest{
		Username: "Test",
		Password: "12345",
	}

	doAction := &example.DoActionRequest{
		Username: "Test",
		Number: 10,
		Counter : 2,
	}

	selfDefinedServer := &UserServer{
		registry: make(map[string]string),
		loggedIn: make(map[string]bool),
		counter: make(map[int32]int32),
	}

	// Started server
	go Server(selfDefinedServer)

	time.Sleep(1*time.Second)

	// Create our client
	client := Client()

	///////////////////////////////////
	// Stress test
	// 1000 Register requests at the same time
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i :=0; i < 1000; i++ {
		go func(i int) {
			defer wg.Done()
			userRegister := &example.RegisterRequest{
				Username: strconv.Itoa(i),
				Password: "12345",
			}

			_, err := client.Register(context.Background(), userRegister)
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()
	///////////////////////////////////

	response2, err := client.Login(context.Background(), userLogin)
	if err != nil {
		panic(err)
	}

	response3, err := client.DoAction(context.Background(), doAction)
	fmt.Println(response2)
	fmt.Println(response3)
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
func Server(userServer *UserServer) {
	lis, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	example.RegisterUserServiceServer(grpcServer, userServer)

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