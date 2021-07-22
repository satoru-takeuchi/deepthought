package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/satoru-takeuchi/deepthought/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	err := subMain()
	if err != nil {
		log.Fatal(err)
	}
}

func handleBoot(wg *sync.WaitGroup, stream deepthought.Compute_BootClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("no more data from boot")
			} else if status.Code(err) == codes.Canceled {
				fmt.Println("boot canceled")
			} else {
				fmt.Printf("receiving boot response: %q\n", err)
			}
			break
		}
		fmt.Printf("Boot: %s\n", resp.Message)
	}
	wg.Done()
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
		time.Sleep(5 * time.Second)
		cancel()
	}(cancel)

	stream, err := cc.Boot(ctx, &deepthought.BootRequest{Silent: true})
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go handleBoot(&wg, stream)

	for {
		queries := []string{"foo", "Life", "Universe", "Everyting"}
		q := queries[rand.Intn(len(queries))]
		n := (rand.Intn(5) + 5) * 100
		fmt.Printf("query %q with deadline(now + %d ms)\n", q, n)
		ctx2, cancel := context.WithDeadline(ctx, time.Now().Add(time.Duration(n)*time.Millisecond))
		defer cancel()

		resp, err := cc.Infer(ctx2, &deepthought.InferRequest{Query: q})
		if err != nil {
			code := status.Code(err)
			if code == codes.InvalidArgument || code == codes.DeadlineExceeded {
				fmt.Printf("couldn't get answer: %q\n", status.Convert(err).Message())
				continue
			} else if code == codes.Canceled {
				fmt.Println("finished query")
				return nil
			}
			return err
		}
		fmt.Printf("answer to %q is %d\n", q, resp.Answer)
		time.Sleep(500 * time.Millisecond)
	}
}
