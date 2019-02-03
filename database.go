package main

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/segmentio/go-athena"
	"log"
	"time"
)

type Database struct {
	Type     string
	Port     string
	Host     string
	DB       string
	User     string
	Password string
	Region   string
	Output   string

	Connection *sql.DB
}

func (d Database) Load() map[string]map[string]*Database {

	data := make(map[string]map[string]*Database)

	data["yuriyk"] = map[string]*Database{}

	data["yuriyk"]["crawler"] = &Database{
		"mysql",
		"3306",
		"127.0.0.1",
		"crawler",
		"test",
		"test",
		"",
		"",

		nil,
	}

	return data

}

func (d *Database) OpenConnection() *sql.DB {

	if d.Type == "athena" {
		return d.GetAthenaConnection()
	}

	var opxConnUrl string
	if d.Type == "postgres" {
		opxConnUrl = "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.DB

	} else {
		opxConnUrl = d.User + ":" + d.Password + "@" + "tcp(" + d.Host + ":" +
			d.Port +
			")/" +
			d.DB + "?" +
			"parseTime=true"

	}

	var err error

	if d.Connection, err = sql.Open(d.Type, opxConnUrl); err != nil {
		log.Println(d.Type, opxConnUrl)
		log.Panic(err)
	}

	log.Println("DB Connection established: " + d.Host + ":" + d.Port + "/" + d.DB)

	d.Connection.SetConnMaxLifetime(time.Minute * 10)
	d.Connection.SetMaxIdleConns(10)
	d.Connection.SetMaxOpenConns(10)

	return d.Connection
}

func (d *Database) GetAthenaConnection() *sql.DB {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(d.Region),
		Credentials: credentials.NewStaticCredentials(d.User, d.Password, ""),
	})

	if err != nil {
		log.Panic(err)
	}

	d.Connection, err = athena.Open(athena.Config{
		Session:        sess,
		Database:       d.DB,
		OutputLocation: d.Output,
	})

	if err != nil {
		log.Panic(err)
	}

	log.Println("Athena DB Connection established: " + d.Region + ":" + d.User + "/" + d.DB)

	return d.Connection

}

func (d *Database) GetConnection() *sql.DB {

	if d.Connection == nil {
		d.Connection = d.OpenConnection()
	}

	err := d.Connection.Ping()
	for err != nil {
		log.Println(err)
		time.Sleep(time.Second)
		err = d.Connection.Ping()
	}

	return d.Connection
}
