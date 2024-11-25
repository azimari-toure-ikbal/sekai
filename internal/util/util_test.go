package util

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	// Test JavaScript/TypeScript parsing
	jsSource := []byte(`
		export default function App() {
			return <div>Hello World</div>;
		}
	`)
	jsResults, err := parseJSFile(jsSource)
	if err != nil || len(jsResults) == 0 {
		t.Errorf("Failed to parse JavaScript file: %v", err)
	}

	// Test Dart parsing
	dartSource := []byte(`
		Widget build(BuildContext context) {
			return Text("Hello World");
		}
	`)
	dartResults, err := parseDartFile(dartSource)
	if err != nil || len(dartResults) == 0 {
		t.Errorf("Failed to parse Dart file: %v", err)
	}
}
