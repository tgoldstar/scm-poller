package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/afero"
	"github.com/tgoldstar/scm-poller/pkg/config"
	"github.com/tgoldstar/scm-poller/pkg/poller"
	"github.com/tgoldstar/scm-poller/pkg/store"
	"github.com/urfave/cli/v2"
)

var (
	configFile     string
	configFileFlag = &cli.StringFlag{
		Name:        "config-file",
		Aliases:     []string{"f"},
		Destination: &configFile,
		Usage:       "YAML configuration file path",
		EnvVars:     []string{"SCM_POLLER_CONFIG_FILE"},
		Required:    true,
	}
)

func handleEvents(ctx context.Context, logger *log.Logger, changes <-chan poller.Change, errs <-chan poller.PollError) {
	for {
		select {
		case change := <-changes:
			log.Println(change)
		case err := <-errs:
			log.Println(err)
		case <-ctx.Done():
			return
		}
	}
}

func startPolling(c *cli.Context) error {
	opts, err := (&config.File{
		FileSystem: afero.NewOsFs(),
		Path:       configFile,
	}).Fetch()

	if err != nil {
		return err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	p := poller.New(logger, store.NewMemoryStore(), poller.Git{})
	changes := make(chan (poller.Change), 1000)
	errs := make(chan poller.PollError, 1000)
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go p.Poll(ctx, opts, changes, errs)
	go handleEvents(ctx, logger, changes, errs)
	<-sigChan
	fmt.Println("Exiting...")
	cancel()

	return nil
}

func main() {
	app := &cli.App{
		Name:  "scm-poller",
		Usage: "Poll remote SCM repositories for changes",
		Flags: []cli.Flag{
			configFileFlag,
		},
		Action: startPolling,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
