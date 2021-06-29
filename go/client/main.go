package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/satoru-takeuchi/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	err := subMain()
	if err != nil {
		log.Fatal(err)
	}
}

func subMain() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: %s <host:port>", os.Args[0])
	}
	addr := os.Args[1]
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	cc := deepthought.NewComputeClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel func()) {
		time.Sleep(2500 * time.Millisecond)
		cancel()
	}(cancel)
	stream, err := cc.Boot(ctx, &deepthought.BootRequest{})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				break
			}
			return fmt.Errorf("receiving boot response: %w", err)
		}
		fmt.Printf("Boot: %s\n", resp.Message)
	}
	return nil
}
