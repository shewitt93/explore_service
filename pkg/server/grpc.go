package server

import (
	"github.com/spf13/explore_service/pkg/grpclibs"
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
	ActorID     uint64 `json:"actor_id"`
	RecipientID uint64 `json:"recipient_id"`
	LikedID     bool
	CreatedAt   time.Time `json:"created_at"`
}
