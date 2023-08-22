package voter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type voterPoll struct {
	PollID   uint64
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint64
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "vote:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterList struct {
	cache
	//Voters map[uint64]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoter() (*VoterList, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterList, error) {
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
	return &VoterList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

// REDIS HELPERS
//------------------------------------------------------------

// We will use this later, you can ignore for now
func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func redisKeyFromId(id uint64) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *VoterList) getItemFromRedis(key string, item *Voter) error {
	itemObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}
	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

func (v *VoterList) PutVoter(VoterID uint64, item Voter) error {
	redisKey := redisKeyFromId(item.VoterID)
	var existingItem Voter
	if err := v.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("item does not exist")
	}
	existingItem.FirstName = item.FirstName
	existingItem.LastName = item.LastName
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", existingItem); err != nil {
		return err
	}
	/*
		list, ok := v.Voters[item.VoterID]
		if !ok {
			return errors.New("item doesn't exists")
		}
		list.FirstName = item.FirstName
		list.LastName = item.LastName
		v.Voters[item.VoterID] = list*/
	return nil
}

func (v *VoterList) AddVoterlist(item Voter) error {
	redisKey := redisKeyFromId(item.VoterID)
	var existingItem Voter
	if err := v.getItemFromRedis(redisKey, &existingItem); err == nil {
		return errors.New("item already exist")
	}
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
		return err
	}
	return nil
}

func (v *Voter) AddPoll(pollID uint64, voteDate time.Time) {

	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: voteDate})
}

func(v *VoterList) AddTopoll(voterID uint64, pollID uint64, voteDate time.Time) error {
	redisKey := redisKeyFromId(voterID)
	var existingItem Voter
	if err := v.getItemFromRedis(redisKey, &existingItem); err !=nil{
		return errors.New("item doesnot exist");
	}
	existingItem.AddPoll(pollID,voteDate);
	if _,err := v.jsonHelper.JSONSet(redisKey, ".", existingItem); err != nil {
		return err
	}
	return nil;
	
}

func (v *VoterList) DeleteItem(id uint64) error {

	pattern := redisKeyFromId(id)
	numDeleted, err := v.cacheClient.Del(v.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

func (v *VoterList) GetItem(id uint64) (Voter, error) {

	var item Voter
	pattern := redisKeyFromId(id)
	err := v.getItemFromRedis(pattern, &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}

func (v *VoterList) GetVoterHistoryItem(id uint64) ([]voterPoll,error) {
	var item Voter
	pattern := redisKeyFromId(id)
	err := v.getItemFromRedis(pattern, &item)
	if err != nil {
		return []voterPoll{},nil
	}
	return item.VoteHistory,nil
}

func (v *VoterList) GetVoterPoolidItem(id uint64,pollID uint64) ([]voterPoll,error) {
	var item Voter
	pattern := redisKeyFromId(id)
	err := v.getItemFromRedis(pattern, &item)
	if err != nil {
		return []voterPoll{},nil
	}
	VoteHistory := item.VoteHistory
	var selectedPollVoteHistory []voterPoll;
	for _, poll := range VoteHistory {
		if poll.PollID == pollID {
			selectedPollVoteHistory = append(selectedPollVoteHistory, poll)
		}
	}
	return selectedPollVoteHistory, nil
}
