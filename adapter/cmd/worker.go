package cmd

import (
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/session"
	"github.com/Miyagawa-Ryohei/mkmicro/app"
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
	"github.com/Miyagawa-Ryohei/mkmicro/infra"
	"github.com/spf13/cobra"
	"log"
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
	config := new(entity.Config)
	configLoader.Load(config)
	factory := session.NewSTSManagerFactory(config.Queue, config.Storage)
	sess,err := factory.Create()
	if err != nil {
		log.Fatal(err)
	}
	return app.NewSubscriber(sess,factory)
}

func init() {
	RootCmd.AddCommand(WorkerCmd.Command)
}
