package cache

import (
	"errors"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/raythorn/falcon/log"
	"time"
)

//Redis engine for falcon Cache
//
//Implement interface Factory and Cache, interface Factory used for init and destroy engine,
//and interface Cache used for manipulating the cache, all function but Ioctrl only support
//for strings and hash, if you want use other data structure or complex operations, please
//refer to Ioctrl, this function is just a wrap of redis.Do()
type Redis struct {
	pool *redigo.Pool
	uri  string
}

//Make create a Redis instance, and return Cache
func (r *Redis) Make(uri string) Cache {
	return &Redis{
		pool: &redigo.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redigo.Conn, error) {
				conn, err := redigo.DialURL(uri, redigo.DialConnectTimeout(5*time.Second), redigo.DialReadTimeout(5*time.Second), redigo.DialWriteTimeout(5*time.Second))
				if err != nil {
					return nil, err
				}

				return conn, nil
			},
			TestOnBorrow: func(conn redigo.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := conn.Do("ping")
				return err
			},
		},
		uri: uri,
	}
}

//Destroy cleanup context if not used
func (r *Redis) Destroy() error {
	return nil
}

//Factory returns interface Factory
func (r *Redis) Factory() Factory {
	return r
}

//Set data to redis with key and args.
//
//Set will auto parse args to determin which data structure to use. you cannot use "ex", "px", "nx", "xx" as keys, they are reserved for engine using.
//	For strings
//		Set("key", 345, "value")                         --> SETEX key 345 value
//		Set("key", "value", "nx|xx")                     --> SET key value nx|xx
//		Set("key", "value", "ex|px" 345)                 --> SET key value ex|px 345
//	For hash
//		Set("key", "field", "value")                     --> HSET key field value
//		Set("key", "field", "valude"[,"filed", "value"]) --> HMSET key field value
//
//Note: You cannot set a value to different data structure
func (r *Redis) Set(key string, args ...interface{}) error {

	conn := r.pool.Get()
	if nil == conn {
		return errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	var err error = nil
	var cmd string = ""

	size := len(args)

	if size == 1 {
		cmd = "SET"
	} else if size == 2 {
		switch args[0].(type) {
		case int, int16, int32, int64, uint, uint16, uint32, uint64:
			cmd = "SETEX"
		case string:
			if val, ok := args[1].(string); ok && (val == "nx" || val == "xx") {
				cmd = "SET"
			} else {
				cmd = "HSET"
			}
		}
	} else if size == 3 {
		if val, ok := args[1].(string); ok && (val == "ex" || val == "px") {
			switch args[2].(type) {
			case int, int16, int32, int64, uint, uint16, uint32, uint64:
				cmd = "SET"
			}
		}
	} else if size >= 4 && size%2 == 0 {
		cmd = "HMSET"
	}

	if cmd != "" {

		args1 := []interface{}{key}
		for _, arg := range args {
			args1 = append(args1, arg)
		}

		_, err = conn.Do(cmd, args1...)

	} else {
		err = errors.New("Redis: invalid args")
	}

	return err
}

//Get cached from redis server
//
//Get will auto parse args to determin which data structure to use.
//	For strings
//		Get("key")                    --> GET	key	(key is key of strings structure)
//	For hash
//		Get("key", "field"[,"field"]) --> HMGET key filed [field]
func (r *Redis) Get(key string, args ...string) interface{} {

	conn := r.pool.Get()
	if nil == conn {
		log.Error("Redis: get connection from pool failed")
		return nil
	}
	defer conn.Close()

	var cmd string = "GET"
	size := len(args)

	if size > 0 {
		cmd = "HMGET"
	}

	var val interface{}
	var err error

	if cmd == "GET" {
		val, err = conn.Do(cmd, key)
	} else {
		args1 := []interface{}{key}
		for _, arg := range args {
			args1 = append(args1, arg)
		}

		val, err = conn.Do(cmd, args1...)
	}

	if err == nil {
		return val
	}

	log.Error(err.Error())

	return nil
}

//Incr increase key's value
//
//	Incr("key")              --> INCR key
//	Incr("key", 10)          --> INCRBY key 10
//	Incr("key", "field", 10) --> HINCRBY key field 10
func (r *Redis) Incr(key string, args ...interface{}) error {
	conn := r.pool.Get()
	if nil == conn {
		return errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	var cmd string = ""
	var err error = nil

	size := len(args)
	if size == 0 {
		cmd = "INCR"
	} else if size == 1 {
		cmd = "INCRBY"
	} else if size == 2 {
		cmd = "HINCRBY"
	} else {
		err = errors.New("Redis: invalid args")
	}

	if cmd != "" && err == nil {
		args1 := []interface{}{key}
		for _, arg := range args {
			args1 = append(args1, arg)
		}
		_, err = conn.Do(cmd, args1...)
		return err
	}

	return errors.New("Redis: invalid args")
}

//Decr decrease key's value
//
//	Decr("key")              --> DECR key
//	Decr("key", 10)	         --> DECRBY key 10
//	Decr("key", "field", 10) --> HDECRBY key field 10
func (r *Redis) Decr(key string, args ...interface{}) error {
	conn := r.pool.Get()
	if nil == conn {
		return errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	var cmd string = ""
	var err error = nil

	size := len(args)
	if size == 0 {
		cmd = "DECR"
	} else if size == 1 {
		cmd = "DECRBY"
	} else if size == 2 {
		cmd = "HDECRBY"
	} else {
		err = errors.New("Redis: invalid args")
	}

	if cmd != "" && err == nil {
		args1 := []interface{}{key}
		for _, arg := range args {
			args1 = append(args1, arg)
		}
		_, err = conn.Do(cmd, args1...)
		return err
	}

	return errors.New("Redis: invalid args")
}

//Exist check if key or field existed
//
//	Exist("key")                     --> EXISTS key
//	Exist("key", "field"[, "field"]) --> EXISTS key field [field]
func (r *Redis) Exist(key string, args ...interface{}) bool {

	conn := r.pool.Get()
	if nil == conn {
		log.Error("Redis: get connection from pool failed")
		return false
	}
	defer conn.Close()

	cmd := ""
	size := len(args)
	if size == 0 {
		cmd = "EXISTS"
	} else if size == 1 {
		cmd = "HEXISTS"
	} else {
		log.Error("Redis: invalid args")
		return false
	}

	args1 := []interface{}{key}
	for _, arg := range args {
		args1 = append(args1, arg)
	}
	if exists, err := redigo.Bool(conn.Do(cmd, args1...)); err == nil {
		return exists
	} else {
		log.Error(err.Error())
		return false
	}
}

//Delete delete a key or field
//
//	Delete("key")          --> DEL key
//	Delete("key", "field") --> HDEL key field
func (r *Redis) Delete(key string, keys ...string) error {
	conn := r.pool.Get()
	if nil == conn {
		return errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	cmd := "DEL"
	size := len(keys)

	if size > 0 {

		cmd = "HDEL"
	}

	keys1 := []interface{}{key}
	for _, k := range keys {
		keys1 = append(keys1, k)
	}
	_, err := conn.Do(cmd, keys1...)

	return err
}

//Expire set expiration of a key
func (r *Redis) Expire(key string, time int64) error {
	conn := r.pool.Get()
	if nil == conn {
		return errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	args := []interface{}{key, time}

	_, err := conn.Do("EXPIRE", args...)

	return err
}

//TTL query a key's ttl
func (r *Redis) TTL(key string) int64 {
	conn := r.pool.Get()
	if nil == conn {
		log.Error("Redis: get connection from pool failed")
		return -2
	}
	defer conn.Close()

	if ttl, err := redigo.Int64(conn.Do("TTL", key)); err == nil {
		return ttl
	}

	return -2
}

//Ioctrl handle all io operations of redis, it just a wrap of redis.Conn.Do()
func (r *Redis) Ioctrl(cmd string, args ...interface{}) (result interface{}, err error) {
	conn := r.pool.Get()
	if nil == conn {
		return nil, errors.New("Redis: get connection from pool failed")
	}
	defer conn.Close()

	return conn.Do(cmd, args...)
}
