package cmd 

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
  Use:   "task",
  Short: "This is a root command",	
}

/*
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
*/