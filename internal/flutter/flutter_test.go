package flutter

import (
	"testing"
)

func TestRunForFlutter(t *testing.T) {
	files := []string{}
	inputLang := "en"

	err := RunForFlutter(&files, &inputLang)
	if err != nil {
		t.Errorf("RunForFlutter failed: %v", err)
	}

	if len(files) == 0 {
		t.Errorf("No files were parsed for the Flutter project")
	}
}

