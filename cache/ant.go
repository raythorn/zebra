package cache

import (
	"errors"
	"sync"
	"time"
)

const (
	//AntClockTick means the clock time of recycling the expired cache in memory.
	AntClockTick = 60 //60 seconds
)

//Ant is a light weight memory cache, all data will be cached in memory, and will be recycled if
//application stopped
type Ant struct {
	sync.RWMutex
	cache  map[string]*particle
	uri    string
	clock  time.Duration
	finish bool
}

type particle struct {
	duration time.Duration
	cat      time.Time
	pit      interface{}
}

func (p *particle) expire() bool {
	if 0 == p.duration {
		return false
	}

	return time.Now().Sub(p.cat) > p.duration
}

func (ant *Ant) Make(uri string) Cache {
	antInstance := &Ant{
		cache:  make(map[string]*particle),
		uri:    uri,
		clock:  AntClockTick,
		finish: false,
	}

	go antInstance.ticker()

	return antInstance
}

func (ant *Ant) Destroy() error {
	ant.finish = true

	return nil
}

//Set set data to cache
//
//	Set("key", value)                   --> cache key with value
//	Set("key", value, "nx|xx")          --> cache key with value if key not exist or exist
//	Set("key", value, "ex", expiration) --> cache key with value and expiration time
func (ant *Ant) Set(key string, args ...interface{}) error {

	ant.Lock()
	defer ant.Unlock()

	p := &particle{cat: time.Now()}
	size := len(args)

	var err error = nil

	if size == 1 {
		p.duration = 0
	} else if size == 2 {
		_, ok := ant.cache[key]

		if (args[1] == "nx" && !ok) || (args[1] == "xx" && ok) {
			p.duration = 0
		} else {
			err = errors.New("Ant: invalid parameters")
		}
	} else if size == 3 {
		if args[1] == "ex" {
			switch args[2].(type) {
			case int:
				p.duration = time.Duration(args[2].(int))
			case int32:
				p.duration = time.Duration(args[2].(int32))
			case int64:
				p.duration = time.Duration(args[2].(int64))
			case uint:
				p.duration = time.Duration(args[2].(uint))
			case uint32:
				p.duration = time.Duration(args[2].(uint32))
			case uint64:
				p.duration = time.Duration(args[2].(uint64))
			case time.Duration:
				p.duration = args[2].(time.Duration)
			default:
				err = errors.New("Ant: expiration valude is not (u)int, (u)int32, (u)int64, or time.Duration")
			}
		} else {
			err = errors.New("Ant: invalid parameters")
		}
	} else {
		err = errors.New("Ant: invalid parameters")
	}

	if err == nil {
		p.pit = args[0]
		ant.cache[key] = p
	}

	return err
}

//Get get data from cache, get accept one or more than one key, if multi-key specified, a slice []interface{} will return,
//otherwise, interface{} will return
func (ant *Ant) Get(key string, args ...string) interface{} {

	ant.Lock()
	defer ant.Unlock()

	var values []interface{} = make([]interface{}, 0)

	keys := []interface{}{key, args}

	for _, k := range keys {
		key := k.(string)
		if val, ok := ant.cache[key]; ok && !val.expire() {
			values = append(values, val.pit)
		} else {
			values = append(values, "nil")
		}
	}

	if len(args) == 0 {
		return values[0]
	}

	return values
}

//Delete delete a data from cache, multi-key available
func (ant *Ant) Delete(key string, args ...string) error {

	ant.Lock()
	defer ant.Unlock()

	keys := []interface{}{key, args}
	for _, k := range keys {
		key := k.(string)
		delete(ant.cache, key)
	}

	return nil
}

//Exist query if key exists, only single key supported, args will be ignored
func (ant *Ant) Exist(key string, args ...interface{}) bool {

	ant.Lock()
	defer ant.Unlock()

	if val, ok := ant.cache[key]; ok && !val.expire() {
		return true
	}

	return false
}

//Incr increase a key's value, if key not exist, ant will auto create a key with 0 value, and increase it
//
//	Incr("key")      --> increase value by 1
//	Incry("key", 10) --> increase value by 10
func (ant *Ant) Incr(key string, args ...interface{}) error {

	ant.Lock()
	defer ant.Unlock()

	size := len(args)

	var count int64 = 0
	var value *particle = nil
	var err error = nil

	ok := false

	if value, ok = ant.cache[key]; ok && !value.expire() {
		count = value.pit.(int64)
	} else {
		value = &particle{
			pit:      count,
			duration: 0,
			cat:      time.Now(),
		}
	}

	if size == 0 {
		count++
	} else if size == 1 {
		switch args[0].(type) {
		case int:
			count += int64(args[0].(int))
		case int32:
			count += int64(args[0].(int32))
		case int64:
			count += args[0].(int64)
		case uint:
			count += int64(args[0].(uint))
		case uint32:
			count += int64(args[0].(uint32))
		case uint64:
			count += int64(args[0].(uint64))
		case time.Duration:
			count += int64(args[0].(time.Duration))
		default:
			err = errors.New("Ant: expiration time MUST be one of (u)int, (u)int32, (u)int64 and time.Duration")
		}
	} else {
		err = errors.New("Ant: too many args")
	}

	if err == nil {
		value.pit = count
		ant.cache[key] = value
	}

	return err
}

//Decr decrease a key's value, if key not exist, ant will auto create a key with 0 value, and decrease it
//
//	Decr("key")      --> decrease value by 1
//	Decry("key", 10) --> decrease value by 10
func (ant *Ant) Decr(key string, args ...interface{}) error {
	ant.Lock()
	defer ant.Unlock()

	size := len(args)

	var count int64 = 0
	var value *particle = nil
	var err error = nil

	ok := false

	if value, ok = ant.cache[key]; ok && !value.expire() {
		count = value.pit.(int64)
	} else {
		value = &particle{
			pit:      count,
			duration: 0,
			cat:      time.Now(),
		}
	}

	if size == 0 {
		count--
	} else if size == 1 {
		switch args[0].(type) {
		case int:
			count -= int64(args[0].(int))
		case int32:
			count -= int64(args[0].(int32))
		case int64:
			count -= args[0].(int64)
		case uint:
			count -= int64(args[0].(uint))
		case uint32:
			count -= int64(args[0].(uint32))
		case uint64:
			count -= int64(args[0].(uint64))
		case time.Duration:
			count -= int64(args[0].(time.Duration))
		default:
			err = errors.New("Ant: expiration time MUST be one of (u)int, (u)int32, (u)int64 and time.Duration")
		}
	} else {
		err = errors.New("Ant: too many args")
	}

	if err == nil {
		value.pit = count
		ant.cache[key] = value
	}

	return err
}

//Expire set a key's expiration time
func (ant *Ant) Expire(key string, exp int64) error {

	ant.Lock()
	defer ant.Unlock()

	if val, ok := ant.cache[key]; ok && !val.expire() {
		val.duration = time.Duration(exp)
	} else {
		return errors.New("Ant: key not exist")
	}

	return nil
}

//TTL query a key's ttl, if expiration time not set return 0, if time expired return -2,
//otherwise, return left time
func (ant *Ant) TTL(key string) int64 {
	ant.Lock()
	defer ant.Unlock()

	if val, ok := ant.cache[key]; ok && !val.expire() {
		return int64(val.duration)
	}

	return -2
}

//Custom io operations, not supported
func (ant *Ant) Ioctrl(cmd string, args ...interface{}) (interface{}, error) {
	return nil, errors.New("Ant: Not support Ioctrl")
}

//Return a factory instance
func (ant *Ant) Factory() Factory {
	return ant
}

//ticker is a recycling goroutine, will recycle expired data in each clock time
func (ant *Ant) ticker() {
	for !ant.finish {

		<-time.After(ant.clock)
		if ant.cache == nil {
			return
		}

		for key := range ant.cache {
			ant.sentinel(key)
		}
	}
}

//sentinel check if a specified key is expired, delete key if expired
func (ant *Ant) sentinel(key string) {

	ant.Lock()
	defer ant.Unlock()

	if particle, ok := ant.cache[key]; ok && particle.expire() {
		delete(ant.cache, key)
	}
}
