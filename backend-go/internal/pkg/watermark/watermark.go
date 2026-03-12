package watermark

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 零宽字符集：用于编码水印
const (
	zeroWidthSpace     = "\u200B" // 00
	zeroWidthNonJoiner = "\u200C" // 01
	zeroWidthJoiner    = "\u200D" // 10
	zeroWidthNoBreak   = "\uFEFF" // 11
)

// Inject 在 markdown 内容中注入零宽字符水印
func Inject(content, clientIP string) string {
	// 生成水印标识：IP hash前8位 + 时间戳
	watermarkID := generateWatermarkID(clientIP)

	// 将标识编码为零宽字符
	encoded := encode(watermarkID)

	// 在段落之间注入水印（避免代码块）
	return injectIntoContent(content, encoded)
}

// generateWatermarkID 生成水印标识
func generateWatermarkID(clientIP string) string {
	// IP hash
	hash := sha256.Sum256([]byte(clientIP))
	ipHash := hex.EncodeToString(hash[:])[:8]

	// Unix 时间戳（base36 编码缩短长度）
	timestamp := strconv.FormatInt(time.Now().Unix(), 36)

	return ipHash + timestamp
}

// encode 将字符串编码为零宽字符序列
func encode(s string) string {
	var result strings.Builder

	for _, ch := range s {
		// 将字符转为二进制
		binary := fmt.Sprintf("%08b", ch)

		// 每2位映射为一个零宽字符
		for i := 0; i < len(binary); i += 2 {
			bits := binary[i : i+2]
			switch bits {
			case "00":
				result.WriteString(zeroWidthSpace)
			case "01":
				result.WriteString(zeroWidthNonJoiner)
			case "10":
				result.WriteString(zeroWidthJoiner)
			case "11":
				result.WriteString(zeroWidthNoBreak)
			}
		}
	}

	return result.String()
}

// injectIntoContent 在内容中注入水印
func injectIntoContent(content, watermark string) string {
	lines := strings.Split(content, "\n")
	inCodeBlock := false
	var eligibleIndexes []int

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		isFence := strings.HasPrefix(trimmed, "```")

		if !inCodeBlock && isWatermarkEligibleLine(trimmed) {
			eligibleIndexes = append(eligibleIndexes, i)
		}

		if isFence {
			inCodeBlock = !inCodeBlock
		}
	}

	if len(eligibleIndexes) == 0 {
		return content
	}

	watermarkRunes := []rune(watermark)
	chunkSize := (len(watermarkRunes) + len(eligibleIndexes) - 1) / len(eligibleIndexes)
	chunks := make(map[int]string, len(eligibleIndexes))
	cursor := 0

	for _, index := range eligibleIndexes {
		if cursor >= len(watermarkRunes) {
			break
		}

		end := cursor + chunkSize
		if end > len(watermarkRunes) {
			end = len(watermarkRunes)
		}

		chunks[index] = string(watermarkRunes[cursor:end])
		cursor = end
	}

	var result strings.Builder

	for i, line := range lines {
		result.WriteString(line)
		if chunk, ok := chunks[i]; ok {
			result.WriteString(chunk)
		}

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func isWatermarkEligibleLine(trimmed string) bool {
	if trimmed == "" {
		return false
	}

	if strings.HasPrefix(trimmed, "```") {
		return false
	}

	// Skip thematic breaks like --- or ***.
	if strings.Trim(trimmed, "-*_ ") == "" && utf8.RuneCountInString(strings.ReplaceAll(trimmed, " ", "")) >= 3 {
		return false
	}

	// Skip markdown table rows so the parser can still recognize the table.
	if strings.Contains(trimmed, "|") {
		return false
	}

	return true
}
