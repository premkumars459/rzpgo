
package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	//"io"
	"log"
	"os"
	"strings"
	"math/rand"
	"flag"
	"time"
)
func getanswer (anschannel chan string ){
	var ans string 
	fmt.Scanln(&ans)
	anschannel<-ans 
}

func shuffleques(records [][]string){
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) {
    	records[i], records[j] = records[j], records[i]
	})
	return 
}

func getCsvfile ( filename string) ([][]string, error){
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))
	
	// Iterate through the records
    records,err := r.ReadAll()
    

	return records ,err
}

func quiz (timerflag time.Duration, records [][]string ){
	score :=0
	fmt.Printf("Your time starts now !! \n")

	mytimer := time.NewTimer(timerflag)
	//anschannel := make(chan string )
	timeflag := 0

	for idx, record := range records {

		fmt.Printf("Question %d : %s Answer: \n", idx ,record[0])
		anschannel := make(chan string )
		go getanswer (anschannel)


		select {
		case <-mytimer.C :
			fmt.Printf("\n Your time is up \n")
			timeflag = 1 
			break 

		case ans:= <-anschannel :
			if ans == strings.Trim(record[1]," "){
			score+=1 ;
		}
		}
		if timeflag == 1 {
			break ;
		}
	}
	fmt.Printf("Congratulations, Your final score is %d" , score)
}



func main() {
	// Open the file

	var timerFlag = flag.Duration("timer", 10*time.Second, "Flag to set test duration. Input format : `<time>s`(without quotes)")
	var testFile = flag.String("test", "problems.csv", "File name of the test set.")
	var shuffle = flag.Bool("shuffle", true, "Boolean flag to shuffle the test")

	records,err := getCsvfile(*testFile)
	if err != nil {
			log.Fatal(err)
	}

    if *shuffle == true {
    	shuffleques(records)
    	
    }

    quiz(*timerFlag, records)

	//for {
		// Read each record from csv
	//	record, err := r.Read()
	//	if err == io.EOF {
	//		break
	//	}
	//}
	
}