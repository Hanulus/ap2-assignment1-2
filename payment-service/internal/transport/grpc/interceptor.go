package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// LoggingInterceptor logs every incoming gRPC request with method name and duration.
// This is a Unary Interceptor (for simple request-response calls like ProcessPayment).
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Call the actual handler
	resp, err := handler(ctx, req)

	duration := time.Since(start)
	fmt.Printf("[gRPC] method=%s duration=%s err=%v\n", info.FullMethod, duration, err)

	return resp, err
}
