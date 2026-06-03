package lin_test

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

const maxSkillDescriptionLength = 1024

func TestSkillMetadataDescriptionLength(t *testing.T) {
	content, err := os.ReadFile("skills/lin/SKILL.md")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	parts := strings.SplitN(string(content), "---\n", 3)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "" {
		t.Fatal("SKILL.md must start with YAML frontmatter")
	}

	var metadata struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		t.Fatalf("frontmatter YAML: %v", err)
	}

	if metadata.Name != "lin" {
		t.Fatalf("name = %q, want lin", metadata.Name)
	}
	if metadata.Description == "" {
		t.Fatal("description must not be empty")
	}
	if got := len(metadata.Description); got > maxSkillDescriptionLength {
		t.Fatalf("description length = %d, max %d", got, maxSkillDescriptionLength)
	}
}
