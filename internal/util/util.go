package util

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func CheckIfNextJS() bool {
	f, err := os.Open("next.config.mjs")
	
	if err != nil {
		return false
	}
	
	defer f.Close()

	return true
}

func ReadConfig() string {
	fmt.Println("Searching for transcore config file...")
	data, err := os.ReadFile(".transcore")

	if err != nil {
		fmt.Println("Couldn't find config file. Proceed with default options")
		return ""
	}
	
	fmt.Println("Found transcore config file!")
	return string(data)
}

func ParseJSX(path string) ([]string, error) {
	// Here we should open the given file and parse it's content to get the tags where we have some text
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("util:ParseJSX: Couldn't open the file %s", path)
	}

	defer file.Close()

	doc, err := html.Parse(file)

	if err != nil {
		return nil, fmt.Errorf("util:ParseJSX: Couldn't parse the file %s", path)
	}

	var texts []string
	extractText(doc, &texts)

		// Afficher les textes extraits
		for _, text := range texts {
			fmt.Println(text)
		}

	return texts, nil
}

// Fonction pour extraire tout texte contenu entre des balises
func extractText(n *html.Node, texts *[]string) {
	if n.Type == html.TextNode {
		trimmedText := strings.TrimSpace(n.Data)
		if len(trimmedText) > 0 {
			*texts = append(*texts, trimmedText)
		}
	}
	// Parcourir les enfants récursivement
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		extractText(child, texts)
	}
}
