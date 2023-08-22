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
	//v := voter.NewVoter(1, "John", "Doe")
	//v.AddPoll(1)
	//v.AddPoll(2)
	//v.AddPoll(3)
	//v.AddPoll(4)

	//b, _ := json.Marshal(v)
	//fmt.Println(string(b))
	vl := api.NewVoterApi()
	/*vl.AddVoter(1, "John", "Doe")
	vl.AddPoll(1, 1)
	vl.AddPoll(1, 2)
	vl.AddVoter(2, "Jane", "Alex")
	vl.AddPoll(2, 1)
	vl.AddPoll(2, 4)*/

	r.POST("/voters", vl.PostVoter)
	r.GET("/voters", vl.GetVoterListJson)
	r.GET("/voters/:id", vl.GetVoterJson)
	r.GET("/voters/:id/polls", vl.GetVoterHistory)
	r.GET("voters/:id/polls/:pollid", vl.GetVoterPoolid)
	r.POST("/voters/:id/polls/", vl.InsertPoll)
	r.GET("voters/health", vl.HealthCheck)
	r.PUT("/voters/:id", vl.UpdateVoter)
	r.DELETE("/voters/:id", vl.DeleteVoter)
	//r.GET("/voters", vl.ListAllVoter)

	//fmt.Println("------------------------")
	//fmt.Println(vl.GetVoterJson(1))
	/*fmt.Println("------------------------")
	fmt.Println(vl.GetVoterJson(2))
	fmt.Println("------------------------")
	fmt.Println(vl.GetVoterListJson())
	fmt.Println("------------------------")*/

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
