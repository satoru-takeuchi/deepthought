package main

import (
	"context"

	"github.com/satoru-takeuchi/deepthought"
)

type Server struct {
	deepthought.UnimplementedComputeServer
}

var _ deepthought.ComputeServer = &Server{}

func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	panic("not implemented yet")
}

func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	panic("not implemented yet")
}
