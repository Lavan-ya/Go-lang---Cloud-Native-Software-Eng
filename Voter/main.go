package main

import (
	"Assignment2/api"
	"flag"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	vl := api.NewVoterApi()
	
	r.POST("/voters", vl.PostVoter)
	r.GET("/voters", vl.GetVoterListJson)
	r.GET("/voters/:id", vl.GetVoterJson)
	r.GET("/voters/:id/polls", vl.GetVoterHistory)
	r.GET("voters/:id/polls/:pollid", vl.GetVoterPoolid)
	r.POST("/voters/:id/polls/", vl.InsertPoll)
	r.GET("voters/health", vl.HealthCheck)
	r.PUT("/voters/:id", vl.UpdateVoter)
	r.DELETE("/voters/:id", vl.DeleteVoter)


	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
