package serve

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var RestServerCmd = &cobra.Command{
	Use:   "rest",
	Short: "Serves REST API",
	Run:   startRestServer,
}

func init() {
	RestServerCmd.Flags().String("config", "config.yaml", "config file")
}

func startRestServer(cmd *cobra.Command, args []string) {

	engine := gin.Default()
	group := engine.Group("/api")

	group.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	group.POST("/user/create")
	group.GET("/login")
	//group.GET("/discover" middleware.JWTAuthMiddleware())
	//group.POST("/swipe")

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: engine,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start the sever: %v", err)
	}
}
