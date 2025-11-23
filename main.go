package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DavidHoenisch/oxidation9/internal/types"
	"github.com/DavidHoenisch/oxidation9/pkg/dns"
	"github.com/DavidHoenisch/oxidation9/pkg/scan"
	"github.com/DavidHoenisch/oxidation9/pkg/spam"
)

func getHomePath() string {
	path, homeErr := os.UserHomeDir()
	if homeErr != nil {
		log.Print(homeErr)
	}
	return path
}

func getExecutablePath() string {
	path, execErr := os.Executable()
	if execErr != nil {
		log.Print(execErr)
	}
	return path
}

var funcMap = map[string]types.Tool{
	"spam": &spam.Spam{},
	"dns":  &dns.Dns{},
	"scan": &scan.Scan{},
}

func main() {
	if len(os.Args) < 1 {
		return
	}

	binaryName := filepath.Base(os.Args[0])

	switch binaryName {

	case "oxidation9", "ox9":
		if len(os.Args) < 2 {
			fmt.Println("Usage: oxidation9 <command> [flags]")
			os.Exit(1)
		}

		subCommand := os.Args[1]

		switch subCommand {
		case "bootstrap":
			runBootstrap()
		case "clean":
			runClean()
		case "spam":
			runTool("spam", os.Args[2:])
		case "dns":
			runTool("dns", os.Args[2:])
		case "scan":
			runTool("scan", os.Args[2:])
		default:
			fmt.Printf("Unknown command: %s\n", subCommand)
			os.Exit(1)
		}

	case "spam":
		runTool("spam", os.Args[1:])
	case "dns":
		runTool("dns", os.Args[1:])
	case "scan":
		runTool("scan", os.Args[1:])
	default:
		if _, ok := funcMap[binaryName]; ok {
			runTool(binaryName, os.Args[1:])
		} else {
			fmt.Println("Unknown binary name or command")
			os.Exit(1)
		}
	}
}

func runTool(name string, args []string) {
	tool := funcMap[name]

	err := tool.Handler(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runBootstrap() {
	ePath := getExecutablePath()
	hPath := getHomePath()
	for key := range funcMap {
		binPath := fmt.Sprintf("%s/.local/bin/%s", hPath, key)
		err := os.Symlink(ePath, binPath)
		if os.IsExist(err) {
			fmt.Printf("%s already boostrapped...skipping\n", key)
			continue
		}
		if err != nil {
			fmt.Printf("Error bootstrapping %s: %v", key, err)
		}
		fmt.Printf("Bootstrapped %s\n", key)
	}
}

func runClean() {
	hPath := getHomePath()
	for key := range funcMap {
		binPath := fmt.Sprintf("%s/.local/bin/%s", hPath, key)
		err := os.Remove(binPath)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("Cleaned %s\n", key)
	}
}
