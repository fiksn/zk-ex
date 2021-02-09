package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	zk "github.com/fiksn/zk"
	cli "github.com/urfave/cli/v2"
)

var (
	GitRevision = "unknown"
)

// TODO: lame handling of quotes, make sure \" is properly interpreted too eventually
func split(s string) []string {
	sq := false
	dq := false

	a := strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' && !sq {
			dq = !dq
		}
		if r == '\'' && !dq {
			sq = !sq
		}
		return !dq && !sq && r == ' '
	})

	return a
}

func quoted(s string) bool {
	return (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"')
}

func unquote(a *[]string) {
	for i, s := range *a {
		if len(s) > 0 && quoted(s) {
			(*a)[i] = s[1 : len(s)-1]
		}
	}
}

func execute(command string, args []string) error {
	unquote(&args)
	cmd := exec.Command(command, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func main() {
	app := &cli.App{
		Name:  fmt.Sprintf("zk-ex (%s)", GitRevision),
		Usage: "Acquire an exclusive lock via zookeeper",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "zkservers",
				Value:   "127.0.0.1:2181",
				Usage:   "Comma separated list of zookeeper servers",
				EnvVars: []string{"ZKSERVERS"},
			},
			&cli.StringFlag{
				Name:    "zkpath",
				Value:   "/lock",
				Usage:   "Absolute zookeeper path for the lock",
				EnvVars: []string{"ZKLOCKPATH"},
			},
			&cli.StringFlag{
				Name:  "lock",
				Value: "",
				Usage: "Action to do when lock acquired (will block until obtaining exclusive lock unless -nolock is present)",
			},
			&cli.StringFlag{
				Name:  "nolock",
				Value: "",
				Usage: "Action to do when lock NOT acquired (when this is present there will be no blocking)",
			},
		},

		Action: func(c *cli.Context) error {
			servers := strings.Split(c.String("zkservers"), ",")
			if len(servers) < 1 {
				log.Fatalf("Invalid zkservers")
			}
			zkpath := c.String("zkpath")
			if len(zkpath) == 0 || zkpath[0] != '/' {
				log.Fatalf("Invalid zkpath")
			}

			lock := split(c.String("lock"))
			nolock := split(c.String("nolock"))
			// The logic is if nolock is present then just try to acquire lock
			try := len(nolock) > 0 && len(nolock[0]) > 0

			conn, _, err := zk.Connect(servers, 5*time.Second)
			if err != nil {
				log.Fatalf("Error %v\n", err)
			}
			defer conn.Close()

			acls := zk.WorldACL(zk.PermAll)
			l := zk.NewLock(conn, zkpath, acls)

			if try {
				ret, err := l.TryLock()
				if err != nil {
					log.Fatalf("Try lock error %v\n", err)
				}
				if ret {
					execute(lock[0], lock[1:])
					l.Unlock()
				} else {
					execute(nolock[0], nolock[1:])
				}
			} else {
				err := l.Lock()
				if err != nil {
					log.Fatalf("Lock error %v\n", err)
				}
				execute(lock[0], lock[1:])
				l.Unlock()
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
