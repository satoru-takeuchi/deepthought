package main

import (
	"context"
	"time"

	"github.com/satoru-takeuchi/deepthought/go/deepthought"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	deepthought.UnimplementedComputeServer
}

var _ deepthought.ComputeServer = &Server{}

func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	if req.Silent {
		return nil
	}
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(1 * time.Second):
		}
		err := stream.Send(&deepthought.BootResponse{
			Message: "I THINK THEREFORE I AM.",
		})
		if err != nil {
			return err
		}
	}
}

func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}
	time.Sleep(41 * time.Second)
	return &deepthought.InferResponse{
		Answer: 42,
	}, nil
}
