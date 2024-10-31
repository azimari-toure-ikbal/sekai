package nextjs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/azimari-toure-ikbal/translate-core/internal/util" // Import the util package
)

func RunForNext(files *[]string) error {
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


	// if util.IsDebugMode() {
	// 	fmt.Printf("The file path is %s\n\n", strings.Join(strings.Split((*files)[105], "/"), "."))
	// }

	originalMap := make(map[string]string)
	re := regexp.MustCompile(`lines (\d+)-(\d+)`)

	// test, _ := util.ParseFile((*files)[105])

	// for _, val := range(test) {

	// 	if val != "" {
	// 		matches := re.FindStringSubmatch(val)
	
	// 		if len(matches) == 3 {
	// 			startLine := matches[1]
	// 			if strings.Split(val, ": ")[1] != "" {
	// 				originalMap[fmt.Sprintf("%s.%s", strings.Join(strings.Split((*files)[105], "/"), "."),startLine)] = strings.Split(val, ": ")[1]

	// 			}


	// 		} else {
	// 			return fmt.Errorf("RunForNext:CheckIfNextJS: Something went wrong when collecting texts")
	// 		}
	// 	}
	// }

	for key, el := range(*files) {
		if util.IsDebugMode() {
			fmt.Printf("The file path is %s\n\n", strings.Join(strings.Split(el, "/"), "."))
		}

		parsed, err := util.ParseFile(el)

		if err != nil {
			return fmt.Errorf("Later")
		}

		for _, val := range(parsed) {
			matches := re.FindStringSubmatch(val)

			if len(matches) == 3 {
				startLine := matches[1]
				originalMap[fmt.Sprintf("%s.%s", strings.Join(strings.Split(el, "/"), "."),startLine)] = fmt.Sprintf(" key is %d : val is %s", key, strings.Split(val, ": ")[1])
			}
		}
	}


	fmt.Println(originalMap)

	return nil
}
