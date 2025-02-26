package sensorparser

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// OpenEulerParser 针对OpenEuler系统的解析器
type OpenEulerParser struct{}

// Parse 实现SensorParser接口
func (p *OpenEulerParser) Parse(r io.Reader) (*Result, error) {
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
				currentDevice = Device{Sensors: []Sensor{}}
			}
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
			// 传感器信息 - OpenEuler格式通常更简单
			sensor, err := parseOpenEulerSensorLine(line)
			if err != nil || sensor == nil {
				continue // 忽略解析错误
			}
			currentDevice.Sensors = append(currentDevice.Sensors, *sensor)
		}
	}

	// 添加最后一个设备（如果有）
	if currentDevice.Name != "" && len(currentDevice.Sensors) > 0 {
		result.Devices = append(result.Devices, currentDevice)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// parseOpenEulerSensorLine 解析OpenEuler格式的传感器行
func parseOpenEulerSensorLine(line string) (*Sensor, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, nil // 不是传感器行
	}

	name := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	// 提取纯数值部分
	var rawValue string
	var unit string

	// 查找第一个非数字、非加减号、非小数点的字符位置
	i := 0
	for i < len(valueStr) {
		if valueStr[i] == '+' || valueStr[i] == '-' || valueStr[i] == '.' || (valueStr[i] >= '0' && valueStr[i] <= '9') {
			i++
		} else {
			break
		}
	}

	if i > 0 {
		rawValue = valueStr[:i]
		if i < len(valueStr) {
			unit = strings.TrimSpace(valueStr[i:])
		}
	} else {
		// 没有找到有效数字
		return nil, nil
	}

	// 移除可能的+号
	rawValue = strings.TrimPrefix(rawValue, "+")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		Name:  name,
		Value: value,
		Unit:  unit,
	}, nil
}
