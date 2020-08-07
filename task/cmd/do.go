package cmd 
import (
	"github.com/spf13/cobra"
	"fmt"
	"strconv"
	"github.com/task/db"
	//"time"
)


var doCmd = &cobra.Command{
  Use:   "do",
  Short: "You can list items using this command",
  Run: func (cmd *cobra.Command, args []string){
  	//fmt.Println("hi, this is do command")
  	var ids  []int 
  	for _,record := range args{
  		id,err := strconv.Atoi(record)
  		if err != nil{
  			fmt.Println("Error")
  		}else {
  			ids = append(ids, id)
  		}
  	}

  	tasks, err := db.AllTasks()
  	if err != nil {
  		fmt.Println("There was an error retrieving data, error: ",err.Error()) 
  	}
  	for _,id := range (ids){
  		if (id <=0 || id >len(tasks)){
  			fmt.Println ("Invalid id :",id)
  			continue 
  		}

  		err := db.DeleteTask(tasks[id-1].Key)
  		if err != nil {
  			fmt.Println("An error occured while deleting the task with id: ", id)
  		}else {
  			err = db.SetTodayCompletedList()
  			if err!=nil {
  				fmt.Println("Error occured in setting today's completed list, error: ",err.Error())
  			}

  			err = db.CreateCompletedToday(tasks[id-1].Value)
  			fmt.Println("Task marked as completed")
  		}
  	}


  	//fmt.Println(ids)
  },

}



func init(){
	RootCmd.AddCommand(doCmd)
}