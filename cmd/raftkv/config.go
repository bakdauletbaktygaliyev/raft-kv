package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

// Peer is one cluster member: its node ID and the address other nodes dial.
type Peer struct {
	ID   uint64
	Addr string
}

// Config is the startup configuration of one raftkv node.
type Config struct {
	ID        uint64
	Peers     []Peer // sorted by ID, includes this node
	DataDir   string
	Bootstrap bool
}

// parseConfig reads flags from args, falling back to environment variables
// (via getenv) for anything not given on the command line.
// Precedence: flag > env > default.
func parseConfig(args []string, getenv func(string) string, out io.Writer) (Config, error) {
	fs := flag.NewFlagSet("raftkv", flag.ContinueOnError)
	fs.SetOutput(out)

	id := fs.Uint64("id", envUint(getenv, "RAFTKV_ID"), "this node's ID (env RAFTKV_ID)")
	peers := fs.String("peers", getenv("RAFTKV_PEERS"), "all cluster members as id=host:port,... (env RAFTKV_PEERS)")
	dataDir := fs.String("data-dir", envOr(getenv, "RAFTKV_DATA_DIR", "data"), "directory for the WAL (env RAFTKV_DATA_DIR)")
	bootstrap := fs.Bool("bootstrap", getenv("RAFTKV_BOOTSTRAP") == "true", "first start of a new cluster: allow an empty WAL (env RAFTKV_BOOTSTRAP)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if fs.NArg() > 0 {
		return Config{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	ps, err := parsePeers(*peers)
	if err != nil {
		return Config{}, fmt.Errorf("--peers: %w", err)
	}

	cfg := Config{ID: *id, Peers: ps, DataDir: *dataDir, Bootstrap: *bootstrap}

	return cfg, cfg.validate()
}

func (c Config) validate() error {
	if c.ID == 0 {
		return errors.New("--id is required and must be > 0")
	}

	if len(c.Peers) == 0 {
		return errors.New("--peers is required")
	}
	if !slices.ContainsFunc(c.Peers, func(p Peer) bool { return p.ID == c.ID }) {
		return fmt.Errorf("--id %d is not in --peers", c.ID)
	}

	if c.DataDir == "" {
		return errors.New("--data-dir must not be empty")
	}

	return nil
}

// parsePeers parses "1=host:port,2=host:port" into peers sorted by ID.
func parsePeers(s string) ([]Peer, error) {
	if s == "" {
		return nil, nil
	}

	var peers []Peer
	for item := range strings.SplitSeq(s, ",") {
		idStr, addr, ok := strings.Cut(strings.TrimSpace(item), "=")
		if !ok || addr == "" {
			return nil, fmt.Errorf("%q: want id=host:port", item)
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			return nil, fmt.Errorf("%q: id must be a positive integer", item)
		}

		if slices.ContainsFunc(peers, func(p Peer) bool { return p.ID == id }) {
			return nil, fmt.Errorf("duplicate id %d", id)
		}

		peers = append(peers, Peer{ID: id, Addr: addr})
	}

	slices.SortFunc(peers, func(a, b Peer) int { return cmp.Compare(a.ID, b.ID) })

	return peers, nil
}

func envOr(getenv func(string) string, key, def string) string {
	if v := getenv(key); v != "" {
		return v
	}

	return def
}

func envUint(getenv func(string) string, key string) uint64 {
	n, _ := strconv.ParseUint(getenv(key), 10, 64)
	return n
}
