package cmd

import (
	"github.com/spf13/cobra"
	"mkmicro/app"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, arg []string) error {
		s := GetServer()
		s.Listen()
		return nil
	},
}

func GetServer () *app.Server {
	s := app.NewServer()
	s.Init()
	return s
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
