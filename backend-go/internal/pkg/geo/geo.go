package geo

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	searcher *xdb.Searcher
	once     sync.Once
)

// Init 初始化 IP 地理位置查询（全量缓存模式，线程安全）
func Init(dbPath string) {
	once.Do(func() {
		cBuff, err := os.ReadFile(dbPath)
		if err != nil {
			log.Printf("[geo] 加载 ip2region 数据库失败: %v", err)
			return
		}
		s, err := xdb.NewWithBuffer(xdb.IPv4, cBuff)
		if err != nil {
			log.Printf("[geo] 初始化 searcher 失败: %v", err)
			return
		}
		searcher = s
		log.Printf("[geo] ip2region 数据库加载成功 (%d bytes)", len(cBuff))
	})
}

// Lookup 查询 IP 地理位置，返回 (country, region)
// ip2region 返回格式: "国家|区域|省份|城市|ISP"
func Lookup(ip string) (country, region string) {
	if searcher == nil || ip == "" {
		return "", ""
	}
	// 跳过本地/内网 IP
	if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "192.168.") || ip == "::1" || ip == "0.0.0.0" {
		return "本地", "内网"
	}
	// IPv6 暂不支持（使用的是 v4 数据库）
	if strings.Contains(ip, ":") {
		return "", ""
	}
	result, err := searcher.SearchByStr(ip)
	if err != nil {
		return "", ""
	}
	return parseResult(result)
}

// parseResult 解析 "国家|区域|省份|城市|ISP" 格式
func parseResult(result string) (country, region string) {
	parts := strings.SplitN(result, "|", 5)
	if len(parts) < 4 {
		return "", ""
	}
	country = parts[0]
	if country == "0" {
		country = ""
	}
	// 省份+城市 作为 region
	province := parts[2]
	city := parts[3]
	if province == "0" {
		province = ""
	}
	if city == "0" {
		city = ""
	}
	switch {
	case province != "" && city != "" && province != city:
		region = province + " " + city
	case province != "":
		region = province
	case city != "":
		region = city
	}
	return country, region
}
