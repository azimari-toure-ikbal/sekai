package util

import (
	"context"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func IsDebugMode() bool {
    debug := os.Getenv("DEBUG")

    return debug == "true"
}

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
    pattern := `^\s*(\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}|\{\{\d+\}\})+$`
    regex := regexp.MustCompile(pattern)

    // Check if the entire string consists only of interpolated variables
    return regex.MatchString(s)
}
