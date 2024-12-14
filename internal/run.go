package internal

import (
	"flag"
	"fmt"

	"github.com/azimari-toure-ikbal/sekai-core/internal/nextjs"
	"github.com/azimari-toure-ikbal/sekai-core/internal/flutter"
	"github.com/azimari-toure-ikbal/sekai-core/internal/util"
)

func Run() error {
	var devEnvFlag, inputLang, outputLang string
	var files []string
	var verbose bool // New verbose flag

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

	// Define flags
	flag.StringVar(&devEnvFlag, "env", "nextjs", "We support NextJS and Flutter projects")
	flag.StringVar(&inputLang, "i", "en", "Please select the input language")
	flag.StringVar(&outputLang, "o", "", "Please select the output language")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode") // Add verbose flag
	flag.Parse()


	// Enable verbose mode for relevant modules
	if verbose {
		fmt.Println("Verbose mode enabled")
		util.EnableVerbose()
	}

	if devEnvFlag != "nextjs" && devEnvFlag != "flutter" {
		return fmt.Errorf("Invalid environment. We support only NextJS and Flutter for now")
	}

	if inputLang == "" {
		return fmt.Errorf("The input language cannot be empty")
	}

	if outputLang == "" {
		return fmt.Errorf("The output language cannot be empty")
	}

	if !isValidLang(outputLang, acceptedLang) {
		return fmt.Errorf("You must choose a valid output language from this list: %v", acceptedLang)
	}

	if inputLang == outputLang {
		return fmt.Errorf("The input and output languages cannot be the same")
	}

	switch devEnvFlag {
	case "nextjs":
		err := nextjs.RunForNext(&files, &inputLang)
		if err != nil {
			return err
		}
	case "flutter":
		err := flutter.RunForFlutter(&files, &inputLang)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid environment. We support only NextJS and Flutter for now")
	}

	return nil
}
