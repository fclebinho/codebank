package services

import (
	"context"
	"fmt"

	"github.com/fclebinho/codebank/src/data/grpc/pb"
	"github.com/fclebinho/codebank/src/domain/dtos"
	"github.com/fclebinho/codebank/src/domain/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionServer struct {
	pb.UnimplementedPaymentServiceServer
	ProcessTranscationUsecase usecases.ProcessTranscationUsecase
}

func (s *TransactionServer) Payment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	transactionDto := dtos.Transaction{
		Name:            in.GetCreditCard().GetName(),
		Number:          in.CreditCard.GetNumber(),
		ExpirationMonth: in.GetCreditCard().GetExpirationMonth(),
		ExpirationYear:  in.GetCreditCard().GetExpirationYear(),
		CVV:             in.GetCreditCard().GetCvv(),
		Amount:          in.GetAmount(),
		Store:           in.GetStore(),
		Description:     in.GetDescription(),
	}

	transaction, err := s.ProcessTranscationUsecase.Process(transactionDto)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if transaction.Status != "approved" {
		return nil, status.Error(codes.FailedPrecondition, "transaction rejected by the bank")
	}

	return &pb.PaymentResponse{Message: fmt.Sprintf("Your payment with card number %s is processing!", in.GetCreditCard().GetNumber())}, nil
}
