package middleware

import (
	"context"
	"fmt"
	sessionclient "task_service/src/internal/adaptors/grpcclient"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func SessionAuthInterceptor(sessionClient *sessionclient.Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Assuming session_id is passed in metadata with key "session-id"
		sessionIDs := md.Get("session-id")
		if len(sessionIDs) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "session-id is not provided")
		}

		valid, userID, err := sessionClient.ValidateSession(ctx, sessionIDs[0])
		fmt.Println("userid : ", userID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "session validation failed: %v", err)
		}

		if !valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid session")
		}

		// You can add userID to context if needed here for handlers

		return handler(ctx, req)
	}
}
