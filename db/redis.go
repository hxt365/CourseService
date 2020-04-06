package db

import (
	"CourseService/utils"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/avast/retry-go"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

type RedisClient interface{}

type Redis struct {
	Client RedisClient
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("cache error: ", err)
	}
	return &Redis{Client: client}
}

func NewTestRedis() *Redis {
	client, err := miniredis.Run()
	if err != nil {
		log.Println("cannot create miniredis")
	}
	return &Redis{Client: client}
}

func (rd *Redis) GetCache(key string) (string, error) {
	var (
		val string
		err error
	)
	switch client := rd.Client.(type) {
	case *redis.Client:
		val, err = client.Get(key).Result()
	case *miniredis.Miniredis:
		val, err = client.Get(key)
	}
	return val, err
}

func (rd *Redis) SetCache(key string, value interface{}) {
	switch client := rd.Client.(type) {
	case *redis.Client:
		go func() {
			err := retry.Do(func() error {
				err := client.Set(key, value, 0).Err()
				return err
			}, retry.Attempts(utils.MAX_RETRY_TIMES), retry.Delay(2*time.Second))
			if err != nil {
				log.Println("cache error: cannot set cache")
			}
		}()
	case *miniredis.Miniredis:
		_ = client.Set(key, value.(string))
	}
}

func (rd *Redis) DelCache(keys ...string) {
	switch client := rd.Client.(type) {
	case *redis.Client:
		go func() {
			err := retry.Do(func() error {
				err := client.Del(keys...).Err()
				if err != nil {
					return err
				}
				return nil
			}, retry.Attempts(utils.MAX_RETRY_TIMES), retry.Delay(2*time.Second))
			if err != nil {
				log.Println("cache error: cannot delete cache")
			}
		}()
	case *miniredis.Miniredis:
		client.Del(keys[0])
	}
}

func (rd *Redis) SetCacheField(key string, values ...interface{}) {
	switch client := rd.Client.(type) {
	case *redis.Client:
		go func() {
			err := retry.Do(func() error {
				err := client.HSet(key, values...).Err()
				if err != nil {
					return err
				}
				return nil
			}, retry.Attempts(utils.MAX_RETRY_TIMES), retry.Delay(2*time.Second))
			if err != nil {
				log.Println("cache error: cannot set cache")
			}
		}()
	case *miniredis.Miniredis:
		client.HSet(key, utils.InterfacesToStrings(values...)...)
	}
}

func (rd *Redis) GetCacheField(key, field string) (string, error) {
	var (
		val string
		err error
	)

	switch client := rd.Client.(type) {
	case *redis.Client:
		val, err = client.HGet(key, field).Result()
	case *miniredis.Miniredis:
		val = client.HGet(key, field)
		if val == "" {
			err = fmt.Errorf("key does not exist")
		}
	}
	return val, err
}

func (rd *Redis) DelCacheField(key string, fields ...string) {
	switch client := rd.Client.(type) {
	case *redis.Client:
		go func() {
			err := retry.Do(func() error {
				err := client.HDel(key, fields...).Err()
				if err != nil {
					return err
				}
				return nil
			}, retry.Attempts(utils.MAX_RETRY_TIMES), retry.Delay(2*time.Second))
			if err != nil {
				log.Println("cache error: cannot delete cache")
			}
		}()
	case *miniredis.Miniredis:
		client.HDel(key, fields[0])
	}
}
