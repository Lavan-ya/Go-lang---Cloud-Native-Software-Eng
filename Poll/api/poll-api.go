package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"Assignment2/poll"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type PollApi struct {
	pollList *poll.PollList
}

func NewPubAPI(location string) (*PubAPI, error) {

	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	//By default, redis manages keys and values, where the values
	//are either strings, sets, maps, etc.  Redis has an extension
	//module called ReJSON that allows us to store JSON objects
	//however, we need a companion library in order to work with it
	//Below we create an instance of the JSON helper and associate
	//it with our redis connnection
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &PubAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
	}, nil
}

func NewPollApi() *PollApi {
	dbHandler, err := poll.NewPoll()
	if err != nil {
		return nil
	}
	return &PollApi{pollList: dbHandler}
}

func (v *PollApi) PostVoter(c *gin.Context) {
	var list poll.Poll
	if err := c.ShouldBindJSON(&list); err != nil {
		fmt.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("Value is -------------------------", list)
	if err := v.pollList.AddItem(list); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (v *PollApi) GetVoterListJson(ctx *gin.Context) {
	voter, err := v.pollList.GetFullItem()
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(voter)
	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", b)
}

func (v *PollApi) DeleteVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := v.pollList.DeleteItem(VoterId); err != nil {
		log.Println("Error deleting item: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)

}
