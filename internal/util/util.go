package util

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

const IsDebugMode = false

func CheckIfNextJS() bool {
	f1, err1 := os.Open("next.config.mjs")
    f2, err2 := os.Open("next.config.js")
	
	if err1 != nil && err2 != nil {
		return false
	}
	
	defer f1.Close()
    defer f2.Close()

	return true
}

func ReadConfig() string {
	fmt.Println("Searching for transcore config file...")
	data, err := os.ReadFile(".sekai.config")

	if err != nil {
		fmt.Println("Couldn't find config file. Proceed with default options")
		return ""
	}
	
	fmt.Println("Found sekai-core config file!")
	return string(data)
}

func ParseFile(file string) ([]string, error){
    source, err := os.ReadFile(file)

    if err != nil {
        return nil, fmt.Errorf("util:ParseFile: couldn't open the file %v", file)
    }

    // Revert any encoded character to its original value
    unescapedSource := []byte(html.UnescapeString(string(source)))

    // fmt.Println(string(unescapedSource))


    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    tree, err := parser.ParseCtx(context.Background(), nil, unescapedSource)

    if err != nil {
        return nil, fmt.Errorf("util:ParseFile: something went wrong when parsing the source")
    }

    rootNode := tree.RootNode()

    jsxElements := []string{}

    // Traverse the AST and collect the elements
    collectJSXElements(rootNode, unescapedSource, &jsxElements)

    // fmt.Println(rootNode)

    // if DEBUG {
    //     fmt.Printf("Length of jsx elements is %v\n\n", len(jsxElements))
    // }

    return jsxElements, nil
}

// collectJSXElements recursively traverses the AST and collects JSX elements' AST representations
func collectJSXElements(n *sitter.Node, sourceCode []byte, jsxElements *[]string) {
    nodeType := n.Type()

    if nodeType == "jsx_element" || nodeType == "jsx_self_closing_element" {
        // Generate the AST representation of this node, skipping nested JSX elements
        astRepresentation := nodeToString(n, sourceCode, 0, true)
        *jsxElements = append(*jsxElements, astRepresentation)
    }

    // Recursively traverse the children
    for i := 0; i < int(n.ChildCount()); i++ {
        child := n.Child(i)
        collectJSXElements(child, sourceCode, jsxElements)
    }
}

// nodeToString recursively builds a string representation of the node, its children, and includes line numbers
func nodeToString(n *sitter.Node, sourceCode []byte, level int, skipNested bool) string {
    var builder strings.Builder
    indent := strings.Repeat("  ", level)

    nodeType := n.Type()

    // Process JSX elements to accumulate text and interpolations
    if nodeType == "jsx_element" || nodeType == "jsx_self_closing_element" {
        var interpolatedText strings.Builder
        varCounter := 1

        for i := 0; i < int(n.ChildCount()); i++ {
            child := n.Child(i)
            childType := child.Type()

            // fmt.Printf("child, child type, child value, children count %v || %s || %s || %d \n\n", child, childType, string(sourceCode[child.StartByte():child.EndByte()]), child.ChildCount())

            if childType == "identifier" || (childType == "jsx_expression" && string(sourceCode[child.StartByte():child.EndByte()]) != `{" "}`){
                // Add an interpolation placeholder for identifiers
                interpolatedText.WriteString(fmt.Sprintf("{{var%d}}", varCounter))
                varCounter++
            } else if childType == "jsx_text" || childType == "string_fragment" {
                // Append plain text for jsx_text and string_fragment
                startByte := child.StartByte()
                endByte := child.EndByte()
                childText := string(sourceCode[startByte:endByte])

                interpolatedText.WriteString(childText)                
            }
        }

        // Get line numbers and add accumulated text with interpolations to the builder
        startPoint := n.StartPoint()
        endPoint := n.EndPoint()
        if emtpyString := interpolatedText.String() == ""; !emtpyString {
            if onlyInterpolated := containsOnlyInterpolatedVariables(interpolatedText.String()); !onlyInterpolated {
                builder.WriteString(fmt.Sprintf("%s- %s (lines %d-%d): %s\n", indent, nodeType, startPoint.Row+1, endPoint.Row+1, interpolatedText.String()))
            }
        }
    }

    // Recursively process non-JSX elements
    for i := 0; i < int(n.NamedChildCount()); i++ {
        child := n.NamedChild(i)
        if skipNested && (child.Type() == "jsx_element" || child.Type() == "jsx_self_closing_element") {
            continue
        }
        childStr := nodeToString(child, sourceCode, level+1, skipNested)
        builder.WriteString(childStr)
    }

    return builder.String()
}

func containsOnlyInterpolatedVariables(s string) bool {
    // Define a regular expression pattern for one or more interpolated variables
    pattern := `^\s*("[{]{2}[a-zA-Z_][a-zA-Z0-9_]*[}]{2}"|"{2}\d+[}]{2}"|[{]{2}[a-zA-Z_][a-zA-Z0-9_]*[}]{2}|[{]{2}\d+[}]{2})+$`
    regex := regexp.MustCompile(pattern)

    // Check if the entire string consists only of interpolated variables
    return regex.MatchString(s)
}

// WriteMapToJSONFile writes a map to a JSON file, sorted by key grouping and line number.
func WriteMapToJSONFile(originalMap map[string]string, inputLang string) error {
	// Ensure the locales directory exists
	err := os.MkdirAll("./locales", os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating locales directory: %v", err)
	}

	// Extract and sort keys
	keys := make([]string, 0, len(originalMap))
	for key := range originalMap {
		keys = append(keys, key)
	}

	// Sort keys logically by path and line number
	sort.Slice(keys, func(i, j int) bool {
		pathI, lineI := splitKey(keys[i])
		pathJ, lineJ := splitKey(keys[j])

		if pathI == pathJ {
			return lineI < lineJ // Sort by line number if paths are equal
		}
		return pathI < pathJ // Sort by path
	})

	// Create the JSON file
	file, err := os.Create(fmt.Sprintf("./locales/%s.json", inputLang))
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("{\n")

	// Write sorted key-value pairs to the file
	for i, key := range keys {
		value := originalMap[key]

		// Write key-value pair without escaping
		_, err := writer.WriteString(fmt.Sprintf("%s: %s", key, value))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}

		// Add a comma except for the last item
		if i < len(keys)-1 {
			writer.WriteString(",\n")
		} else {
			writer.WriteString("\n")
		}
	}

	writer.WriteString("}\n")

	// Flush the buffered writer
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}

// splitKey splits a key into its path and line number for sorting.
func splitKey(key string) (string, int) {
	// Matches the last segment of the key as a number
	re := regexp.MustCompile(`^(.*)\.(\d+)$`)
	matches := re.FindStringSubmatch(key)

	if len(matches) != 3 {
		return key, 0 // Return the original key and a default number if it doesn't match the format
	}

	// Extract path and line number
	line, _ := strconv.Atoi(matches[2]) // Convert line number to int
	return matches[1], line
}

type TranslationChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// translate sends a single translation request and reads the streamed response.
func translate(url, model, input, output, text string) (string, error) {
	payload := map[string]interface{}{
		"model":  model,
		"prompt": fmt.Sprintf("%s:%s: %s", input, output, text),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var translationResult string
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading response line: %w", err)
		}

		var chunk TranslationChunk
		if err := json.Unmarshal(line, &chunk); err != nil {
			return "", fmt.Errorf("error unmarshalling response chunk: %w", err)
		}

		translationResult += chunk.Response
		if chunk.Done {
			break
		}
	}

	return translationResult, nil
}

// TranslateConcurrently processes translations in parallel and tracks progress.
func TranslateConcurrently(url, model, input, output string, texts map[string]string) map[string]string {
    var wg sync.WaitGroup
    sem := make(chan struct{}, 10) // Limit to 10 concurrent requests
    progress := make(chan string) // Channel to track progress
    results := make(map[string]string)
    resultsMu := sync.Mutex{}
    malformedTranslations := make(map[string]string) // To log malformed responses

    // Goroutine to monitor progress
    go func() {
        total := len(texts)
        completed := 0
        for range progress {
            completed++
            fmt.Printf("\rProgress: %d/%d translations completed", completed, total)
        }
        fmt.Println() // New line after all translations are complete
    }()

    for key, text := range texts {
        wg.Add(1)
        sem <- struct{}{} // Acquire a spot

        go func(key, text string) {
            defer wg.Done()
            defer func() { <-sem }() // Release the spot

            translatedText, err := translate(url, model, input, output, text)
            if err != nil {
                fmt.Printf("Error translating key %q: %v\n", key, err)
                return
            }

            // Validate split result
            parts := strings.Split(translatedText, ":")
            if len(parts) < 2 {
                // fmt.Printf("\n /?\\ Malformed response for key %s: %s /?\\ \n\n", key, translatedText)

                // Store malformed translation for further inspection
                resultsMu.Lock()
                malformedTranslations[key] = translatedText
                resultsMu.Unlock()
                return
            }

            // Safely store the result
            cleanedResult := ReplaceApos(strings.TrimSpace(parts[1]))
            resultsMu.Lock()
            results[key] = fmt.Sprintf(`"%s"`, cleanedResult)
            resultsMu.Unlock()

            // Send progress update
            progress <- key
        }(key, text)
    }

    // Wait for all translations to complete
    wg.Wait()
    close(progress)

    // Log all malformed translations after processing
    if len(malformedTranslations) > 0 {
        fmt.Printf("\n%d Malformed translations encountered:\n", len(malformedTranslations))
        for key, value := range malformedTranslations {
            fmt.Printf("- Key: %s, Response: %s\n", key, value)
        }
    }

    return results
}


func ReplaceApos(text string) string {
    r := strings.NewReplacer(`"`, `'`)
    return r.Replace(text)
}