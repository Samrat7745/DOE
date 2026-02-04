package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/Samrat/collector/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMetricCollectorServer
}

// StreamMetrics implements the gRPC service
func (s *server) StreamMetrics(stream pb.MetricCollector_StreamMetricsServer) error {
	for {
		// 1. Receive the next message from the stream
		req, err := stream.Recv()

		// 2. Check if the agent closed the connection
		if err == io.EOF {
			return stream.SendAndClose(&pb.MetricResponse{
				Success: true,
				Message: "Stream closed gracefully",
			})
		}
		if err != nil {
			return err
		}

		// 3. Process the data (For now, we just log it)
		fmt.Printf("[%s] CPU: %.2f%% | Mem: %d/%d KB\n",
			req.AgentId, req.CpuUsage, req.MemAvail, req.MemTotal)

		// ENGINEERING TIP: In a real app, you would push 'req'
		// into a Worker Pool or a Message Queue here.
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMetricCollectorServer(s, &server{})

	fmt.Println("Collector listening on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
