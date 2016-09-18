package bootstrap

import (
	"fmt"
	"strconv"

	"github.com/blueprint/blueprint/dba"
	"github.com/spf13/viper"
)

// DBA bootstraps the dba object and dials the specified databases.
func DBA() error {

	// get database configuration information
	dbs := viper.GetStringMapString("databases")

	for typ := range dbs {

		switch typ {
		case "mongo", "mongodb":
			db := viper.GetStringMapString("databases." + typ)
			mongo, err := dba.NewMongoDB(
				db["name"],
				db["user"],
				db["pass"],
				db["addr"],
			)
			if err != nil {
				return err
			}
			dba.AddMongoDB(mongo)

		case "redis", "redisdb":
			db := viper.GetStringMapString("databases." + typ)
			var poolsize int
			if db["pool_size"] != "" {
				i, err := strconv.Atoi(db["pool_size"])
				if err != nil {
					return err
				}
				poolsize = i
			}
			var dbint int64
			if db["db"] != "" {
				i, err := strconv.ParseInt(db["db"], 10, 64)
				if err != nil {
					return err
				}
				dbint = i
			}
			redis, err := dba.NewRedis(
				dbint,
				db["user"],
				db["pass"],
				db["addr"],
				poolsize,
			)
			if err != nil {
				return err
			}
			dba.AddRedisDB(redis)

		case "mysql", "postgres", "sqlite", "sql":
			db := viper.GetStringMapString("databases." + typ)
			sql, err := dba.NewSQL(
				db["name"],
				db["user"],
				db["pass"],
				db["addr"],
				typ,
			)
			if err != nil {
				return err
			}
			dba.AddSQL(sql)
		default:
			return fmt.Errorf("unsupported database %s", typ)
		}
	}
	return nil
}
