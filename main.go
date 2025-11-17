package main

import (
	"fmt"
	"github.com/DavidHoenisch/oxidation9/internal/models"
	"github.com/DavidHoenisch/oxidation9/pkg/spam"
	"log"
	"os"
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

var funcMap = map[string]func(models.FuncParams) error{
	"spam": func(params models.FuncParams) error {
		return spam.Run(params)
	},
}

func main() {
	if len(os.Args) < 1 {
		log.Println("not enough args")
	}

	caller := os.Args[0]

	switch caller {
	case "oxidation9", "ox9", "./oxidation9":
		args := os.Args[1]
		switch args {
		case "bootstrap":
			ePath := getExecutablePath()
			hPath := getHomePath()
			for key, _ := range funcMap {
				binPath := fmt.Sprintf("%s/.local/bin/%s", hPath, key)

				err := os.Symlink(ePath, binPath)
				if err != nil {
					log.Print(err)
				}
			}
		case "clean":
			hPath := getHomePath()
			for key, _ := range funcMap {
				binPath := fmt.Sprintf("%s/.local/bin/%s", hPath, key)

				err := os.Remove(binPath)
				if err != nil {
					log.Print(err)
				}
			}
		default:
			fmt.Println("Unknown command")
		}
	case "spam":
		err := funcMap["spam"](models.FuncParams{
			Args: map[string]any{
				"count": os.Args[1],
				"url":   os.Args[2],
			},
		})
		if err != nil {
			fmt.Println("something")
		}
	default:
		fmt.Println("Unknown command")
	}
}
