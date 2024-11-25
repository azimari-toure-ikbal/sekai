package util

import (
	"bufio"
	"context"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
    "encoding/json"
    "strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

const IsDebugMode = false
var verbose bool = false

// EnableVerbose sets the verbose mode
func EnableVerbose() {
	verbose = true
}


// logVerbose logs messages only if verbose mode is enabled
func LogVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func CheckIfNextJS() bool {
	f, err := os.Open("next.config.mjs")
	
	if err != nil {
		return false
	}
	
	defer f.Close()

	return true
}

func CheckIfFlutter() bool {
    //look for a lib folder or a pubspec.yaml file
    _, err := os.Stat("lib")
    if err != nil {
        return false
    }

    _, err = os.Stat("pubspec.yaml")
    if err != nil {
        return false
    }
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

func parseJSFile(source []byte) ([]string, error) {
    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    tree, err := parser.ParseCtx(context.Background(), nil, source)
    if err != nil {
        return nil, fmt.Errorf("util:parseJSFile: error parsing JavaScript source")
    }

    rootNode := tree.RootNode()
    jsxElements := []string{}

    // Traverse the AST and collect JSX elements
    collectJSXElements(rootNode, source, &jsxElements)

    return jsxElements, nil
}

func parseDartFile(source []byte) ([]string, error) {
    LogVerbose("\n\nParsing Dart file...")
    
    // Comprehensive regex patterns for various text-displaying widgets and components
    textPatterns := []string{
        // Direct Text widgets
        `Text\s*\(\s*([^)]+)\)`,
        
        // TextSpan patterns
        `TextSpan\s*\(\s*text:\s*([^,)]+)`,
        
        // RichText patterns
        `RichText\s*\(\s*.*?text:\s*([^,)]+)`,
        
        // TextField patterns
        `TextField\s*\(\s*(?:decoration:\s*InputDecoration\s*\(\s*)?(?:label|hint|error)?Text:\s*([^,)]+)`,
        
        // AppBar title
        `AppBar\s*\(\s*title:\s*Text\s*\(\s*([^)]+)\)`,
        
        // SnackBar content
        `SnackBar\s*\(\s*content:\s*Text\s*\(\s*([^)]+)\)`,
        
        // Dialog text
        `(?:AlertDialog|SimpleDialog)\s*\(\s*(?:title|content):\s*Text\s*\(\s*([^)]+)\)`,
        
        // Tooltip patterns
        `Tooltip\s*\(\s*message:\s*([^,)]+)`,
        
        // Tab patterns
        `Tab\s*\(\s*text:\s*([^,)]+)`,
    }

    parsedStrings := []string{}
    sourceStr := string(source)

    for _, pattern := range textPatterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllStringSubmatch(sourceStr, -1)
        
        for _, match := range matches {
            if len(match) > 1 {
                extractedText := strings.TrimSpace(match[1])
                
                // More robust string extraction
                stringRe := regexp.MustCompile(`'(?:\\.|[^'\\])*'|"(?:\\.|[^"\\])*"`)
                stringMatches := stringRe.FindAllString(extractedText, -1)
                
                for _, stringMatch := range stringMatches {
                    // Remove surrounding quotes
                    cleaned := stringMatch[1 : len(stringMatch)-1]
                    
                    // Unescape any escaped characters
                    unescaped, err := strconv.Unquote(stringMatch)
                    if err == nil {
                        cleaned = unescaped
                    }
                    
                    // Avoid empty strings and duplicates
                    if cleaned != "" && !contains(parsedStrings, cleaned) {
                        LogVerbose("Captured string: %s", cleaned)
                        parsedStrings = append(parsedStrings, cleaned)
                    }
                }            }
        }
    }

    LogVerbose("Parsed strings: %v", parsedStrings)
    return parsedStrings, nil
}

// Helper function to check for duplicates
func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func ParseFile(file string) ([]string, error) {
	LogVerbose("\n\nParsing file: %s", file)
	source, err := os.ReadFile(file)
	if err != nil {
		LogVerbose("Error opening file %s: %v", file, err)
		return nil, fmt.Errorf("util:ParseFile: couldn't open the file %v", file)
	}

	unescapedSource := []byte(html.UnescapeString(string(source)))

	var parsedElements []string
	if strings.HasSuffix(file, ".dart") {
		LogVerbose("Detected Dart file. Parsing...")
		parsedElements, err = parseDartFile(unescapedSource)
		if err != nil {
			LogVerbose("Error parsing Dart file %s: %v", file, err)
			return nil, fmt.Errorf("util:ParseFile: error parsing Dart file: %v", err)
		}
	} else if strings.HasSuffix(file, ".js") || strings.HasSuffix(file, ".jsx") || strings.HasSuffix(file, ".tsx") {
		LogVerbose("Detected JavaScript/TypeScript file. Parsing...")
		parsedElements, err = parseJSFile(unescapedSource)
		if err != nil {
			LogVerbose("Error parsing JavaScript/TypeScript file %s: %v", file, err)
			return nil, fmt.Errorf("util:ParseFile: error parsing JavaScript/TypeScript file: %v", err)
		}
	} else {
		LogVerbose("Unsupported file type: %s", file)
		return nil, fmt.Errorf("util:ParseFile: unsupported file type for %v", file)
	}

	LogVerbose("Successfully parsed file: %s", file)
	return parsedElements, nil
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

func WriteMapToJSONFile(originalMap map[string]string, inputLang string) error {
    file, err := os.Create(fmt.Sprintf("%s.json", inputLang))
    if err != nil {
        return fmt.Errorf("Error creating file: %v", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    writer.WriteString("{\n")

    count := 0
    total := len(originalMap)

    // Iterate over the map and write each key-value pair to the file
    for key, value := range originalMap {
        count++
        if count == total {
            // Last item, don't add a comma
            _, err := writer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
            if err != nil {
                return fmt.Errorf("Error writing to file: %v", err)
            }
        } else {
            // Not the last item, add a comma
            _, err := writer.WriteString(fmt.Sprintf("%s: %s,\n", key, value))
            if err != nil {
                return fmt.Errorf("Error writing to file: %v", err)
            }
        }
    }

    writer.WriteString("}\n")

    // Flush the buffered writer to the file
    err = writer.Flush()
    if err != nil {
        return fmt.Errorf("Error flushing writer: %v", err)
    }

    return nil
}

func SanitizeKey(input string) string {
	// Sanitize text for ARB keys (remove special characters, spaces, etc.)
	re := regexp.MustCompile(`[^\w]+`)
	return re.ReplaceAllString(strings.ToLower(input), "_")
}

func WriteMapToJSONFileFlutter(originalMap map[string]string, inputLang string) error {
    LogVerbose("\n\nStarting WriteMapToJSONFileFlutter...")
    outputFile := fmt.Sprintf("%s", inputLang) // Use inputLang to define the output file name
    LogVerbose("Output file: %s", outputFile)

    // Step 1: Prepare data
    LogVerbose("Preparing data for the JSON file...")
    data := map[string]interface{}{
        "@@locale": inputLang, // Use inputLang as the locale
    }

    LogVerbose("\nOriginal map: %v", originalMap)

    for key, value := range originalMap {
        LogVerbose("Adding key: '%s', value: '%s'", key, value)
        data[key] = value
        data[fmt.Sprintf("@%s", key)] = map[string]string{
            "description": value,
        }
    }
    LogVerbose("Data preparation complete. Total keys: %d", len(originalMap))

    // Step 2: Create the file
    LogVerbose("Creating file: %s", outputFile)
    file, err := os.Create(outputFile)
    if err != nil {
        LogVerbose("Error creating file '%s': %v", outputFile, err)
        return fmt.Errorf("Error creating file: %v", err)
    }
    defer func() {
        LogVerbose("Closing file: %s", outputFile)
        file.Close()
    }()

    // Step 3: Encode the data to JSON
    LogVerbose("Encoding data to JSON...")
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(data); err != nil {
        LogVerbose("Error encoding data to JSON: %v", err)
        return fmt.Errorf("Error writing to JSON file: %v", err)
    }

    LogVerbose("Successfully wrote JSON data to file: %s", outputFile)
    return nil
}

