package main

import (
	pollapi "Assignment2/poll-api"

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

	
	v1 :=pollapi.NewPollApi()
	
	
	

	//poll API

	r.POST("/poll",v1.PostVoter)
	r.GET("/poll",v1.GetVoterListJson)
	r.DELETE("/poll/:id",v1.DeleteVoter)


	

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
