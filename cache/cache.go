//Package cache is a simple cache machnism for zebra, it support redis for now,
//and you can add your own cache engine with implement Factory and Cache.
package cache

import (
	"errors"
)

func init() {
	cacheInstance = &cache{nil}
}

var (
	cacheInstance *cache
)

//Cache is a interface which is used to interact with cache, you MUST implement this to use zebra's cache mechanism
type Cache interface {
	//Set data to cache
	Set(key string, args ...interface{}) error

	//Get data from cache
	Get(key string, args ...string) interface{}

	//Delete a data from cache
	Delete(key string, args ...string) error

	//Query if key exists
	Exist(key string, args ...interface{}) bool

	//Increase a key's value
	Incr(key string, args ...interface{}) error

	//Decrease a key's value
	Decr(key string, args ...interface{}) error

	//Set a key's expiration time
	Expire(key string, time int64) error

	//Query a key's ttl
	TTL(key string) int64

	//Custom io operations
	Ioctrl(cmd string, args ...interface{}) (interface{}, error)

	//Return a factory instance
	Factory() Factory
}

//Factory is a interface which is used init and cleanup cache context, you MUST implement this to use zebra's cache mechanism
type Factory interface {
	//Make initialise a Cache instance
	Make(uri string) Cache

	//Destroy cleanup cache context
	Destroy() error
}

type cache struct {
	engine Cache
}

//Register register cache engine and make it ready to use
func Register(uri string, factory Factory) error {

	if engine := factory.Make(uri); engine != nil {
		cacheInstance.engine = engine
		return nil
	}

	return errors.New("Cache: Register cache engine failed")
}

//UnRegister unregister a cache engine and cleanup context
func UnRegister(engine string) error {

	if cacheInstance.engine != nil {
		if err := cacheInstance.engine.Factory().Destroy(); err == nil {
			cacheInstance.engine = nil
			return nil
		}
	}

	return errors.New("Cache: UnRegister failed")
}

//Set data to cache
func Set(key string, args ...interface{}) error {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Set(key, args...)
}

//Get data from cache
func Get(key string, args ...string) interface{} {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Get(key, args...)
}

//Delete a data from cache
func Delete(key string, args ...string) error {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Delete(key, args...)
}

//Query if key exists
func Exist(key string, args ...interface{}) bool {
	if cacheInstance.engine == nil {
		return false
	}

	return cacheInstance.engine.Exist(key, args...)
}

//Increase a key's value
func Incr(key string, args ...interface{}) error {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Incr(key, args...)
}

//Decrease a key's value
func Decr(key string, args ...interface{}) error {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Decr(key, args...)
}

//Set a key's expiration time
func Expire(key string, time int64) error {
	if cacheInstance.engine == nil {
		return errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Expire(key, time)
}

//Query a key's ttl
func TTL(key string) int64 {
	if cacheInstance.engine == nil {
		return -2
	}

	return cacheInstance.engine.TTL(key)
}

//Custom io operations
func Ioctrl(cmd string, args ...interface{}) (interface{}, error) {
	if cacheInstance.engine == nil {
		return nil, errors.New("Cache: engine invalid")
	}

	return cacheInstance.engine.Ioctrl(cmd, args...)
}
