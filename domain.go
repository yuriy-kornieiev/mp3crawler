package main

import (
	"database/sql"
	"github.com/go-redis/redis"
	"log"
	"net/url"
	"sync"
	"time"
)

type Domain struct {
	sync.Mutex

	ID          sql.NullInt64
	Host        sql.NullString
	DateCreated time.Time
}

var domainSkipList = []string{"yandex.ru", "vk.com", "facebook.com", "mail.ru", "google.com", "google.ru", "google.com.ua"}

func (d Domain) GetDomain(url url.URL) bool {
	client := config.GetRedisConnection()

	currentTime := time.Now().Format(layout)

	val, err := client.Get("domain:" + url.Host).Result()
	if err == redis.Nil {
		domain := d.Get(url)
		if domain.ID.Valid == false {
			d.Insert(url)
		}
		client.Set("domain:"+url.Host, currentTime, 0)
		val = currentTime

	} else if err != nil {
		panic(err)
	}

	t1, _ := time.Parse(layout, val)
	t2, _ := time.Parse(layout, currentTime)
	t1 = t1.Add(time.Second * 20)

	if t2.After(t1) {
		client.Set("domain:"+url.Host, currentTime, 0)
		return true
	}

	return false
}

func (d Domain) Get(url url.URL) Domain {
	var err error
	var rows *sql.Rows
	conn := config.GetDBConnection("crawler")

	sql := "SELECT id, host, date_created FROM domains where host = ?"

	rows, err = conn.Query(sql, url.Host)
	if err != nil {
		log.Panic(err.Error())
	}
	defer rows.Close()

	var row Domain

	if rows.Next() {
		err := rows.Scan(&row.ID, &row.Host, &row.DateCreated)
		if err != nil {
			log.Panic(err)
		}
	}

	return row
}

func (d Domain) Insert(url url.URL) {

	// validate exists in redis
	var rows *sql.Rows
	conn := config.GetDBConnection("crawler")

	query := "INSERT IGNORE INTO domains (host) VALUES(?)"
	rows, err := conn.Query(query, url.Host)
	if err != nil {
		log.Panic(err.Error())
	}
	defer rows.Close()

}
