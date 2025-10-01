package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/letmeconsoleyou/grpc/grpc-chat/chatpb"
)

type Message struct {
	username  string
	content   string
	timestamp int64
}

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not create a client: %v", err)
	}

	defer conn.Close()

	c := pb.NewChatClient(conn)
	// res, err := c.RegisterUser(context.Background(), &pb.UserInfo{Username: "letmeconsoleyou"})
	// if err != nil {
	// 	log.Fatalf("Could not register user: %v", err)
	// }
	// fmt.Printf("Response: %s", res.GetMessage())

	var messages []Message
	messages = append(messages, Message{username: "John", content: "Hi", timestamp: 1})
	messages = append(messages, Message{username: "Harry", content: "Hey", timestamp: 2})
	messages = append(messages, Message{username: "Messi", content: "Bye", timestamp: 3})

	// stream, err := c.SendOffineMessages(context.Background())
	// if err != nil {
	// 	log.Fatalf("could not open stream: %v", err)
	// }

	// for _, message := range messages {
	// 	err := stream.Send(&pb.ChatMessage{Username: message.username, Content: message.content, Timestamp: message.timestamp})
	// 	if err != nil {
	// 		log.Fatal("Could not send: %v", err)
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

	// res, err := stream.CloseAndRecv()
	// if err != nil {
	// 	log.Fatal("Could not receive: %v", err)
	// }
	// fmt.Printf("Response: %d\n", res.GetSuccessfullyStoredMessage())
	// chat_stream, err := c.GetChatHistory(context.Background(), &pb.ChatHistoryRequest{Timestamp: 3})
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	// for {
	// 	res, err := chat_stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}

	// 	fmt.Printf("Chat history: %s, %s, %d\n", res.GetUsername(), res.GetContent(), res.GetTimestamp())
	// }

	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Could not open stream channel: %v", err)
	}

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not get res: %v", err)
			}
			log.Printf("Received %s, %s, %d\n", res.GetUsername(), res.GetContent(), res.GetTimestamp())
		}
	}()

	for _, message := range messages {
		err := stream.Send(&pb.ChatMessage{Username: message.username, Content: message.content, Timestamp: message.timestamp})
		if err != nil {
			log.Fatalf("Could not send stream of msg from client: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
