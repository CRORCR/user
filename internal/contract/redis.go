package contract

// redis 搭建

import (
	"fmt"
	"time"

	"github.com/CRORCR/user/internal/model"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client      *redis.Client
	RedisLocker *redislock.Client
}

var (
	RedisClient Redis
	RainKey     = "user"
)

/*
   host: "tx6-inno-frenzy-db-test01.bj"
   port: 6379
   password: "redisinkePassWd"
   max_active: 500
   idle_timeout: 1000
   connect_timeout: 1000
   read_timeout: 1000
   write_timeout: 1000
   db: 1
   retry: 2
*/
func CreateRedisConnection(config model.RedisConfig) *redis.Client {
	hostAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     hostAddr,
		Password: config.Password,
		DB:       int(config.DB),
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		panic(err)
	}

	if nil == redisClient {
		panic(fmt.Sprintf("InitRedis Error:%v", hostAddr))
	}

	return redisClient
}

func init() {
	RedisClient.Client = CreateRedisConnection(global.GlobalConfig.Conf.Redis.Host, global.GlobalConfig.Conf.Redis.Port, "", global.GlobalConfig.Conf.Redis.DB)
	RedisClient.RedisLocker = redislock.New(RedisClient.Client)
}

func (self *Redis) GetLock(key string, ttl time.Duration) (*redislock.Lock, error) {
	return self.RedisLocker.Obtain(key, ttl, nil)
}

func (self *Redis) ReleaseLock(lock *redislock.Lock) {
	if lock != nil {
		lock.Release()
	}
}

func (self *Redis) GetValue(key string, v interface{}) error {
	return self.Client.Get(key).Scan(v)
}

func (self *Redis) SetValue(key string, v interface{}, expiration time.Duration) error {
	return self.Client.Set(key, v, expiration).Err()
}

func (this *Redis) Set(key string, value string, expiration time.Duration) {
	log.Logger.Debug(fmt.Sprintf(`redis set. key: %s, val: %s, expiration: %v`, key, value, expiration))
	if err := this.Client.Set(key, value, expiration).Err(); err != nil {
		log.Logger.Trace(`redis set error`, err)
	}
}

func (this *Redis) Get(key string) string {
	log.Logger.Debug(fmt.Sprintf(`redis get. key: %s`, key))
	result, err := this.Client.Get(key).Result()
	if err != nil {
		if err.Error() == `redis: nil` {
			return ``
		}
		log.Logger.Trace(`redis get error`, err)
	}
	log.Logger.Debug(fmt.Sprintf(`redis get. result: %s`, result))
	return result
}

func (self *Redis) SetHash(key string, fields map[string]interface{}) error {
	if err := self.Client.HMSet(key, fields).Err(); err != nil {
		return err
	}
	self.Client.Expire(key, time.Second*30) //30秒
	return nil
}

func (self *Redis) GetHash(key string) (map[string]string, error) {
	return self.Client.HGetAll(key).Result()
}

func (self *Redis) HGet(key, field string) *string {
	result, err := self.Client.HGet(key, field).Result()
	if err != nil {
		if err.Error() == `redis: nil` {
			return nil
		}
		log.Logger.Error(`redis hget error`, err)
	}
	log.Logger.Debugf(`redis hget. key: %s, field: %s, value: %s`, key, field, result)
	return &result
}

func (self *Redis) HSet(key, field string, value interface{}) {
	log.Logger.Debugf(`redis hset. key: %s, field: %s, value: %s`, key, field, value)
	if err := self.Client.HSet(key, field, value).Err(); err != nil {
		log.Logger.Error(`redis hset error`, err)
	}
}

func (self *Redis) HSetExpiration(key, field string, value interface{}, expiration time.Duration) {
	log.Logger.Debug(fmt.Sprintf(`redis HSetExpiration. key: %s, val: %s`, key, field))
	_, err := self.Client.HSet(key, field, value).Result()
	if err != nil {
		log.Logger.Error(`redis HSetExpiration error`, err)
	}
	self.Client.Expire(key, expiration)
}

//添加元素
func (self *Redis) SetRainValue(id string, key string, v string) error {
	return self.Client.HSet(fmt.Sprintf("%v_%v", RainKey, id), key, v).Err()
}

//获得红包金额
func (self *Redis) GetRainValue(id, key string, v interface{}) error {
	return self.Client.HGet(fmt.Sprintf("%v_%v", RainKey, id), key).Scan(v)
}

func (self *Redis) Length(id string) int64 {
	return self.Client.HLen(fmt.Sprintf("%v_%v", RainKey, id)).Val()
}

//设置红包超时时间
func (self *Redis) Expire(id string, t time.Duration) error {
	return self.Client.Expire(fmt.Sprintf("%v_%v", RainKey, id), t).Err()
}

func (this *Redis) SetNx(key string, value string, expiration time.Duration) bool {
	result := this.Client.SetNX(key, value, expiration)
	if err := result.Err(); err != nil {
		log.Logger.Error(`redis set error`, err)
		return false
	}
	return result.Val()
}
