// Package dba is an singleton that manages databases as would a Database
// Administrator.
package dba

import (
	"strconv"

	"gopkg.in/gorp.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v3"
)

// admin is a database administrator object that manages different SQL and NoSQL
// databases. The admin object follows the singleton pattern in the db package
// so that admin is easy to access across a project.
var admin *DatabaseAdmin

func init() {
	admin = &DatabaseAdmin{}
}

// Database describes a methods to connect and close database connections.
type Database interface {
	Dial(name string) error
	Name() string
	Close()
}

// DatabaseAdmin represents a database administrator object who's duty is to
// maintain all database values, connections, and state.
type DatabaseAdmin struct {
	Mongo []Database
	Redis []Database
	SQL   []Database
}

// MGO takes a database name to search all mongodb databases db.DBA currently
// maintains and returns a pointer to a mgo.Database instance if found.
// mgo.Database represents a mongodb orm-ish driver for executing queries.
// If no database is found nil is returned.
func (d DatabaseAdmin) MGO(name string) *mgo.Database {
	for _, db := range d.Mongo {
		if db.Name() == name {
			if m, ok := db.(*MongoDB); ok {
				return m.mongodb
			}
			return nil
		}
	}
	return nil
}

// AddMongoDB adds a MongoDB database object to the DatabaseAdmin databases.
func (d *DatabaseAdmin) AddMongoDB(m *MongoDB) {
	d.Mongo = append(d.Mongo, m)
}

// RedisClient takes a database name to search all redis databases db.DBA
// currently maintains and returns a pointer to a redis.Client instance if
// found. redis.Client represents a redis driver for executing queries.
// If no database is found nil is returned.
func (d DatabaseAdmin) RedisClient(db int64) *redis.Client {
	name := strconv.FormatInt(db, 10)
	for _, db := range d.Redis {
		if db.Name() == name {
			if r, ok := db.(*RedisDB); ok {
				return r.client
			}
			return nil
		}
	}
	return nil
}

// AddRedisDB adds a Redis database object to the DatabaseAdmin databases.
func (d *DatabaseAdmin) AddRedisDB(r *RedisDB) {
	r.db = int64(len(d.Redis))
	d.Redis = append(d.Redis, r)
}

// GORP takes a database name to search all SQL databases DatabaseAdmin
// currently maintains and returns a pointer to a gorp.DbMap instance if found.
// gorp.DbMap represents an SQL orm-ish driver for multiple different SQL
// databases (sqlite, mysql, postgres). If no database is found nil is returned.
func (d DatabaseAdmin) GORP(name string) *gorp.DbMap {
	for _, db := range d.SQL {
		if db.Name() == name {
			if m, ok := db.(*SQL); ok {
				return m.sql
			}
			return nil
		}
	}
	return nil
}

// AddSQL adds a sql database object to the DatabaseAdmin databases.
func (d *DatabaseAdmin) AddSQL(s *SQL) {
	d.SQL = append(d.SQL, s)
}

// All returns all the databases this DatabaseAdmin maintains.
func (d *DatabaseAdmin) All() []Database {
	dbs := append(d.Mongo, d.SQL...)
	dbs = append(dbs, d.Redis...)
	return dbs
}

// MGO takes a database name to search all mongodb databases the db package
// currently maintains and returns a pointer to a mgo.Database instance if
// found. mgo.Database represents a mongodb orm-ish driver for executing queries.
// If no database is found nil is returned.
func MGO(name string) *mgo.Database {
	return admin.MGO(name)
}

// AddMongoDB adds a MongoDB database object to the db package.
func AddMongoDB(m *MongoDB) {
	admin.AddMongoDB(m)
}

// RedisClient takes a database name to search all redis databases the db package
// currently maintains and returns a pointer to a redis.Client instance if
// found. redis.Client represents a RedisClient client for executing queries.
// If no database is found nil is returned.
func RedisClient(db int64) *redis.Client {
	return admin.RedisClient(db)
}

// AddRedisDB adds a Redis database client to the db package.
func AddRedisDB(r *RedisDB) {
	admin.AddRedisDB(r)
}

// GORP takes a database name to search all SQL databases the db package
// currently maintains and returns a pointer to a gorp.DbMap instance if found.
// gorp.DbMap represents an SQL orm-ish driver for multiple different SQL
// databases (sqlite, mysql, postgres). If no database is found nil is returned.
func GORP(name string) *gorp.DbMap {
	return admin.GORP(name)
}

// AddSQL adds a sql database object to the singleton DatabaseAdmin's databases.
func AddSQL(s *SQL) {
	admin.AddSQL(s)
}

// All returns all the databases this singleton DatabaseAdmin maintains.
func All() []Database {
	return admin.All()
}
