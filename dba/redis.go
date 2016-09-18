package dba

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/redis.v3"
)

// RedisDB describes a redis database object.
type RedisDB struct {
	client           *redis.Client
	user, pass, addr string
	poolSize         int
	db               int64
}

// NewRedis creates a new redis database object.
func NewRedis(db int64, user, pass, addr string, poolSize int) (*RedisDB, error) {
	if addr == "" {
		return nil, redisError(db, errors.New("missing redis address"))
	}
	return &RedisDB{
		db:       db,
		user:     user,
		pass:     pass,
		addr:     addr,
		poolSize: poolSize,
	}, nil
}

// Name satisfies the Database interface.  The redis database object returns
// the name of the last database the Dial method attempted to connect with.
func (r *RedisDB) Name() string {
	return strconv.FormatInt(r.db, 10)
}

// Dial satisfies the Database interface. The Redis database object attempts
// to dial and start a session with a local or remote Redis instance.
func (r *RedisDB) Dial(name string) error {

	// catch any potential errors
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(redisError(r.db, fmt.Errorf("%s", err)))
		}
	}()

	i, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		return redisError(r.db, err)
	}

	// establish a client
	client := redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: r.pass,
		DB:       i,
		PoolSize: r.poolSize,
	})

	// ping test the connection
	_, err = client.Ping().Result()
	if err != nil {
		return redisError(r.db, err)
	}

	// set database client
	r.client = client

	// double check the database session connection was sucessful
	if r.client == nil {
		return redisError(r.db, errors.New("client is nil"))
	}

	return nil
}

// Close satisfies the Database interface.  The redis database object closes
// the connection the Dial method created.
func (r *RedisDB) Close() {
	// if there is an error, recover, log, and stop the world
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(redisError(r.db, fmt.Errorf("%s", err)))
		}
	}()
	r.client.Close()
}

func redisError(db int64, err error) error {
	return fmt.Errorf("redis db %d error: %s", db, err)
}
