package grpc

import (
	"context"
	"time"

	pb "github.com/Hanulus/ap2-generated/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"payment-service/internal/usecase"
)

// PaymentGRPCServer implements the gRPC PaymentService interface
type PaymentGRPCServer struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentGRPCServer(uc *usecase.PaymentUseCase) *PaymentGRPCServer {
	return &PaymentGRPCServer{uc: uc}
}

// ProcessPayment handles incoming gRPC payment requests
func (s *PaymentGRPCServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	// Validate input
	if req.OrderId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id and amount are required")
	}

	// Call the use case (business logic stays here, untouched)
	payment, err := s.uc.Authorize(req.OrderId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return gRPC response
	return &pb.PaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     timestamppb.New(time.Now()),
	}, nil
}
