package db 
import (
	//"log"
	"encoding/binary"
	"time"
	"github.com/boltdb/bolt"

)


var taskBucket = []byte("tasks")
var completedtaskBucket = []byte("completedtasks")
var dateFlag = []byte("dateflag")
var db *bolt.DB

type Task struct {
	Key int 
	Value string 
}

func Init (dbPath string) error {
	var err error 
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err!= nil {
		return err 
	}

	return db.Update(func(tx *bolt.Tx) error {
		var err error 
		_, err = tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err 
		}
		_, err = tx.CreateBucketIfNotExists(completedtaskBucket)
		if err != nil {
			return err 
		}
		_, err = tx.CreateBucketIfNotExists(dateFlag)

		return err 	
	})
}
func SetTodayCompletedList()error {
	current_date := time.Now().String()[0:10]
  	//fmt.Println(current_date)
  	prevdate,err := DateFlagCheck(current_date)
  	if err!= nil {
  		return err 
  	}
  	if (prevdate != current_date){
  		err = AdjustCompletedBucket()
  	}
  	return err 
}
func DateFlagCheck (taskdate string)(string, error){
	id :=1
	date:= ""
	err := db.Update(func(tx *bolt.Tx)error{
		b:= tx.Bucket(dateFlag)
		key := itob(id)
		d := b.Get(key)
		if d!= nil {
			date = string(d)
		}
		if date != taskdate {
			return b.Put(key , []byte(taskdate))
		}
		return nil 
		
	})

	if err!= nil {
		return "" , err 
	}
	return date, nil 
}


func AdjustCompletedBucket () error {
	return db.Update(func(tx *bolt.Tx) error {
		var err error 
		err = tx.DeleteBucket(completedtaskBucket)
		if err != nil {
			return err 
		}
		_, err = tx.CreateBucketIfNotExists(completedtaskBucket)
		

		return err 	
	})
	
}


func CreateTask (task string)(int , error){
	var id int
	err := db.Update(func(tx *bolt.Tx)error{
		b:= tx.Bucket(taskBucket)
		id64,_ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key , []byte(task))
	})

	if err!= nil {
		return -1 , err 
	}
	return id, nil 
}

func CreateCompletedToday (task string)error{
	err := db.Update(func(tx *bolt.Tx)error{
		b:= tx.Bucket(completedtaskBucket)
		id64,_ := b.NextSequence()
		key := itob(int(id64))
		return b.Put(key , []byte(task))
	})
	return err 
}


func AllTasks () ([]Task,error){
	var tasks []Task 
	err := db.View(func(tx *bolt.Tx) error {
		b:= tx.Bucket(taskBucket)
		c := b.Cursor()
		for key, value:= c.First() ; key != nil; key,value= c.Next(){
			tasks = append (tasks, Task{
				Key : btoi(key),
				Value : string(value),
			})
		}
		return nil 
	})
	if err != nil {
		return nil,err
	}
	return tasks , nil 
}

func AllCompletedToday () ([]Task,error){
	var tasks []Task 
	err := db.View(func(tx *bolt.Tx) error {
		b:= tx.Bucket(completedtaskBucket)
		c := b.Cursor()
		for key, value:= c.First() ; key != nil; key,value= c.Next(){
			tasks = append (tasks, Task{
				Key : btoi(key),
				Value : string(value),
			})
		}
		return nil 
	})
	if err != nil {
		return nil,err
	}
	return tasks , nil 
}

func DeleteTask (key int )(error) {
	err:= db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		err := b.Delete(itob(key))
		if err != nil {
			return err
		}
		
		return err 
	})

	return err  
}



/*

func CompletedToday ()([]Task,error){
	var tasks []Task 
	err := db.View(func(tx *bolt.Tx) error {
		b:= tx.Bucket(completedtaskBucket)
		c := b.Cursor()
		for key, value:= c.First() ; key != nil; key,value= c.Next(){
			tasks = append (tasks, Task{
				Key : btoi(key),
				Value : string(value),
			})
		}
		return nil 
	})
	if err != nil {
		return nil,err
	}
	return tasks , nil 
}
*/


func itob (v int) []byte{
	b := make([]byte,8)
	binary.BigEndian.PutUint64(b,uint64(v))
	return b 
}

func btoi (b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

