package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fclebinho/codebank/src/data/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *Server) Payment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	return &pb.PaymentResponse{Message: fmt.Sprintf("Your payment with card number %s is processing!", in.GetCreditCard().GetNumber())}, nil
}

func main() {
	println("Running gRPC server")

	listener, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &Server{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed ro serve: %v", err)
	}

}
