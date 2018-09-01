# Broadcast

[![CircleCI](https://circleci.com/gh/solcates/broadcast/tree/master.svg?style=svg)](https://circleci.com/gh/solcates/broadcast/tree/master) [![codecov](https://codecov.io/gh/solcates/broadcast/branch/master/graph/badge.svg)](https://codecov.io/gh/solcates/broadcast) [![Go Report Card](https://goreportcard.com/badge/github.com/solcates/broadcast)](https://goreportcard.com/report/github.com/solcates/broadcast)

**broadcast** is a go module for discovering devices/services on a LAN using UDP4 broadcasting(not multicasting).

## Requirements

- Go 1.11+ for go module support

## Usage

### library

To use **broadcast** in your Go application:

```
package main

import (
	"github.com/solcates/broadcast"
	"log"
)

func main() {

	// Create a broadcaster instance
	bc := broadcast.NewBroadcaster(8080, "Is there anybody out there?")

	// to include this node as well, set findself to true
	bc.SetFindself(true)

	// discover the nodes on this LAN
	nodes, err := bc.Discover()
	if err != nil {
		log.Fatal(err)
	}
	// iterate over the nodes found.
	for _, node := range nodes {
		log.Printf("Found Node @ %v", node)
	}
}

```

### binary

To install broadcast the cli...

`go get -u github.com/solcates/broadcast/broadcast`
