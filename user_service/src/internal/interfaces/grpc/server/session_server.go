package server

import (
	"context"
	"fmt"

	pb "user_service/src/internal/interfaces/grpc/generated"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SessionServer struct {
	pb.UnimplementedSessionValidatorServer
	redisClient *redis.Client
}

func NewSessionServer(redisClient *redis.Client) *SessionServer {
	return &SessionServer{redisClient: redisClient}
}

func (s *SessionServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	// userID, err := s.redisClient.Get(ctx, req.SessionId).Result()
	userID, err := s.redisClient.Get(ctx, "session:"+req.SessionId).Result()

	fmt.Println("user req session id : ", req.SessionId)

	// fmt.Println()
	if err == redis.Nil {
		return &pb.ValidateSessionResponse{Valid: false, Error: "Session not found"}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "Redis error: %v", err)
	}

	return &pb.ValidateSessionResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
