package nextjs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/azimari-toure-ikbal/translate-core/internal/util" // Import the util package
)

func RunForNext(files *[]string, inputLang, outputLang *string) error {
	if !util.CheckIfNextJS() {
		return fmt.Errorf("RunForNext:CheckIfNextJS: You must be at the root of a valid NextJS project")
	}

	dirToSkip := []string{
		"node_modules",
		".next",
		"public",
	}

	fileExtensions := []string {
		".tsx",
		".jsx",
	}

	if util.ReadConfig() != "" {
		dirToSkip = append(dirToSkip, strings.Split(util.ReadConfig(), "\n")...)
	}
	
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {

		// Skip directories
		for _, dir := range(dirToSkip) {
			if f.IsDir() && strings.Contains(path, dir) {
				return filepath.SkipDir
			}
		}

		// We check for files here and their extensions
		if !f.IsDir() {
			for _, extension := range(fileExtensions) {
				if strings.HasSuffix(f.Name(), extension) {
					*files = append(*files, path)
				}
			}
		}
		
		return nil
	 })

	 if err != nil {
		return fmt.Errorf("RunForNext:filepath.walk: Something went wrong when traversing the directories %v", err)
	 }

	fmt.Printf("We found a total of : %v files\n", len(*files))

	originalMap := make(map[string]string)
	re := regexp.MustCompile(`lines (\d+)-(\d+)`)
	reSpace := regexp.MustCompile(`\s+`)

	for key, el := range(*files) {
		
		r := strings.NewReplacer("[", "", "]", "", "(", "", ")", "", ".page.tsx", "", ".layout.tsx", "", ".page.jsx", "", ".layout.jsx", "", ".tsx", "", ".jsx", "")
		tradKey := r.Replace(strings.Join(strings.Split(el, "/"), "."))

		if util.IsDebugMode {
			fmt.Printf("The key to be used is %s\n\n", tradKey)
		}

		parsed, err := util.ParseFile(el)

		if err != nil {
			return fmt.Errorf("Later")
		}

		for _, val := range(parsed) {
			matches := re.FindStringSubmatch(val)

			if len(matches) == 3 {
				startLine := matches[1]
				if util.IsDebugMode {
					originalMap[fmt.Sprintf("%s.%s", tradKey, startLine)] = fmt.Sprintf(" key is %d : val is %s", key, strings.Split(val, ": ")[1])
				}
				originalMap[fmt.Sprintf(`"%s.%s"`, tradKey, startLine)] = fmt.Sprintf(`"%s"`, util.ReplaceApos(reSpace.ReplaceAllString(strings.TrimSpace(strings.Split(val, ": ")[1]), " ")))
			}
		}
	}

	err = util.WriteMapToJSONFile(originalMap, *inputLang)

	if err != nil {
		return fmt.Errorf("Something went wrong while writing the input file: %v", err)
	}

	url := "http://localhost:11434/api/generate"
	model := "trad"
	start := time.Now()
	results := util.TranslateConcurrently(url, model, *inputLang, *outputLang, originalMap)
	fmt.Printf("All translations completed in %v\n", time.Since(start))

	// for key, val := range results {
	// 	results[key]
	// }

	err = util.WriteMapToJSONFile(results, *outputLang)

	if err != nil {
		return fmt.Errorf("Something went wrong while writing the output file: %v", err)
	}
		
	return nil
}
