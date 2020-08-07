package main 

import (
	"github.com/task/cmd"
	"github.com/task/db"
)


func main (){

	dbPath:= "tasks.db"
	err := db.Init(dbPath)
	if err != nil{
		panic(err)
	}

	cmd.RootCmd.Execute()
	 
}

/*
func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// db, err := bolt.Open("my.db", 0600, nil)
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}



	defer db.Close()
}
*/