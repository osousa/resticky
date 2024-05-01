package main

import (
	"context"
	"log"

	"github.com/osousa/resticky/internal"
	"github.com/osousa/resticky/internal/protobuf/resticky"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		conn *grpc.ClientConn
		ctx  context.Context
		err  error
	)

	ctx = context.Background()

	// Connect to the server
	conn, err = grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	r := resticky.NewRestickyServiceClient(conn)

	z := zap.NewExample()

	c := internal.ProvideCLI(ctx, z, r)

	c.Start()
}
