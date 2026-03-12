package frontmatter

import (
	"bytes"
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata Front-matter 元数据
type Metadata struct {
	Title      string   `yaml:"title"`
	Slug       string   `yaml:"slug"`
	Date       string   `yaml:"date"`
	Category   string   `yaml:"category"`
	Tags       []string `yaml:"tags"`
	Summary    string   `yaml:"summary"`
	Cover      string   `yaml:"cover"`
	Status     string   `yaml:"status"` // draft/published
}

// Parse 解析 Markdown Front-matter
func Parse(content []byte) (*Metadata, string, error) {
	// 检查是否以 --- 开头
	if !bytes.HasPrefix(content, []byte("---\n")) && !bytes.HasPrefix(content, []byte("---\r\n")) {
		return nil, string(content), nil // 没有 front-matter
	}

	// 查找第二个 ---
	lines := bytes.Split(content, []byte("\n"))
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		line := bytes.TrimSpace(lines[i])
		if bytes.Equal(line, []byte("---")) {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return nil, "", errors.New("invalid front-matter format")
	}

	// 提取 YAML 部分
	yamlContent := bytes.Join(lines[1:endIndex], []byte("\n"))

	// 提取正文部分
	bodyContent := bytes.Join(lines[endIndex+1:], []byte("\n"))

	// 解析 YAML
	var meta Metadata
	if err := yaml.Unmarshal(yamlContent, &meta); err != nil {
		return nil, "", err
	}

	return &meta, strings.TrimSpace(string(bodyContent)), nil
}
