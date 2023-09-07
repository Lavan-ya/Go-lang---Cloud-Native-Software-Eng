package poll

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type pollOption struct {
	PollOptionID    uint
	PollOptionText string
}

type Poll struct {
	PollID       uint
	PollTitle    string
	PollQuestion string
	PollOptions  []pollOption
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

type PollList struct {
	cache
}

func NewPoll() (*PollList,error) {
	redisUrl := os.Getenv("REDIS_URL")
    if redisUrl == "" {
        redisUrl = RedisDefaultLocation
    }
    return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*PollList, error) {
    
    client := redis.NewClient(&redis.Options{
        Addr: location,
    })

    ctx := context.Background()

  
    err := client.Ping(ctx).Err()
    if err != nil {
        log.Println("Error connecting to redis" + err.Error())
        return nil, err
    }

    
    jsonHelper := rejson.NewReJSONHandler()
    jsonHelper.SetGoRedisClientWithContext(ctx, client)

  
    return &PollList{
        cache: cache{
            cacheClient: client,
            jsonHelper:  jsonHelper,
            context:     ctx,
        },
    }, nil
}

func isRedisNilError(err error) bool {
    return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func redisKeyFromId(id uint64) string {
    return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *PollList) getItemFromRedis(key string, item *Poll) error {
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

func (v *PollList) AddItem(item Poll) error {
    redisKey := redisKeyFromId(uint64(item.PollID))
    var existingItem Poll
    if err := v.getItemFromRedis(redisKey, &existingItem); err == nil {
        return errors.New("item already exist")
    }
    if _, err := v.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
        return err
    }
    return nil
}

func(v *PollList) GetFullItem() ([]Poll,error){
    var toDoList []Poll
    var toDoItem Poll

    pattern := RedisKeyPrefix + "*"
    ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
    for _, key := range ks {
        err := v.getItemFromRedis(key, &toDoItem)
        if err != nil {
            return []Poll{},nil
        }
        toDoList = append(toDoList, toDoItem)
    }

    return toDoList,nil

}

func (v *PollList) DeleteItem(id uint64) error {

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



