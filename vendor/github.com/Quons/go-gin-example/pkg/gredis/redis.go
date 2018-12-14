package gredis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Quons/go-gin-example/pkg/setting"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

//key前缀，防止key冲突
var globalPrefix = setting.RedisSetting.Prefix
var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

//添加缓存方法。注意参数顺序：缓存对象，过期时间（单位为秒）,key前缀，key值数组
func Set(data interface{}, expire int, prefix string, keyValue ...interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()
	encodeData, err := Encode(data)
	if err != nil {
		logrus.Error(err)
		return err
	}
	key := getCacheKey(prefix, keyValue)
	_, err = conn.Do("SET", key, encodeData)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, expire)
	if err != nil {
		return err
	}

	return nil
}

/*判断key是否存在方法，存在返回true，不存在返回false*/
func Exists(prefix string, keyValue ...interface{}) bool {
	conn := RedisConn.Get()
	defer conn.Close()
	key := getCacheKey(prefix, keyValue)
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

//获取key方法，存在返回true，否则返回false。注意参数顺序：接受对象，key前缀，key值数组
func Get(to interface{}, prefix string, keyValue ...interface{}) bool {
	conn := RedisConn.Get()
	defer conn.Close()
	key := getCacheKey(prefix, keyValue)
	exist := Exists(key)
	if !exist {
		return false
	}
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logrus.Error(err)
		return false
	}
	if err = Decode(reply, to); err != nil {
		logrus.Error(err)
		return false
	}
	return true
}

/*删除key方法*/
func Delete(prefix string, keyValue ...interface{}) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	key := getCacheKey(prefix, keyValue)
	return redis.Bool(conn.Do("DEL", key))
}

/*关键字删除key方法*/
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

/*拼接缓存key*/
func getCacheKey(prefix string, keyValue ...interface{}) string {
	key := prefix
	for _, v := range keyValue {
		key = fmt.Sprintf(globalPrefix+key+"%v", v)
	}
	return key
}

// --------------------
// Encode
// 用gob进行数据编码
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
