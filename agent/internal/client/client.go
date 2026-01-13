package agentclient

import (
	"context"
	"log"
	agentpb "netshield/agent/proto"

	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	stream agentpb.AgentService_StreamMetricsClient
}

func New(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := agentpb.NewAgentServiceClient(conn)
	stream, err := c.StreamMetrics(context.Background())
	if err != nil {
		log.Print("[agent] failed to connect to server:", err)
	}
	cli := &Client{
		conn:   conn,
		stream: stream,
	}
	cli.listenControlAsync()
	return cli, nil
}

func (c *Client) listenControlAsync() {
	go func() {
		for {
			msg, err := c.stream.Recv()
			if err != nil {
				log.Println("[agent] control stream closed:", err)
				return
			}
			log.Printf("[agent] control message: type=%s data=%s\n", msg.Type, msg.Data)
		}
	}()
}

func (c *Client) ReportMetric(m *agentpb.NetworkMetric) error {
	return c.stream.Send(m)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
