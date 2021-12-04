package cmd

import (
	"github.com/spf13/cobra"
	"mkmicro/adapter/gateway/session"
	"mkmicro/app"
	"mkmicro/entity"
	"mkmicro/infra"
)

var WorkerCmd = entity.WorkerCommand{
	Command: &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, arg []string) error {
			subscriber := GetSubscriber()
			subscriber.Listen()
			return nil
		},
	},
}

func GetSubscriber () *app.Subscriber {
	configLoader := infra.ConfigLoader{}
	config := new(infra.Config)
	configLoader.Load(config)
	factory := session.NewSTSManagerFactory(config.Queue, config.Session)
	sess := factory.Create()
	return app.NewSubscriber(sess,factory)
}

func init() {
	RootCmd.AddCommand(WorkerCmd.Command)
}
