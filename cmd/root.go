package cmd

import (
	"github.com/omerkaya1/watcher/internal/config"
	"github.com/omerkaya1/watcher/internal/db"
	"github.com/omerkaya1/watcher/internal/errors"
	"github.com/omerkaya1/watcher/internal/mq"
	"github.com/spf13/cobra"
	"log"
)

var configFile string

var RootCmd = &cobra.Command{
	Use:     "watcher",
	Short:   "simple DB watcher service that queries a DB with calendar events and enqueues them in a message queue",
	Example: "  watcher -c /path/to/config.json",
	Run:     startNotificationService,
}

func init() {
	RootCmd.Flags().StringVarP(&configFile, "config", "c", "", "specifies the path to a configuration file")
}

func startNotificationService(cmd *cobra.Command, args []string) {
	if configFile == "" {
		log.Fatalf("%s:%s", errors.ErrCMDPrefix, errors.ErrBadConfigFile)
	}

	conf, err := config.InitConfig(configFile)
	if err != nil {
		log.Fatalf("%s:%s", errors.ErrCMDPrefix, err)
	}

	esp, err := db.NewMainEventStorage(conf.DB)
	if err != nil {
		log.Fatalf("%s:%s", errors.ErrDBPrefix, err)
	}

	messageQueue, err := mq.NewEventMQProducer(conf.Queue, esp)
	if err != nil {
		log.Fatalf("%s:%s", errors.ErrMQPrefix, err)
	}
	log.Println("Watcher service initialisation")
	if err := messageQueue.ProduceMessages(); err != nil {
		log.Fatalf("%s:%s", errors.ErrMQPrefix, err)
	}
}
