package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/sainiajay/backend.ajaysaini.dev/services/bot"

	"google.golang.org/grpc"
)

type server struct {
	bot.UnimplementedBotServiceServer
}

func (s *server) HandleUserMessage(ctx context.Context, message *bot.Message) (*bot.Message, error) {
	body := message.GetBody()
	log.Printf("Received: %v", body)
	return &bot.Message{Body: "ACK " + message.GetBody()}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	log.Printf("Server listening on port %s", port)
	grpc_server := grpc.NewServer()
	bot.RegisterBotServiceServer(grpc_server, &server{})
	if err := grpc_server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
