package cmd 
import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/task/db"
)


var completeCmd = &cobra.Command{
  Use:   "completed",
  Short: "You can list items using this command",
  Run: func (cmd *cobra.Command, args []string){
  	//fmt.Println("hi, this is list command")
    errr := db.SetTodayCompletedList()
    if errr!=nil {
      fmt.Println("Error occured in setting today's completed list, error: ",errr.Error())
    }
  	tasks, err := db.AllCompletedToday()
  	if err != nil {
  		fmt.Println("Something went wrong, error :",err.Error())
  		os.Exit (1) 
  		return 
  	}
  	if (len(tasks) == 0 ){
  		fmt.Println("You have no completed tasks")
  		return 
  	}
  	fmt.Println("You have completed these tasks")
  	for i , task := range(tasks){
  		fmt.Println( i+1,". ",task.Value)
  	}
  	return 
  },

}



func init(){
	RootCmd.AddCommand(completeCmd)
}