package server

import (
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"time"
)

type ExploreGRPCServer struct {
	grpclibs.ExploreServiceServer
}

func NewExploreGRPCServer() *ExploreGRPCServer {
	return &ExploreGRPCServer{}
}

type UserDecision struct {
	ID          uint64 `json:"id"`
	ActorID     uint64
	RecipientID uint64
	LikedID     bool
	CreatedAt   time.Time
}
