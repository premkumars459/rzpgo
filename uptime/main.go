package main

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/uptime/db"
	"github.com/uptime/uc"
	"github.com/uptime/urlhandler"
	"github.com/uptime/ut"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Hello world")
	urlChannel := make(map[string]chan ut.Urltable)
	deactivateChannel := make(map[int]chan int)
	//var m sync.Mutex
	r := gin.Default()
	r.GET("/urls/:id", urlhandler.GetHomePage)
	r.POST("/urls/", urlhandler.PostHomepage(urlChannel, deactivateChannel))
	r.PATCH("/urls/", urlhandler.PatchHomepage(urlChannel))
	r.POST("/urls/:id/activate", urlhandler.ActivateURL(urlChannel, deactivateChannel))
	r.POST("/urls/:id/deactivate", urlhandler.DeactivateURL(deactivateChannel))
	r.DELETE("/urls/:id", urlhandler.DeleteURL(deactivateChannel))

	go urlcheck(urlChannel, deactivateChannel)

	r.Run()
}

func urlcheck(urlChannel map[string]chan ut.Urltable, deactivateChannel map[int]chan int) {
	dbc, errr := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/uptimeurls")
	if errr != nil {
		log.Println("Connection Failed to Open, Error: " + errr.Error())
	}
	log.Println("Connection Established")
	if dbc.HasTable("urltable") == false {
		dbc.Table("urltable").CreateTable(&ut.Urltable{})
	}
	
	defer dbc.Close()

	sqlQuery := "SELECT * FROM urltable"
	records, err := db.Dbquery("raw", sqlQuery)
	check(err)
	//db.Table("urltable").Find(&records)
	if len(records) > 0 {
		for _, record := range records {
			if record.Status == "active" {
				urlChannel[record.Url] = make(chan ut.Urltable, 1)
				deactivateChannel[record.Id] = make(chan int, 1)
				go uc.CheckURL(record.Url, record.CrawlTimeout, record.Frequency, record.FailureCount, record.FailureThreshold, urlChannel[record.Url], deactivateChannel[record.Id])
			}
		}
	} else {
		log.Println("you don't have any urls to check.... try posting some")
	}

}
