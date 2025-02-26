package sensorparser

import (
	"bufio"
	"os"
	"testing"
)

func TestGenericParser(t *testing.T) {
	// 测试通用解析器
	parser := &GenericParser{}

	// 测试Ubuntu传感器输出
	ubuntuFile, err := os.Open("testdata/sensors_ubuntu.txt")
	if err != nil {
		t.Fatalf("无法打开Ubuntu测试文件: %v", err)
	}
	defer ubuntuFile.Close()

	ubuntuResult, err := parser.Parse(ubuntuFile)
	if err != nil {
		t.Fatalf("解析Ubuntu传感器数据失败: %v", err)
	}

	// 验证Ubuntu结果
	if len(ubuntuResult.Devices) != 4 {
		t.Errorf("Ubuntu设备数量错误: 期望4, 得到%d", len(ubuntuResult.Devices))
	}

	// 检查第一个设备
	if ubuntuResult.Devices[0].Name != "nvme-pci-0200" {
		t.Errorf("设备名称错误: 期望nvme-pci-0200, 得到%s", ubuntuResult.Devices[0].Name)
	}

	if ubuntuResult.Devices[0].Adapter != "PCI adapter" {
		t.Errorf("适配器错误: 期望PCI adapter, 得到%s", ubuntuResult.Devices[0].Adapter)
	}

	if len(ubuntuResult.Devices[0].Sensors) != 1 {
		t.Errorf("传感器数量错误: 期望1, 得到%d", len(ubuntuResult.Devices[0].Sensors))
	}
}

func TestOpenEulerParser(t *testing.T) {
	// 测试OpenEuler解析器
	parser := &OpenEulerParser{}

	// 测试openEuler传感器输出
	openEulerFile, err := os.Open("testdata/sensors_openEuler.txt")
	if err != nil {
		t.Fatalf("无法打开openEuler测试文件: %v", err)
	}
	defer openEulerFile.Close()

	openEulerResult, err := parser.Parse(openEulerFile)
	if err != nil {
		t.Fatalf("解析openEuler传感器数据失败: %v", err)
	}

	// 验证openEuler结果
	if len(openEulerResult.Devices) != 1 {
		t.Errorf("openEuler设备数量错误: 期望1, 得到%d", len(openEulerResult.Devices))
		return // 避免后续访问不存在的设备
	}

	// 检查设备
	if openEulerResult.Devices[0].Name != "k10temp-pci-00c3" {
		t.Errorf("设备名称错误: 期望k10temp-pci-00c3, 得到%s", openEulerResult.Devices[0].Name)
	}

	if len(openEulerResult.Devices[0].Sensors) != 2 {
		t.Errorf("传感器数量错误: 期望2, 得到%d", len(openEulerResult.Devices[0].Sensors))
		return // 避免后续访问不存在的传感器
	}

	// 检查传感器值
	if len(openEulerResult.Devices[0].Sensors) > 0 {
		tctl := openEulerResult.Devices[0].Sensors[0]
		if tctl.Name != "Tctl" || tctl.Value != 49.0 {
			t.Errorf("Tctl传感器错误: 期望49.0, 得到%f", tctl.Value)
		}
	}
}

func TestUbuntuParser(t *testing.T) {
	// 测试Ubuntu解析器
	parser := &UbuntuParser{}

	// 测试Ubuntu传感器输出
	ubuntuFile, err := os.Open("testdata/sensors_ubuntu.txt")
	if err != nil {
		t.Fatalf("无法打开Ubuntu测试文件: %v", err)
	}
	defer ubuntuFile.Close()

	ubuntuResult, err := parser.Parse(ubuntuFile)
	if err != nil {
		t.Fatalf("解析Ubuntu传感器数据失败: %v", err)
	}

	// 验证Ubuntu结果
	if len(ubuntuResult.Devices) != 4 {
		t.Errorf("Ubuntu设备数量错误: 期望4, 得到%d", len(ubuntuResult.Devices))
	}

	// 检查第一个设备
	if ubuntuResult.Devices[0].Name != "nvme-pci-0200" {
		t.Errorf("设备名称错误: 期望nvme-pci-0200, 得到%s", ubuntuResult.Devices[0].Name)
	}
}

func TestNewParser(t *testing.T) {
	// 测试工厂方法
	ubuntuParser := NewParser("ubuntu")
	openEulerParser := NewParser("openEuler")
	genericParser := NewParser("unknown")

	// 验证类型
	_, isUbuntuParser := ubuntuParser.(*UbuntuParser)
	if !isUbuntuParser {
		t.Errorf("期望UbuntuParser类型")
	}

	_, isOpenEulerParser := openEulerParser.(*OpenEulerParser)
	if !isOpenEulerParser {
		t.Errorf("期望OpenEulerParser类型")
	}

	_, isGenericParser := genericParser.(*GenericParser)
	if !isGenericParser {
		t.Errorf("期望GenericParser类型")
	}
}

func TestLegacyParse(t *testing.T) {
	// 测试向后兼容的Parse函数

	// 测试Ubuntu传感器输出
	ubuntuFile, err := os.Open("testdata/sensors_ubuntu.txt")
	if err != nil {
		t.Fatalf("无法打开Ubuntu测试文件: %v", err)
	}
	defer ubuntuFile.Close()

	ubuntuResult, err := Parse(ubuntuFile)
	if err != nil {
		t.Fatalf("解析Ubuntu传感器数据失败: %v", err)
	}

	// 验证Ubuntu结果
	if len(ubuntuResult.Devices) != 4 {
		t.Errorf("Ubuntu设备数量错误: 期望4, 得到%d", len(ubuntuResult.Devices))
	}
}

func TestDebugOpenEulerData(t *testing.T) {
	// 打印openEuler测试数据
	openEulerFile, err := os.Open("testdata/sensors_openEuler.txt")
	if err != nil {
		t.Fatalf("无法打开openEuler测试文件: %v", err)
	}
	defer openEulerFile.Close()

	scanner := bufio.NewScanner(openEulerFile)
	t.Log("OpenEuler测试数据内容:")
	for scanner.Scan() {
		t.Logf("行: %q", scanner.Text())
	}
}
