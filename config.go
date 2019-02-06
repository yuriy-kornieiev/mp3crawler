package main

import (
	"database/sql"
	"github.com/go-redis/redis"
)

//const (
//	BASE_PARSE_NAME     = "base"
//	SPIDERS_STATUS_INIT = iota
//	SPIDERS_STATUS_RUNNING
//	SPIDERS_STATUS_STOP
//	SPIDERS_BASIC_COOKIE
//)

type Config struct {
	Databases map[string]map[string]*Database
	//ActiveMQ    map[string]string
	//AWS         map[string]*AWS
	Redis       map[string]*Redis
	Environment string
	Source      string
}

func (c *Config) Load() {
	c.Databases = Database{}.Load()
	//c.ActiveMQ = ActiveMQ{}.Load()
	//c.AWS = AWS{}.Load()
	c.Redis = Redis{}.Load()
}

func (c *Config) GetDBConnection(source string) *sql.DB {
	return c.Databases[c.Environment][source].GetConnection()
}

func (c *Config) GetRedisConnection() *redis.Client {
	return c.Redis[c.Environment].GetConnection()
}

//func (c *Config) GetAWSConnection() *session.Session {
//	return c.AWS[c.Environment].GetConnection()
//}
