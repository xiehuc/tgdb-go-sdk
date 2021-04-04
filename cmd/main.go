//
// main.go
// Copyright (C) 2018 toraxie <toraxie@tencent.com>
//
// Distributed under terms of the Tencent license.
//

package main

import (
	"bufio"
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
)

func read(cmds chan string) {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		cmds <- line
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

	log.Println("connection successful")

	cmds := make(chan string, 24)
	fmt.Print("cypher> ")
	go read(cmds)

	for line := range cmds {
		if len(line) == 0 {
			fmt.Fprint(os.Stdout, "\n")
			continue
		}
		all, m, err := s.Cypher(line, nil)
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		} else {
			fmt.Fprint(os.Stdout, all, m, "\n")
		}
		fmt.Print("cypher> ")
	}
}
