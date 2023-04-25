package cmd

import (
	"github.com/souravbiswassanto/api-bookserver/apiHandler"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var (
	// Port stores port number for starting a connection
	Port     int
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start cmd starts the server on a port",
		Long: `It starts the server on a given port number, 
				Port number will be given in the cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			apiHandler.RunServer(Port)
		},
	}
)

func init() {

	startCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8081, "Port number for starting server")
	rootCmd.AddCommand(startCmd)
}
