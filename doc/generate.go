package main

import (
	"log"

	"github.com/davinash/yados/cmd/cli/commands/store"

	"github.com/davinash/yados/cmd/cli/commands/server"
	"github.com/spf13/cobra/doc"
)

func main() {
	serverCommands := server.NewServerCommands()
	err := doc.GenMarkdownTree(serverCommands, "./")
	if err != nil {
		log.Fatal(err)
	}
	storeCommands := store.NewStoreCommands()
	err = doc.GenMarkdownTree(storeCommands, "./")
	if err != nil {
		log.Fatal(err)
	}
}
