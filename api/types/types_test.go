cat > api/types/types_test.go << 'EOF'
package types

import (
	"encoding/json"
	"testing"
)

func TestCodeFileJSON(t *testing.T) {
	file := CodeFile{
		Path:     "test.go",
		Language: Go,
		Content:  "package main\n\nfunc main() {}",
		Hash:     "abc123",
	}

	data, err := json.Marshal(file)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshalled CodeFile
	err = json.Unmarshal(data, &unmarshalled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshalled.Path != file.Path {
		t.Errorf("expected %s, got %s", file.Path, unmarshalled.Path)
	}
	if unmarshalled.Language != file.Language {
		t.Errorf("expected %s, got %s", file.Language, unmarshalled.Language)
	}
}

func TestTestSuggestion(t *testing.T) {
	suggestion := TestSuggestion{
		Language:     Python,
		FunctionName: "add",
		TestCode:     "def test_add(): assert add(1, 2) == 3",
		TestFilePath: "test_add.py",
		Framework:    TestFrameworkPytest,
		Confidence:   0.95,
	}

	if suggestion.Framework != TestFrameworkPytest {
		t.Errorf("expected pytest, got %s", suggestion.Framework)
	}
}

func TestReadmeSection(t *testing.T) {
	section := ReadmeSection{
		Type:    DocSectionInstallation,
		Title:   "Installation",
		Content: "npm install my-package",
		Order:   1,
	}

	if section.Type != DocSectionInstallation {
		t.Errorf("expected installation section, got %s", section.Type)
	}
}
EOF