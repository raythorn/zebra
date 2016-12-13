//Package db is a database mechanism, which helps accessing data more convenient and efficient.
//
//This mechanism provide two interfaces, Database and Factory, and you MUST implement these two
//interfaces to add your own database engine. Currently, falcon support ONLY MongoDB with driver
//mgo(http://labix.org/mgo).
package db

import (
	"errors"
)

var (
	dbInstance *database
)

func init() {
	dbInstance = &database{nil}
}

//Database is interface which MUST be implemented to use falcon database mechanism
//
//This interface provide several base operations for database access, such as, Insert, Delete, Update, Query
//, etc, with which you can access the data with unified APIs, and it's very convenient and efficient.
type Database interface {
	//Use change current active database
	Use(db string) error

	//Drop delete a database or table/collection of database
	Drop(db string, args ...interface{}) error

	//Insert add a data to database
	Insert(collection string, args ...interface{}) error

	//Delete remove a record from database
	Delete(collection string, args ...interface{}) error

	//Update update a record already in database
	Update(collection string, args ...interface{}) error

	//Query retrieve a record from database
	Query(collection string, args ...interface{}) error

	//Count returns total number of records in collection
	Count(collection string) int

	//Ioctrl implement other operations
	Ioctrl(cmd string, args ...interface{}) (interface{}, error)

	//Factory return a Factory interface instance, which can create and destroy a database engine
	Factory() Factory
}

//Factory is a interface MUST be implemented to use falcon database mechanism
//
//This interface provide two APIs to create and destroy database engine instance
type Factory interface {
	//Make create a new database engine with URI string, which implemented Database interface.
	//URI string MUST like: mongodb://user:password@192.168.100.2:2014, depends on each implementation of engine
	Make(uri string) Database

	//Destroy destroy and cleanup engine instance if no further use
	Destroy() error
}

type database struct {
	engine Database
}

func Register(uri string, factory Factory) error {

	if engine := factory.Make(uri); engine != nil {
		dbInstance.engine = engine
		return nil
	}

	return errors.New("DB: register failed")
}

func UnRegister() error {
	if dbInstance != nil {
		if err := dbInstance.engine.Factory().Destroy(); err == nil {
			return nil
		}
	}

	return errors.New("DB: unregister failed")
}

func Use(db string) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}

	return dbInstance.engine.Use(db)
}

func Drop(db string, args ...interface{}) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Drop(db, args...)
}

func Insert(collection string, args ...interface{}) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Insert(collection, args...)
}

func Delete(collection string, args ...interface{}) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Delete(collection, args...)
}

func Update(collection string, args ...interface{}) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Update(collection, args...)
}

func Query(collection string, args ...interface{}) error {
	if dbInstance.engine == nil {
		return errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Query(collection, args...)
}

func Ioctrl(cmd string, args ...interface{}) (interface{}, error) {
	if dbInstance.engine == nil {
		return nil, errors.New("DB: engine invalid")
	}
	return dbInstance.engine.Ioctrl(cmd, args...)
}
