package utils

import (
	"strings"
)

// UserAgentInfo User-Agent解析信息
type UserAgentInfo struct {
	DeviceType string
	OS         string
	Browser    string
}

// ParseUserAgent 解析User-Agent字符串
func ParseUserAgent(userAgentString string) UserAgentInfo {
	if userAgentString == "" {
		return UserAgentInfo{
			DeviceType: "Unknown",
			OS:         "Unknown",
			Browser:    "Unknown",
		}
	}

	// 简化的User-Agent解析逻辑
	deviceType := "Unknown"
	os := "Unknown"
	browser := "Unknown"

	// 设备类型判断
	if strings.Contains(strings.ToLower(userAgentString), "mobile") {
		deviceType = "Mobile"
	} else if strings.Contains(strings.ToLower(userAgentString), "tablet") {
		deviceType = "Tablet"
	} else {
		deviceType = "PC"
	}

	// 操作系统判断
	if strings.Contains(strings.ToLower(userAgentString), "windows") {
		os = "Windows"
	} else if strings.Contains(strings.ToLower(userAgentString), "macintosh") {
		os = "Mac OS"
	} else if strings.Contains(strings.ToLower(userAgentString), "linux") {
		os = "Linux"
	} else if strings.Contains(strings.ToLower(userAgentString), "android") {
		os = "Android"
	} else if strings.Contains(strings.ToLower(userAgentString), "iphone") || strings.Contains(strings.ToLower(userAgentString), "ipad") {
		os = "iOS"
	}

	// 浏览器判断
	if strings.Contains(strings.ToLower(userAgentString), "chrome") {
		browser = "Chrome"
	} else if strings.Contains(strings.ToLower(userAgentString), "firefox") {
		browser = "Firefox"
	} else if strings.Contains(strings.ToLower(userAgentString), "safari") {
		browser = "Safari"
	} else if strings.Contains(strings.ToLower(userAgentString), "edge") {
		browser = "Edge"
	}

	return UserAgentInfo{
		DeviceType: deviceType,
		OS:         os,
		Browser:    browser,
	}
}
