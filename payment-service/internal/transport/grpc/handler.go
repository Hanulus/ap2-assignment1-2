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
	if req.OrderId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id and amount are required")
	}

	payment, err := s.uc.Authorize(req.OrderId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     timestamppb.New(time.Now()),
	}, nil
}

// ListPayments returns payments filtered by amount range
func (s *PaymentGRPCServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	payments, err := s.uc.ListByAmountRange(req.MinAmount, req.MaxAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert domain payments to protobuf responses
	var pbPayments []*pb.PaymentResponse
	for _, p := range payments {
		pbPayments = append(pbPayments, &pb.PaymentResponse{
			TransactionId: p.TransactionID,
			Status:        p.Status,
			CreatedAt:     timestamppb.New(time.Now()),
		})
	}

	return &pb.ListPaymentsResponse{Payments: pbPayments}, nil
}
