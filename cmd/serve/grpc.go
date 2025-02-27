package serve

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shewitt93/explore_service/internal/database"
	"github.com/shewitt93/explore_service/internal/repository"
	"github.com/shewitt93/explore_service/internal/server"
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "GRPC server",
	Run:   startGrpcServer,
}

func init() {
	GrpcServerCmd.Flags().String("config", "config.yml", "config file (default is config.yml)")

}

func startGrpcServer(cmd *cobra.Command, args []string) {

	log.Println("starting grpc server")
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	decisionRepository := repository.NewDecisionRepositoryImpl(db)

	s := grpc.NewServer()

	grpcServer := server.NewExploreGRPCServer(decisionRepository)
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

func initDB() (*sql.DB, error) {
	// Get database connection details from environment variables
	dbConfig := database.ConfigDatabase{
		User:     getEnvWithDefault("DB_USER", "test"),
		Password: getEnvWithDefault("DB_PASS", "test"),
		Host:     getEnvWithDefault("DB_HOST", "mysqldb"),
		Port:     getEnvWithDefault("DB_PORT", "3306"),
	}

	dbName := getEnvWithDefault("DB_NAME", "explore_service")

	// Generate DSN
	dsn := database.GenerateDSN(dbConfig, dbName)

	// Create database connection
	return database.NewMysqlConnection(dsn)
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
