package repository

import (
	"context"

	pb "github.com/Hanulus/ap2-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PaymentGRPCClient replaces the old REST client.
// It implements the PaymentClient interface from usecase/ports.go.
type PaymentGRPCClient struct {
	client pb.PaymentServiceClient
}

// NewPaymentGRPCClient connects to the Payment Service gRPC server
func NewPaymentGRPCClient(address string) (*PaymentGRPCClient, error) {
	// insecure.NewCredentials() means no TLS — fine for local/docker
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PaymentGRPCClient{client: pb.NewPaymentServiceClient(conn)}, nil
}

// Authorize calls ProcessPayment on the Payment Service via gRPC
func (c *PaymentGRPCClient) Authorize(orderID string, amount int64) (string, error) {
	resp, err := c.client.ProcessPayment(context.Background(), &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", err
	}

	// Return empty string if declined — use case will mark order as Failed
	if resp.Status != "Authorized" {
		return "", nil
	}
	return resp.TransactionId, nil
}
