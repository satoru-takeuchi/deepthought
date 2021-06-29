package main

import (
	"fmt"
	"log"
	"net"

	"github.com/satoru-takeuchi/deepthought"
	"google.golang.org/grpc"
)

var portNumber = 13333

func main() {
	serv := grpc.NewServer()
	deepthought.RegisterComputeServer(serv, &Server{})
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatal(err)
	}
	serv.Serve(l)
}
