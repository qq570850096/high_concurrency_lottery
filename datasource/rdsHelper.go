package datasource

import (
	"High/conf"
	"fmt"
	"log"
	"github.com/gomodule/redigo/redis"
	"time"
)

var cacheInstance *RedisConn

func InstanceCache() *RedisConn {
	once.Do(func() {
		cacheInstance = NewCache()
	})
	return cacheInstance
}
type RedisConn struct {
	pool *redis.Pool
	showDebug bool
}

func (rds *RedisConn) Do(commandName string,
	args ...interface{}) (reply interface{}, err error) {
	// 从连接池里拿一个连接
	conn := rds.pool.Get()
	// 用完后放回
	defer conn.Close()

	t1 := time.Now().UnixNano()
	reply,err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("rdsHelper.Do",err,e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		fmt.Printf("[redis] [info] [%dus] comd = %s," +
			"err = %s, args = %v, reply = %s\n",
			(t2-t1)/1000,commandName,err,args,reply)
	}
	return reply,err
}

// 设置是否需要Debug
func (rds *RedisConn) ShowDebug(b bool) {
	rds.showDebug = b
}
// 新实例
func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			c,err := redis.Dial("tcp" ,fmt.Sprintf("%s:%d",
				conf.RdsCache.Host,conf.RdsCache.Port))
			if err != nil {
				log.Fatal("rdsHelper NewCache error=",err)
				return nil,err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute{
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000,
		MaxActive:       10000,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	instance := &RedisConn{
		pool:      &pool,
		showDebug: false,
	}
	instance.ShowDebug(true)
	return instance
}
