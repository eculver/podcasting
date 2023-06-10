package main

import (
	"fmt"
	"log"

	"github.com/eculver/tdtv2/pkg/contentorigin"
	"github.com/segmentio/cli"
)

type config struct {
	Host string `flag:"-H,--host" help:"Hostname"`
	Port string `flag:"-P,--port" help:"Port" default:"22"`
	User string `flag:"-u,--user" help:"Username"`
	Pass string `flag:"-p,--pass" help:"Password"`
}

func main() {
	cli.Exec(cli.Command(func(cfg config, root string) {
		client, err := contentorigin.New(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		// list files
		if err := client.Walk(root, func(path string) error {
			fmt.Println(path)
			return nil
		}); err != nil {
			fmt.Printf("got error walking: %s\n", err)
		}
	}))
}
