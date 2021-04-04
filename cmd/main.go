//
// main.go
// Copyright (C) 2018 toraxie <toraxie@tencent.com>
//
// Distributed under terms of the Tencent license.
//

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"git.code.oa.com/opdd/tgdb-go-sdk"
)

var (
	user = flag.String("u", "", "user")
	pass = flag.String("p", "", "password")
	host = flag.String("h", "", "host")
	port = flag.Int("P", 7687, "Port")
	db   = flag.String("db", "", "database name")

	interactive bool
)

func init() {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		interactive = true
	}
}

func read(cmds chan string) {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		cmds <- line
	}
	close(cmds)
}

func prompt() {
	if interactive {
		fmt.Print("cypher> ")
	}
}

func main() {
	flag.Parse()

	client, err := tgdb.New(*host, *port, *user, *pass)
	if err != nil {
		log.Fatal(err)
	}
	s, err := client.Open(*db, tgdb.ReadMode)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	if interactive {
		log.Println("connection successful")
	}

	cmds := make(chan string, 24)
	prompt()
	enc := json.NewEncoder(os.Stdout)
	go read(cmds)

	for line := range cmds {
		if len(line) == 0 {
			fmt.Fprint(os.Stdout, "\n")
			continue
		}
		all, _, err := s.Cypher(line, nil)
		if err != nil {
			log.Println(err)
		} else {
			_ = enc.Encode(all)
		}
		prompt()
	}

}
