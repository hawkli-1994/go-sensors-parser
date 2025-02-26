package sensorparser

import (
	"os"
	"testing"
)

func TestGetSystemDistro(t *testing.T) {
	// 测试Ubuntu发行版
	os.Setenv("ID", "ubuntu")
	os.Setenv("ID_LIKE", "debian")

	distro := GetSystemDistro()
	if distro != "ubuntu" {
		t.Errorf("期望Ubuntu发行版, 得到%s", distro)
	}
}
