package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/example/my-mcp-server/internal" // TODO: Replace with your module path
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func getVersionInfo() (string, string, string) {
	if version != "dev" {
		return version, commit, date
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v == "" || v == "(devel)" {
			v = "dev"
		}

		var rev, buildTime string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				rev = setting.Value
			case "vcs.time":
				buildTime = setting.Value
			}
		}

		if rev == "" {
			rev = "unknown"
		}
		if buildTime == "" {
			buildTime = "unknown"
		}

		return v, rev, buildTime
	}

	return "dev", "unknown", "unknown"
}

func printVersion() {
	version, commit, date := getVersionInfo()
	fmt.Printf("my-mcp-server version %s\n", version) // TODO: Replace with your server name
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built: %s\n", date)
}

func main() {
	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "version" {
		printVersion()
		os.Exit(0)
	}

	version, _, _ := getVersionInfo()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := internal.GetServer(version)
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Printf("Server error: %v", err)
	}
}
