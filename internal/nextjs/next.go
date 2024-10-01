package nextjs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/azimari-toure-ikbal/translate-core/internal/util" // Import the util package
)

func RunForNext(files *[]string) error {
	if !util.CheckIfNextJS() {
		return fmt.Errorf("next:CheckIfNextJS: You must be at the root of a valid NextJS project")
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
		return fmt.Errorf("next:filepath.walk: Something went wrong when traversing the directories %v", err)
	 }

	fmt.Printf("We found a total of : %v files\n", len(*files))

	file, err := os.ReadFile((*files)[10])

	if err != nil {
		return fmt.Errorf("util:ParseJSX: Couldn't open the file %s", (*files)[10])
	}

	return nil
}
