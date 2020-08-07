package cmd 
import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
	"os"
	"github.com/task/db"
)


var addcmd = &cobra.Command{
  Use:   "add",
  Short: "You can add items using this command",
  Run: func (cmd *cobra.Command, args []string){
  	//fmt.Println("hi, this is add command")
  	addTask := strings.Join(args," ")
  	_, err := db.CreateTask(addTask)
  	if err != nil {
  		fmt.Println("Something went wrong, error :",err.Error())
  		os.Exit (1) 
  		return 
  	}
  	fmt.Println("task add ")
  	return 
  },

}



func init(){
	RootCmd.AddCommand(addcmd)
}