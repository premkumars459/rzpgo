package uc

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/uptime/db"
	"github.com/uptime/ut"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//CheckURL : urls are sent into go routines and are checked for status in an infinite loop. these go routines can be interupted by update or deactivate commands.
func CheckURL(url string, crawlTime int, frequency int, failureCount int, failureThreshold int, updateChannel chan ut.Urltable, deactivateChannel chan int) {
	failedCases := failureCount
	for {
		if failedCases >= failureThreshold {
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
			failedCases++
			log.Println("inactive " + url)
			go markFailure(url, failedCases)

			//log.Println(err.Error(), resp)
		} else {
			//fmt.Println("no error")
			log.Println("active ******************************************************************************" + url)

			log.Println(resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		flag := 0
		for {
			select {
			case k := <-updateChannel:
				log.Println("update recieved for " + url)
				if k.CrawlTimeout != -1 {
					crawlTime = k.CrawlTimeout
				}
				if k.Frequency != -1 {
					frequency = k.Frequency
				}
				if k.FailureThreshold != -1 {
					failureThreshold = k.FailureThreshold
				}
				log.Println("update recieved for " + url)
				break
			case <-mytimer.C:
				flag = 1
			}
			if flag == 1 {
				break
			}

		}
		select {
		case _ = <-deactivateChannel:
			failedCases = failureThreshold
			log.Println("deactivate called ---------------------------------------------------------")

		default:
		}

	}

	log.Println("Threshold reached...")
}

func markFailure(url string, failedCases int) {
	failureCount := strconv.Itoa(failedCases)
	sqlQuery := "UPDATE urltable SET failure_count = " + failureCount + " WHERE url ='" + url + "'"
	_, errdb := db.Dbquery("exec", sqlQuery)
	check(errdb)

}
