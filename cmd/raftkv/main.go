// Command raftkv runs one node of the raft-kv replicated key-value store.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	cfg, err := parseConfig(os.Args[1:], os.Getenv, os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "raftkv:", err)
		os.Exit(2)
	}

	fmt.Printf("id:        %d\n", cfg.ID)
	fmt.Printf("data-dir:  %s\n", cfg.DataDir)
	fmt.Printf("bootstrap: %t\n", cfg.Bootstrap)
	fmt.Println("peers:")

	for _, p := range cfg.Peers {
		fmt.Printf("  %d = %s\n", p.ID, p.Addr)
	}
}
