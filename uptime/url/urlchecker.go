package uc

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/uptime/db"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Urltable struct {
	Id                int    `gorm:"AUTO_INCREMENT"`
	Url               string `gorm:"unique;not null"`
	Crawl_timeout     int    `gorm:"not null"`
	Frequency         int    `gorm:"not null"`
	Failure_threshold int    `gorm:"not null"`
	Status            string `gorm:"not null"`
	Failure_count     int    `gorm:"not null"`
}

func CheckUrl(url string, crawlTime int, frequency int, failure_count int, failure_threshold int, updateChannel chan Urltable, deactivateChannel chan int) {
	failedCases := failure_count
	for {
		if failedCases >= failure_threshold {
			sqlQuery := "UPDATE urltable SET status = 'inactive' WHERE url ='" + url + "' ;"
			_, errdb := db.Dbquery("exec", sqlQuery)
			check(errdb)
			break
		}

		mytimer := time.NewTimer(time.Duration(frequency) * time.Second)
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		crawlTime = 1
		//fmt.Println("starting")
		resp, err := client.Get(url)
		if err != nil {
			//fmt.Println("error")
			failedCases += 1
			log.Println("inactive " + url)
			go markFailure(url, failedCases)

			//log.Println(err.Error(), resp)
		} else {
			//fmt.Println("no error")
			log.Println("active ******************************************************************************" + url)

			log.Println(resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		select {
		case k := <-updateChannel:
			log.Println("update recieved for " + url)
			if k.Crawl_timeout != -1 {
				crawlTime = k.Crawl_timeout
			}
			if k.Frequency != -1 {
				frequency = k.Frequency
			}
			log.Println("update recieved for " + url)
			break
		case <-mytimer.C:

		}
		select {
		case _ = <-deactivateChannel:
			failedCases = failure_threshold
			log.Println("deactivate called ---------------------------------------------------------")

		default:
		}

	}

	log.Println("Threshold reached...")
}

func markFailure(url string, failedCases int) {
	failure_count := strconv.Itoa(failedCases)
	sqlQuery := "UPDATE urltable SET failure_count = " + failure_count + " WHERE url ='" + url + "'"
	_, errdb := db.Dbquery("exec", sqlQuery)
	check(errdb)

}
