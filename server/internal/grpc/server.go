package grpcserver

import (
	"context"
	"fmt"
	"io"
	"log"
	agentpb "netshield/agent/proto"
	"netshield/server/internal/db"
)

type AgentServiceServer struct {
	agentpb.UnimplementedAgentServiceServer
	store *db.Store
}

func NewAgentServiceServer(store *db.Store) *AgentServiceServer {
	return &AgentServiceServer{store: store}
}

// StreamMetrics: agents send metrics; we save to Postgres.
func (s *AgentServiceServer) StreamMetrics(stream agentpb.AgentService_StreamMetricsServer) error {
	log.Println("[server] StreamMetrics connected")
	defer log.Println("[server] StreamMetrics disconnected")

	for {
		m, err := stream.Recv()
		if err == io.EOF {
			
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream recv: %w", err)
		}

		log.Printf("[server] metric from device=%s user=%s domain=%s score=%d\n",
			m.DeviceId, m.UserId, m.Domain, m.ExperienceScore)

		if err := s.store.SaveMetric(context.Background(), m); err != nil {
			log.Println("[server] SaveMetric error:", err)
		}
	}
}
