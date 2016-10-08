package db

import (
	"errors"
	mgo "gopkg.in/mgo.v2"
	"log"
	"sync"
)

type FC func(*mgo.Collection) (interface{}, error)
type FD func(*mgo.Database) (interface{}, error)

type MongoDB struct {
	sync.RWMutex
	sess     *mgo.Session
	uri      string
	database string
}

func (m *MongoDB) session() *mgo.Session {

	var err error = nil
	if m.sess == nil {
		m.sess, err = mgo.Dial(m.uri)
		if err != nil {
			return nil
		}
	}

	return m.sess.Clone()
}

func (m *MongoDB) Make(uri string) Database {

	m.uri = uri
	db := &MongoDB{
		sess:     m.session(),
		uri:      uri,
		database: "",
	}

	return db
}

func (m *MongoDB) Destroy() error {

	m.sess.Close()
	return nil
}

//Use change current active database
func (m *MongoDB) Use(db string) error {

	m.database = db
	return nil
}

//Drop delete a database or table/collection of database
func (m *MongoDB) Drop(db string, args ...interface{}) error {

	session := m.session()
	defer session.Close()

	size := len(args)
	if size == 0 {
		return session.DB(m.database).DropDatabase()
	} else if size == 1 {
		if collection, ok := args[0].(string); ok {
			return session.DB(m.database).C(collection).DropCollection()
		}
	}

	return errors.New("MongoDB: invalid args")
}

//Insert add a data to database
func (m *MongoDB) Insert(collection string, args ...interface{}) error {
	session := m.session()
	defer session.Close()

	c := session.DB(m.database).C(collection)

	return c.Insert(args...)
}

//Delete remove a record from database
func (m *MongoDB) Delete(collection string, args ...interface{}) error {
	if len(args) != 1 {
		return errors.New("MongoDB: invalid args")
	}

	session := m.session()
	defer session.Close()

	c := session.DB(m.database).C(collection)

	return c.Remove(args[0])
}

//Update update a record already in database
func (m *MongoDB) Update(collection string, args ...interface{}) error {

	if len(args) != 2 {
		return errors.New("MongoDB: invalid args")
	}

	session := m.session()
	defer session.Close()

	c := session.DB(m.database).C(collection)

	return c.Update(args[0], args[1])
}

//Query retrieve a record from database
func (m *MongoDB) Query(collection string, args ...interface{}) interface{} {

	if len(args) != 1 {
		return errors.New("MongoDB: invalid args")
	}

	session := m.session()
	defer session.Close()

	c := session.DB(m.database).C(collection)

	return c.Find(args[0])
}

//Count returns total number of records in collection, otherwise, -1 returned
func (m *MongoDB) Count(collection string) int {

	session := m.session()
	defer session.Close()
	c := session.DB(m.database).C(collection)

	if n, err := c.Count(); err == nil {
		return n
	} else {
		log.Println(err)
		return -1
	}
}

//Ioctrl implement other operations, cmd support "database" and "collection", if "database" specified, Ioctrl accept a
//function as second parameter, whose prototype MUST be func(*mgo.Database) (interface{}, error), and if "collection"
//specified, Ioctrl accept a collection name as second parameter, and a function as third parameter, whose prototype
//MUST be func(*mgo.Collection) (interface{}, error).
func (m *MongoDB) Ioctrl(cmd string, args ...interface{}) (interface{}, error) {

	session := m.session()
	defer session.Close()

	if cmd == "database" && len(args) == 1 {
		db := session.DB(m.database)
		if fd, ok := args[0].(FD); ok {
			return fd(db)
		}
	} else if cmd == "collection" && len(args) == 2 {
		if collection, ok := args[0].(string); ok {
			c := session.DB(m.database).C(collection)
			if fc, ok := args[1].(FC); ok {
				return fc(c)
			}
		}
	}

	return nil, errors.New("MongDB: invalid args")
}

//Factory return a Factory interface instance, which can create and destroy a database engine
func (m *MongoDB) Factory() Factory {
	return m
}
