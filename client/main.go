package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"time"

	"github.com/sainiajay/backend.ajaysaini.dev/services/bot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(host, opts...)
}

const (
	remote = "api.ajaysaini.dev:443"
	local  = ":9000"
)

func main() {
	address := flag.String("address", remote, "Address where service is running.")
	flag.Parse()

	conn, err := NewConn(*address, false)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := bot.NewBotServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := client.HandleUserMessage(ctx, &bot.Message{Body: "Hi"})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Got: %s", resp.GetBody())
}
