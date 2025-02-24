package serve

import (
	"fmt"
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"github.com/shewitt93/explore_service/pkg/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "pkg",
	Short: "GRPC server",
	Run:   startGrpcServer,
}

func init() {
	GrpcServerCmd.Flags().String("config", "config.yml", "config file (default is config.yml)")

}

func startGrpcServer(cmd *cobra.Command, args []string) {

	s := grpc.NewServer()

	grpcServer := server.NewExploreGRPCServer()
	grpclibs.RegisterExploreServiceServer(s, grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// If the server stops, that send the kill signal
		// The containers are set to auto-restart, so this will restart the container
		defer func() {
			log.Println("Sending kill signal")

			signalChan <- syscall.SIGTERM
		}()

		log.Printf("Starting server on port %s\n", os.Getenv("GRPC_PORT"))

		if err := s.Serve(listener); err != nil {
			log.Printf("Failed to serve: %v\n", err)
		}
	}()

	// Wait for the signal
	<-signalChan

	// This essentially stops the server, but only after all current requests have been completed
	s.GracefulStop()

	log.Println("Server stopped")
}
