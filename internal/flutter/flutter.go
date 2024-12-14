package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/azimari-toure-ikbal/sekai-core/internal/util"
)


func RunForFlutter(files *[]string, inputLang *string, outputLang *string) error {
	if !util.CheckIfFlutter() {
		return fmt.Errorf("RunForFlutter:CheckIfFlutter: You must be at the root of a valid Flutter project")
	}

	dirToSkip := []string{
		".dart_tool",
		".gradle",
		".idea",
		".pub",
		"fonts",
		"build",
		"android",
		"ios",
		"web",
		"macos",
		"linux",
		"windows",
		"test",
	}

	util.LogVerbose("Skipping directories: %v", dirToSkip)

	// Walk through Flutter project to gather `.dart` files
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			util.LogVerbose("Error accessing path %s: %v", path, err)
			return err
		}
		// Skip specified directories
		for _, dir := range dirToSkip {
			if f.IsDir() && strings.Contains(path, dir) {
				util.LogVerbose("Skipping directory: %s", path)
				return filepath.SkipDir
			}
		}
		// Include `.dart` files
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".dart") {
			util.LogVerbose("Found Dart file: %s", path)
			*files = append(*files, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("RunForFlutter:filepath.walk: Something went wrong when traversing directories: %v", err)
	}

	util.LogVerbose("Found %d Dart files for translation.", len(*files))

	originalMap := make(map[string]string)
	for _, file := range *files {
		util.LogVerbose("Parsing file: %s", file)
		parsed, err := util.ParseFile(file)
		if err != nil {
			return fmt.Errorf("RunForFlutter:ParseFile: error parsing file %v: %v", file, err)
		}
		
		for _, element := range parsed {
			// Generate keys from parsed strings
			key := generateFlutterKey(file, element)
			util.LogVerbose("Generated key: %s for text: %s", key, element)
			originalMap[key] = element
		}
	}



	// Write translations to `l10n/app_{inputLang}.arb`
	inputFile := fmt.Sprintf("lib/l10n/app_%s.arb", *inputLang)
	util.LogVerbose("Writing translations to file: %s", inputFile)
	err = util.WriteMapToJSONFileFlutter(originalMap, inputFile)
	if err != nil {
		return fmt.Errorf("RunForFlutter: error writing translations to file: %v", err)
	}

	//Make translations
	url := "http://localhost:11434/api/generate"
	model := "trad"
	start := time.Now()
	results := util.TranslateConcurrently(url, model, *inputLang, *outputLang, originalMap)
	fmt.Printf("All translations completed in %v\n", time.Since(start))

	// Write translations to `l10n/app_{outputLang}.arb`
	outputFile := fmt.Sprintf("lib/l10n/app_%s.arb", *outputLang)
	util.LogVerbose("Writing translations to file: %s", outputFile)
	err = util.WriteMapToJSONFileFlutter(results, outputFile)

	if err != nil {
		return fmt.Errorf("RunForFlutter: error writing translations to file: %v", err)
	}

	util.LogVerbose("Translations successfully written to %s", outputFile)
	return nil
}

func generateFlutterKey(filePath, text string) string {
	// Generate unique keys based on file path and text content
	baseKey := strings.ReplaceAll(filepath.Base(filePath), ".dart", "")
	key := fmt.Sprintf("%s_%s", baseKey, util.SanitizeKey(text))
	util.LogVerbose("Generated key %s from filePath %s and text %s", key, filePath, text)
	return key
}
