package cmd 
import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/task/db"
)


var listCmd = &cobra.Command{
  Use:   "list",
  Short: "You can list items using this command",
  Run: func (cmd *cobra.Command, args []string){
  	//fmt.Println("hi, this is list command")

  	tasks, err := db.AllTasks()
  	if err != nil {
  		fmt.Println("Something went wrong, error :",err.Error())
  		os.Exit (1) 
  		return 
  	}
  	if (len(tasks) == 0 ){
  		fmt.Println("You have no tasks to complete")
  		return 
  	}
  	fmt.Println("You have the following tasks to complete")
  	for i , task := range(tasks){
  		fmt.Println( i+1,". ",task.Value," key = ",task.Key)
  	}
  	return 
  },

}



func init(){
	RootCmd.AddCommand(listCmd)
}