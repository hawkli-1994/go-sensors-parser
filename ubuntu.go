package sensorparser

import (
	"bufio"
	"io"
	"strings"
)

// UbuntuParser 针对Ubuntu系统的解析器
type UbuntuParser struct{}

// Parse 实现SensorParser接口
func (p *UbuntuParser) Parse(r io.Reader) (*Result, error) {
	scanner := bufio.NewScanner(r)
	result := &Result{
		Devices: []Device{},
	}

	var currentDevice Device

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			// 空行表示设备信息结束
			if currentDevice.Name != "" {
				result.Devices = append(result.Devices, currentDevice)
			}
			currentDevice = Device{Sensors: []Sensor{}}
			continue
		}

		if currentDevice.Name == "" {
			// 新设备开始
			currentDevice = Device{
				Name:    line,
				Sensors: []Sensor{},
			}
		} else if strings.HasPrefix(line, "Adapter:") {
			// 适配器信息
			currentDevice.Adapter = strings.TrimSpace(strings.TrimPrefix(line, "Adapter:"))
		} else {
			// 传感器信息 - Ubuntu格式通常包含更多限制值
			sensor, err := parseUbuntuSensorLine(line)
			if err != nil || sensor == nil {
				continue // 忽略解析错误
			}
			currentDevice.Sensors = append(currentDevice.Sensors, *sensor)
		}
	}

	// 添加最后一个设备（如果有）
	if currentDevice.Name != "" {
		result.Devices = append(result.Devices, currentDevice)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}


// parseUbuntuSensorLine 解析Ubuntu格式的传感器行
func parseUbuntuSensorLine(line string) (*Sensor, error) {
	return parseSensorLine(line)
}
