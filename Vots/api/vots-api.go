package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Assignment2/schema"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nitishm/go-rejson/v4"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type VoteAPI struct {
	cache
	pubAPIURL string
	apiClient *resty.Client
}

func NewVoteAPI(location string, pubAPIurl string) (*VoteAPI, error) {

	apiClient := resty.New()
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
	return &VoteAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
		pubAPIURL: pubAPIurl,
		apiClient: apiClient,
	}, nil
}

func (p *VoteAPI) GetVote(c *gin.Context) {

	voteid := c.Param("id")
	if voteid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No publication ID provided"})
		return
	}

	cacheKey := "votes:" + voteid
	pubBytes, err := p.helper.JSONGet(cacheKey, ".")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find publication in cache with id=" + cacheKey})
		return
	}

	var vote schema.Vote
	err = json.Unmarshal(pubBytes.([]byte), &vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	c.JSON(http.StatusOK, vote)
}

func (p *VoteAPI) GetVotes(c *gin.Context) {

	var pubList []schema.Vote
	var pubItem schema.Vote

	//Lets query redis for all of the items
	pattern := "votes:*"
	ks, _ := p.client.Keys(p.context, pattern).Result()
	for _, key := range ks {
		err := p.getItemFromRedis(key, &pubItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find publication in cache with id=" + key})
			return
		}
		pubList = append(pubList, pubItem)
	}

	c.JSON(http.StatusOK, pubList)
}

func (r *VoteAPI) GetPollfromVotes(c *gin.Context) {
	voteId := c.Param("id")
	if voteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "votes:" + voteId
	var vote schema.Vote
	err := r.getItemFromRedis(cacheKey, &vote)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	// Here, we retrieve the PollID from the Vote struct
	pollID := vote.PollID

	fmt.Printf("%d", pollID)

	// Construct the URL to fetch the Poll using the PollID
	pubURL := r.pubAPIURL + fmt.Sprintf("%d", pollID)

	fmt.Printf("%s", pubURL)

	c.Redirect(http.StatusMovedPermanently, pubURL)
}

func (p *VoteAPI) DeleteVote(c *gin.Context) {
	pollid := c.Param("id")
	if pollid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	cacheKey := "votes:" + pollid

	// Check if voter exists
	exists, err := p.client.Exists(p.context, cacheKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking voter existence"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found in cache"})
		return
	}

	// Delete the voter from the cache
	err = p.client.Del(p.context, cacheKey).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete poll from cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "poll deleted successfully"})
}

// Helper to return a ToDoItem from redis provided a key
func (p *VoteAPI) getItemFromRedis(key string, pub *schema.Vote) error {

	//Lets query redis for the item, note we can return parts of the
	//json structure, the second parameter "." means return the entire
	//json structure
	itemObject, err := p.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	//JSONGet returns an "any" object, or empty interface,
	//we need to convert it to a byte array, which is the
	//underlying type of the object, then we can unmarshal
	//it into our ToDoItem struct
	err = json.Unmarshal(itemObject.([]byte), pub)
	if err != nil {
		return err
	}

	return nil
}
