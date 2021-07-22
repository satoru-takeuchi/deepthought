package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/satoru-takeuchi/deepthought/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var portNumber = 13333

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	serv := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kep))
	deepthought.RegisterComputeServer(serv, &Server{})
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatal(err)
	}
	serv.Serve(l)
}
