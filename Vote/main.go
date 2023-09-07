package main

import (
	"flag"
	"fmt"

	"Assignment2/api"

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

	v2 := api.NewVoteApi()

	//vote API

	r.POST("/vote", v2.PostVoter)
	r.GET("/vote", v2.GetVoterListJson)
	r.DELETE("/vote/:id", v2.DeleteVoter)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
