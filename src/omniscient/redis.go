package omniscient

import (
	"fmt"
	"log"
	"time"

	"backoff"

	"gopkg.in/redis.v3"
)

// RedisClient is an interface which can interfact with a redis server.
type RedisClient interface {
	Delete(keys ...string) (int64, error)
	HGetAllMap(key string) (map[string]string, error)
	HMSet(key, field, value string, pairs ...string) (string, error)
	LPush(key string, values ...string) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
	LRem(key string, count int64, value interface{}) (int64, error)
	Ping() (string, error)
	Set(key string, value interface{}, expiration time.Duration) (string, error)
}

type redisClient struct {
	client *redis.Client
}

var _ RedisClient = (*redisClient)(nil)

// NewRedisClient creates an instance of RedisClient. It will attempt multiple times
// using a backoff policy.
func NewRedisClient(addr string) (RedisClient, error) {
	var client *redis.Client
	var attempt int

	for {
		sleepTime, err := backoff.DefaultPolicy.Duration(attempt)
		if err != nil {
			return nil, err
		}

		time.Sleep(sleepTime)

		log.Printf("connecting to redis server %s", addr)
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		})

		_, err = client.Ping().Result()
		if err == nil {
			break
		}

		attempt++
		log.Printf("backing off because redis didn't respond to ping: %s", err)

	}
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("unable to ping redis server: %v", err)
	}

	return &redisClient{
		client: client,
	}, nil
}

func (rc *redisClient) Delete(keys ...string) (int64, error) {
	cmd := rc.client.Del(keys...)
	return cmd.Result()
}

func (rc *redisClient) HGetAllMap(key string) (map[string]string, error) {
	cmd := rc.client.HGetAllMap(key)
	return cmd.Result()
}

func (rc *redisClient) HMSet(key, field, value string, pairs ...string) (string, error) {
	cmd := rc.client.HMSet(key, field, value, pairs...)
	return cmd.Result()
}

func (rc *redisClient) LPush(key string, values ...string) (int64, error) {
	cmd := rc.client.LPush(key, values...)
	return cmd.Result()
}

func (rc *redisClient) LRange(key string, start, stop int64) ([]string, error) {
	cmd := rc.client.LRange(key, start, stop)
	return cmd.Result()
}

func (rc *redisClient) LRem(key string, count int64, value interface{}) (int64, error) {
	cmd := rc.client.LRem(key, count, value)
	return cmd.Result()
}

func (rc *redisClient) Ping() (string, error) {
	cmd := rc.client.Ping()
	return cmd.Result()
}

func (rc *redisClient) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	cmd := rc.client.Set(key, value, expiration)
	return cmd.Result()
}
