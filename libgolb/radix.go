package libgolb

import (
	"github.com/fzzy/radix/redis"
	"time"
)

var (
	LBClient    *redis.Client
)

func ConnectToRedis() (err error) {
	LBClient, err = redis.DialTimeout("tcp", Conf.RedisLB.Hostname+":"+Conf.RedisLB.Port, time.Duration(10)*time.Second)
	if err != nil {
		return err
	}
	r := LBClient.Cmd("select", Conf.RedisLB.Database)
	if r.Err != nil {
		return r.Err
	}
	return err
}

func RadixCheck(c *redis.Client, key string) (err error) {
	_, err = c.Cmd("get", key).Str()
	return err
}

func RadixSet(c *redis.Client, key, value string) error {
	r := c.Cmd("set", key, value)
	return r.Err
}

func RadixExpire(c *redis.Client, key, ttl string) error {
	r := c.Cmd("EXPIRE", key, ttl)
	return r.Err
}

func RadixUpdate(c *redis.Client, key, value string) error {
	r := c.Cmd("set", key, value)
	return r.Err
}

func RadixGetString(c *redis.Client, key string) (s string, err error) {
	s, err = c.Cmd("get", key).Str()
	return s, err
}

func RadixDel(c *redis.Client, key string) error {
	r := c.Cmd("del", key)
	return r.Err
}

func RadixList(c *redis.Client) (s []string, err error) {
	s, err = c.Cmd("KEYS", "*").List()
	return s, err
}

func SetMutex(c *redis.Client, key string) (result int, err error) {
	result, err = c.Cmd("SETNX", "lock:" + key, key).Int()
	return
}

func DelMutex(c *redis.Client, key string) (err error) {
	_, err = c.Cmd("SETNX", "lock:" + key).Int()
	return
}
