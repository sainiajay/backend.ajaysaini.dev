package main

import (
	"ajaysaini-dev/services/bot"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"time"

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

func main() {
	address := flag.String("address", "apibackend-rvovnjllja-as.a.run.app:443", "Address where service is running.")
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
