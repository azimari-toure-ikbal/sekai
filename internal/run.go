package internal

import (
	"flag"
	"fmt"

	"github.com/azimari-toure-ikbal/translate-core/internal/nextjs"
)

func Run() error {
	var devEnvFlag, inputLang, outputLang string
	var files []string

	acceptedLang := []string{"fr", "en", "de", "it", "jp", "po", "es"}
	// Helper function to check if the selected language is valid
	isValidLang := func(lang string, acceptedLangs []string) bool {
		for _, l := range acceptedLangs {
			if lang == l {
				return true
			}
		}
		return false
	}

	flag.StringVar(&devEnvFlag, "env", "nextjs", "We only support NextJS for the moment")
	flag.StringVar(&inputLang, "i", "en", "Please select the input language")
	flag.StringVar(&outputLang, "o", "", "Please select the ouput language")
	flag.Parse()

	if devEnvFlag != "nextjs" && devEnvFlag != "flutter" {
		return fmt.Errorf("You can't do that brother")
	}

	if inputLang == "" {
		return fmt.Errorf("The input language cannot be empty")
	}

	if outputLang == "" {
		return fmt.Errorf("The output language cannot be empty")
	}

	if !isValidLang(outputLang, acceptedLang) {
		return fmt.Errorf("You must chose a valid output language in this list %v", acceptedLang)
	}

	if inputLang == outputLang {
		return fmt.Errorf("The input and output languages cannot be the same")
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