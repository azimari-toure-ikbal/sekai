package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var devEnvFlag string
	var files []string

	flag.StringVar(&devEnvFlag, "env", "nextjs", "This should either be nextjs or flutter")
	flag.Parse()

	if devEnvFlag != "nextjs" && devEnvFlag != "flutter" {
		log.Fatal("You can't do that brother.")
	}

	if devEnvFlag == "flutter" {
		fmt.Println("Work in progress")
		return;
	}

	if devEnvFlag == "nextjs" {

		isNext := checkIfNextJS()

		if !isNext {
			log.Fatal("You must be at the root of a valid NextJS project.")
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

		if readConfig() != "" {
			dirToSkip = append(dirToSkip, strings.Split(readConfig(), "\n")...)
		}

		fmt.Println(dirToSkip)

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
						files = append(files, path)
					}
				}
			}
			
			return nil
		 })
	
		 if err != nil {
			fmt.Println("Error:", err)
			return
		 }
	}

	for _, file := range(files) {
		fmt.Printf("File found : %v\n", file)
	}

	fmt.Printf("We found a total of : %v files\n", len(files))
}

func checkIfNextJS() bool {
	f, error := os.Open("next.config.mjs")
	
	if error != nil {
		return false
	}
	
	defer f.Close()

	return true
}

func readConfig() string {
	fmt.Println("Searching for transcore config file...")
	data, error := os.ReadFile(".transcore")

	if error != nil {
		fmt.Println("Couldn't find config file. Proceed with default options")
		return ""
	}
	
	fmt.Println("Found transcore config file!")
	return string(data)
}