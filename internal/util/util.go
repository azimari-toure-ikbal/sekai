package util

import (
	"context"
	"fmt"
	"os"
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

    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    tree, err := parser.ParseCtx(context.Background(), nil, source)

    if err != nil {
        return nil, fmt.Errorf("util:ParseFile: something went wrong when parsing the source")
    }

    rootNode := tree.RootNode()

    jsxElements := []string{}

    // Traverse the AST and collect the elements
    collectJSXElements(rootNode, source, &jsxElements)

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

    // Get the type of the node
    nodeType := n.Type()

    // Get the text of the node if it's a leaf
    nodeText := ""
    if n.ChildCount() == 0 {
        startByte := n.StartByte()
        endByte := n.EndByte()
        nodeText = string(sourceCode[startByte:endByte])
        nodeText = strings.TrimSpace(nodeText)
    }

    // Check if the node is a `string_fragment` and its parent property is `className`
    shouldSkipLineNumber := false
    if nodeType == "string_fragment" {
        parent := n.Parent().PrevNamedSibling()

        if parent != nil && parent.Type() == "property_identifier" && (string(sourceCode[parent.StartByte():parent.EndByte()]) == "className" || string(sourceCode[parent.StartByte():parent.EndByte()]) == "variant" || string(sourceCode[parent.StartByte():parent.EndByte()]) == "size") {
            shouldSkipLineNumber = true
        }
    }

    // Build the string representation with or without line numbers
    if nodeType == "string_fragment" || nodeType == "jsx_text" {
        if !shouldSkipLineNumber {
            startPoint := n.StartPoint()
            endPoint := n.EndPoint()
            builder.WriteString(fmt.Sprintf("%s- %s (lines %d-%d): %s\n", indent, nodeType, startPoint.Row+1, endPoint.Row+1, nodeText))
        } 
    } 

    // Recursively process the children nodes
    for i := 0; i < int(n.NamedChildCount()); i++ {
        child := n.NamedChild(i)
        childType := child.Type()

        if skipNested && !(childType == "jsx_element" || childType == "jsx_self_closing_element") {
            // Indicate that a nested JSX element exists
            childStr := nodeToString(child, sourceCode, level+1, skipNested)
            builder.WriteString(childStr)
        }
    }

    return builder.String()
}


func printNode(n *sitter.Node, sourceCode []byte, level int) {
    indent := strings.Repeat("  ", level)

    // Get the type of the node
    nodeType := n.Type()

    // Get the text of the node if it's a leaf
    nodeText := ""
    if n.ChildCount() == 0 {
        startByte := n.StartByte()
        endByte := n.EndByte()
        nodeText = string(sourceCode[startByte:endByte])
        nodeText = strings.TrimSpace(nodeText)
    }

    // Print the node with indentation
    if nodeText != "" {
        fmt.Printf("%s- %s: %s\n", indent, nodeType, nodeText)
    } else {
        fmt.Printf("%s- %s\n", indent, nodeType)
    }

    // Recursively print the children nodes
    for i := 0; i < int(n.NamedChildCount()); i++ {
        child := n.NamedChild(i)
        printNode(child, sourceCode, level+1)
    }
}