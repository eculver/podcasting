package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eculver/tdtv2/pkg/contentorigin"
	"github.com/segmentio/cli"
)

// TODO:
// - promptable fields
// - progress meter

type config struct {
	Host string `flag:"-H,--host" help:"Hostname"`
	Port string `flag:"-P,--port" help:"Port" default:"22"`
	User string `flag:"-u,--user" help:"Username"`
	Pass string `flag:"-p,--pass" help:"Password"`
}

func main() {
	cli.Exec(cli.Command(func(cfg config, src, dst string) {
		client, err := contentorigin.New(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		// open source file
		srcFile, err := os.Open(src)
		if err != nil {
			log.Fatal(err)
		}
		// upload file
		num, err := client.Put(srcFile, dst)
		if err != nil {
			fmt.Printf("got error uploading: %s\n", err)
		}
		fmt.Printf("%d bytes copied\n", num)
	}))
}
