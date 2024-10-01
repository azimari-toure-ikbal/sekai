package internal

import (
	"flag"
	"fmt"

	"github.com/azimari-toure-ikbal/translate-core/internal/nextjs"
)

func Run() error {
	var devEnvFlag string
	var files []string

	flag.StringVar(&devEnvFlag, "env", "nextjs", "We only support NextJS for the moment")
	flag.Parse()

	if devEnvFlag != "nextjs" && devEnvFlag != "flutter" {
		return fmt.Errorf("You can't do that brother")
	}

	switch devEnvFlag {
		case "nextjs":
			err := nextjs.RunForNext(&files)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("You provided a wrong value of env. We only support --env nextjs for the moment")
	}

	return nil
}