package urlhandler

import (
	"log"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql" // if required
	"github.com/uptime/db"
	"github.com/uptime/uc"
	"github.com/uptime/ut"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//GetHomePage : Retrieves the url data from database.
func GetHomePage(c *gin.Context) {
	id := c.Param("id")

	sqlQuery := "SELECT * FROM urltable WHERE id = " + id
	records, err := db.Dbquery("raw", sqlQuery)
	check(err)

	if len(records) > 0 {
		c.JSON(200, gin.H{
			"message":           "Hi GET!",
			"requested id":      records[0].Id,
			"url found":         records[0].Url,
			"crawl_timeout":     records[0].CrawlTimeout,
			"frequency":         records[0].Frequency,        // every 30 seconds
			"failure_threshold": records[0].FailureThreshold, // mark as inactive once failure count reaches 50
			"status":            records[0].Status,           // active or inactive
			"failure_count":     records[0].FailureCount,
		})
	} else {
		c.JSON(200, gin.H{
			"message":      "Hi from get, No such id found",
			"requested id": id,
		})
	}
}

//isValidURL : checks whether the requested url is in the correct format.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// PostHomepage : Recieves the url to be freshly started and process it.
func PostHomepage(urlChannel map[string]chan ut.Urltable, deactivateChannel map[int]chan int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		url, err := c.GetPostForm("url")
		msg := "no url"
		crawlTimeout := c.DefaultPostForm("crawl_timeout", "20")
		frequency := c.DefaultPostForm("frequency", "30")
		failureThreshold := c.DefaultPostForm("failure_threshold", "50")
		if err != false {
			if isValidURL(url) == false {
				c.JSON(200, gin.H{
					"url":     url,
					"message": "check your url... maintain format as in example shown",
					"example": "http://example.com/",
				})
			} else {
				msg = url

				sqlQuery := "SELECT * FROM urltable WHERE url = '" + url + "'"
				s, err := db.Dbquery("raw", sqlQuery)
				check(err)

				if len(s) == 0 {
					var ct, f, ft int
					var err error
					ct, err = strconv.Atoi(crawlTimeout)
					check(err)
					f, err = strconv.Atoi(frequency)
					check(err)
					ft, err = strconv.Atoi(failureThreshold)
					check(err)

					sqlQuery = "INSERT INTO urltable (url , crawl_timeout,frequency, failure_threshold,status,failure_count ) VALUES ('" + url + "', " + crawlTimeout + "," + frequency + ", " + failureThreshold + ",'active', 0);"
					_, err = db.Dbquery("exec", sqlQuery)
					check(err)

					sqlQuery = "SELECT * FROM urltable WHERE url = '" + url + "'"
					s, err = db.Dbquery("raw", sqlQuery)
					check(err)

					urlChannel[url] = make(chan ut.Urltable, 1)
					deactivateChannel[s[0].Id] = make(chan int, 1)

					go uc.CheckURL(url, ct, f, 0, ft, urlChannel[url], deactivateChannel[s[0].Id])

					c.JSON(200, gin.H{
						"message":           msg,
						"id":                "                        b7f32a21-b863-4dd1-bd86-e99e8961ffc6",
						"url":               "abc.com",
						"crawl_timeout":     crawlTimeout,
						"frequency":         frequency,        // every 30 seconds
						"failure_threshold": failureThreshold, // mark as inactive once failure count reaches 50
						"status":            "active",         // active or inactive
						"failure_count":     0,
					})
				} else {
					c.JSON(200, gin.H{
						"message": "url already exist",
						"url":     url,
					})
				}
			}

		} else {
			c.JSON(200, gin.H{
				"message": "Didn't recieve the request",
				"url":     url,
			})

		}
	}
	return gin.HandlerFunc(fn)

}
func getPatchDetails(url string, crawlTimeout string, frequency string, failureThreshold string) (string, ut.Urltable, error) {
	var k ut.Urltable
	var err error
	k.Id = -1
	query := "UPDATE urltable SET "
	if crawlTimeout != "" {
		query += "crawl_timeout = " + crawlTimeout + ", "
		k.CrawlTimeout, err = strconv.Atoi(crawlTimeout)
		if err != nil {
			return "", k, err
		}
	} else {
		k.CrawlTimeout = -1
	}
	if frequency != "" {
		query += "frequency = " + frequency + ", "
		k.Frequency, err = strconv.Atoi(frequency)
		if err != nil {
			return "", k, err
		}
	} else {
		k.Frequency = -1
	}
	if failureThreshold != "" {
		query += "failure_threshold = " + failureThreshold + ", "
		k.FailureThreshold, err = strconv.Atoi(failureThreshold)
		if err != nil {
			return "", k, err
		}
	} else {
		k.FailureThreshold = -1
	}
	query += "failure_count = 0 WHERE url = '" + url + "'"
	k.FailureCount = 0
	k.Status = ""
	k.Url = url

	return query, k, nil

}

// PatchHomepage : this function is used to update the checking factors of a url.
func PatchHomepage(urlChannel map[string]chan ut.Urltable) gin.HandlerFunc {

	fn := func(c *gin.Context) {

		url, err := c.GetPostForm("url")
		msg := "error in getting url"

		if err != false {
			msg = "got url"

			crawlTimeout := c.DefaultPostForm("crawl_timeout", "")
			frequency := c.DefaultPostForm("frequency", "")
			failureThreshold := c.DefaultPostForm("failure_threshold", "")
			sqlQuery, k, errPatch := getPatchDetails(url, crawlTimeout, frequency, failureThreshold)

			if isValidURL(url) == true && errPatch == nil {
				sql := "SELECT * FROM urltable WHERE url = '" + url + "'"
				s, errdb := db.Dbquery("raw", sql)
				check(errdb)
				if len(s) > 0 {
					if s[0].Status == "active" {
						urlChannel[k.Url] <- k
					}
					_, errdb = db.Dbquery("exec", sqlQuery)
					check(errdb)

					c.JSON(200, gin.H{
						"message":           msg,
						"id":                k.Id,
						"url":               k.Url,
						"crawl_timeout":     k.CrawlTimeout,
						"frequency":         k.Frequency,        // every 30 seconds
						"failure_threshold": k.FailureThreshold, // mark as inactive once failure count reaches 50
						"status":            k.Status,           // active or inactive
						"failure_count":     k.FailureCount,
					})
				}
			} else {
				if isValidURL(url) == false {
					msg = "invalid url"
				}
				if errPatch != nil {
					msg = "invalid details"
				}
				if isValidURL(url) == false && errPatch != nil {
					msg = "invalid urls and details"
				}
				c.JSON(200, gin.H{
					"message": msg,
					"url":     url,
				})
			}
		} else {
			c.JSON(200, gin.H{
				"message": msg,
				"url":     url,
			})

		}
	}
	return gin.HandlerFunc(fn)
}

//ActivateURL : activate an inactive url.
func ActivateURL(urlChannel map[string]chan ut.Urltable, deactivateChannel map[int]chan int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		sqlQuery := "SELECT * FROM urltable WHERE id = " + id
		records, err := db.Dbquery("raw", sqlQuery)
		check(err)
		msg := "Hi from Post activate, No such id found"
		//db.Raw(sql).Scan(&records)
		if len(records) > 0 {
			if records[0].Status == "active" {
				msg = "error -> Url already active"
			} else {
				sqlQuery := "UPDATE urltable SET status = 'active' WHERE id = " + id
				_, errdb := db.Dbquery("exec", sqlQuery)
				check(errdb)
				sqlQuery = "SELECT * FROM urltable WHERE id = " + id
				records, err = db.Dbquery("raw", sqlQuery)
				check(err)
				go uc.CheckURL(records[0].Url, records[0].CrawlTimeout, records[0].Frequency, records[0].FailureCount, records[0].FailureThreshold, urlChannel[records[0].Url], deactivateChannel[records[0].Id])

				msg = "url activated"
			}

		}
		c.JSON(200, gin.H{
			"message":      msg,
			"requested id": id,
		})
	}
	return gin.HandlerFunc(fn)
}

// DeactivateURL : deactivate an active url
func DeactivateURL(deactivateChannel map[int]chan int) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		id := c.Param("id")
		log.Println("deactivate called")
		msg := "Hi from Post deactivate, No such id found"
		sqlQuery := "SELECT * FROM urltable WHERE id = " + id
		records, err := db.Dbquery("raw", sqlQuery)
		check(err)

		if len(records) > 0 {
			if records[0].Status == "inactive" {
				msg = "error -> Url already inactive"
			} else {
				idint, _ := strconv.Atoi(id)
				deactivateChannel[idint] <- 1
				sqlQuery := "UPDATE urltable SET status = 'inactive' WHERE id = " + id
				_, errdb := db.Dbquery("exec", sqlQuery)
				check(errdb)
				msg = "url deactivated"
			}

		}
		c.JSON(200, gin.H{
			"message":      msg,
			"requested id": id,
		})
	}
	return gin.HandlerFunc(fn)
}

//DeleteURL : it will remove the url from the database and stop checking for it.
func DeleteURL(deactivateChannel map[int]chan int) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		id := c.Param("id")
		log.Println("delete called")
		msg := "Hi from delete, No such id found"
		sqlQuery := "SELECT * FROM urltable WHERE id = " + id
		records, err := db.Dbquery("raw", sqlQuery)
		check(err)

		if len(records) > 0 {

			idint, _ := strconv.Atoi(id)
			deactivateChannel[idint] <- 1
			sqlQuery := "DELETE FROM urltable WHERE id = " + id
			_, errdb := db.Dbquery("exec", sqlQuery)
			check(errdb)
			//idint = 1 //<-deactivateChannel[idint]
			msg = "url deleted"

		}
		c.JSON(200, gin.H{
			"message":      msg,
			"requested id": id,
		})
	}
	return gin.HandlerFunc(fn)
}
