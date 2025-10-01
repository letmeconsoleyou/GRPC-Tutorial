package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/letmeconsoleyou/grpc/grpc-chat/chatpb"
	"google.golang.org/grpc"
)

var registeredUsers []string

type server struct {
	pb.UnimplementedChatServer
}

type Message struct {
	username  string
	content   string
	timestamp int64
}

var messages []Message

// Unary
func (s *server) RegisterUser(ctx context.Context, req *pb.UserInfo) (*pb.RegisterResponse, error) {
	registeredUsers = append(registeredUsers, req.GetUsername())
	return &pb.RegisterResponse{Message: req.GetUsername() + " registered successfully!"}, nil
}

// Client Streaming
func (s *server) SendOffineMessages(stream pb.Chat_SendOffineMessagesServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// return
			return stream.SendAndClose(&pb.UploadSummary{SuccessfullyStoredMessage: int64(len(messages))})
		}
		if err != nil {
			return err
		}
		fmt.Println("Received message")
		messages = append(messages, Message{
			username:  req.GetUsername(),
			content:   req.GetContent(),
			timestamp: req.GetTimestamp(),
		})
	}
}

// Server streaming
func (s *server) GetChatHistory(req *pb.ChatHistoryRequest, stream pb.Chat_GetChatHistoryServer) error {
	for idx, message := range messages {
		if idx >= int(req.GetTimestamp())-1 {
			err := stream.Send(&pb.ChatMessage{Username: message.username, Content: message.content, Timestamp: message.timestamp})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *server) Chat(stream pb.Chat_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received %s, %s, %d", req.GetUsername(), req.GetContent(), req.GetTimestamp())
		err = stream.Send(&pb.ChatMessage{Username: req.GetUsername() + " RESPOSNE", Content: req.GetContent() + " RESPONSE", Timestamp: req.GetTimestamp()})
		if err != nil {
			log.Fatalf("Could not send message from server via stream: %v", err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Error creating tcp listener: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})

	fmt.Println("Server listening...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Could not server: %v", err)
	}
}
