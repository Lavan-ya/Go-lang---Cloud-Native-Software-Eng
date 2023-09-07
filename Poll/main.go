package main

import (
	pollapi "Assignment2/api"

	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
	cacheURL string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Cache URL")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func envVarOrDefault(envVar string, defaultVal string) string {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func setupParms() {
	//first process any command line flags
	processCmdLineFlags()

	//now process any environment variables
	cacheURL = envVarOrDefault("PUBAPI_CACHE_URL", cacheURL)
	hostFlag = envVarOrDefault("PUBAPI_HOST", hostFlag)
	pfNew, err := strconv.Atoi(envVarOrDefault("PUBAPI_PORT", fmt.Sprintf("%d", portFlag)))
	//only update the port if we were able to convert the env var to an int, else
	//we will use the default we got from the command line, or command line defaults
	if err == nil {
		portFlag = uint(pfNew)
	}

}

func main() {
	setupParms()
	apiHandler, err := pollapi.NewPubApi(cacheURL)

	if err != nil {
		panic(err)
	}

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	v1 := pollapi.NewPollApi()

	//poll API

	r.POST("/poll", v1.PostVoter)
	r.GET("/poll", v1.GetVoterListJson)
	r.DELETE("/poll/:id", v1.DeleteVoter)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
