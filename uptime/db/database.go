package db

import (
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // if required
	"github.com/uptime/ut"
)

var m sync.Mutex

// Dbquery : all databse queries are done using Dbquery function only
func Dbquery(action string, query string) ([]ut.Urltable, error) {
	var s []ut.Urltable
	m.Lock()
	dbs, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/uptimeurls")
	if err != nil {
		log.Println("Connection Failed to Open, Error: " + err.Error())
		m.Unlock()
		return nil, err

	}
	log.Println("Connection Established")

	if action == "raw" {
		dbs.Raw(query).Scan(&s)
	} else {
		dbs.Exec(query)
	}

	defer dbs.Close()
	m.Unlock()
	if err != nil {
		return nil, err
	}
	return s, nil

}
