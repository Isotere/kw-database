package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Isotere/kw-database/apps/server/app/commands/start"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Server",
	Long:  `Longer Description Will be Later`,
	Run: func(cmd *cobra.Command, _ []string) {
		hostFlag := cmd.Flags().Lookup("host")
		if hostFlag == nil {
			fmt.Println("Host is required")
			os.Exit(1)
		}

		portFlag := cmd.Flags().Lookup("port")
		if portFlag == nil {
			fmt.Println("Port is required")
			os.Exit(1)
		}

		host := hostFlag.Value.String()
		port, _ := strconv.Atoi(portFlag.Value.String())

		start.Handle(host, port)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("host", "H", defaultHost, "server host")
	startCmd.Flags().Int32P("port", "p", defaultPort, "server port")
}
