package sensorparser

import (
	"bufio"
	"os"
	"strings"
)

// GetSystemDistro 检测当前操作系统的发行版
// 返回小写的发行版名称，如 "ubuntu", "openeuler" 等
func GetSystemDistro() string {
	// 尝试读取/etc/os-release文件，这是大多数现代Linux发行版的标准
	file, err := os.Open("/etc/os-release")
	if err != nil {
		// 如果无法打开该文件，尝试其他文件
		return tryAlternativeFiles()
	}
	defer file.Close()

	// 读取文件寻找ID字段
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 首先检查ID字段
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			// 去除可能的引号
			id = strings.Trim(id, "\"'")
			return strings.ToLower(id)
		}
	}

	// 如果没有找到ID，尝试其他方法
	return tryAlternativeFiles()
}

// 尝试其他文件或方法来识别系统
func tryAlternativeFiles() string {
	// 尝试/etc/lsb-release (Ubuntu常用)
	if distro := readLsbRelease(); distro != "" {
		return distro
	}

	// 尝试检查常见的发行版特定文件
	if fileExists("/etc/openEuler-release") {
		return "openeuler"
	}
	if fileExists("/etc/ubuntu-release") || fileExists("/etc/debian_version") {
		return "ubuntu"
	}

	// 默认返回空字符串，表示无法识别
	return ""
}

// 从/etc/lsb-release读取发行版信息
func readLsbRelease() string {
	file, err := os.Open("/etc/lsb-release")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DISTRIB_ID=") {
			id := strings.TrimPrefix(line, "DISTRIB_ID=")
			id = strings.Trim(id, "\"'")
			return strings.ToLower(id)
		}
	}

	return ""
}

// 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
