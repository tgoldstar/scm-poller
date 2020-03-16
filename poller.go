package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var commitHistroy map[string]plumbing.Hash = make(map[string]plumbing.Hash)

func poll(remote string, tracking []string) error {
	trackingSet := make(map[string]struct{})
	for _, branch := range tracking {
		trackingSet[branch] = struct{}{}
	}

	log.Print("Initializing cache...")

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return err
	}

	for _, ref := range refs {
		name := ref.Name().Short()
		hash := ref.Hash()

		if _, ok := trackingSet[name]; ok {
			log.Printf("Setting %s to %v", name, hash)
			commitHistroy[name] = hash
		}
	}

	log.Println("Polling remote for changes...")
	pollTicker := time.NewTicker(1 * time.Second)

	for _ = range pollTicker.C {
		refs, err := rem.List(&git.ListOptions{})
		if err != nil {
			return err
		}
		for _, ref := range refs {
			name := ref.Name().Short()
			current := ref.Hash()

			if previous, ok := commitHistroy[name]; ok {
				if current != previous {
					log.Printf("Branch %s was changed from %v to %v", name, previous, current)
					commitHistroy[name] = current
				}
			}
		}
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:                 "scm-poller",
		EnableBashCompletion: true,
		BashComplete:         cli.ShowCompletions,
		Usage:                "Poll remote SCM for changes",
		ArgsUsage:            "<remote> <branch>",
		Action: func(c *cli.Context) error {
			remote := c.Args().Get(0)
			branch := c.Args().Get(1)
			if remote == "" || branch == "" {
				cli.ShowAppHelpAndExit(c, 1)
			}
			return poll(remote, []string{branch})
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
