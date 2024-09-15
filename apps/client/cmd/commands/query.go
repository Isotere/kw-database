package commands

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Isotere/kw-database/apps/client/cmd/commands/query"
)

// queryCmd represents the start command
var queryCmd = &cobra.Command{
	Use:   "query \"<quoted string>\"",
	Short: "Make Query",
	Long:  `Longer Description Will be Later`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Usage()
			os.Exit(1)
		}

		hostFlag := cmd.Flags().Lookup("host")

		portFlag := cmd.Flags().Lookup("port")

		host := ""
		if hostFlag != nil {
			host = hostFlag.Value.String()
		}

		port, _ := strconv.Atoi(portFlag.Value.String())

		query.Handle(args[0], host, port)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringP("host", "H", "", "server host")
	_ = queryCmd.MarkFlagRequired("host")
	queryCmd.Flags().Int32P("port", "p", 0, "server port")
	_ = queryCmd.MarkFlagRequired("port")
}
