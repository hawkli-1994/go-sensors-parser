// Package sensorparser 提供解析Linux sensors命令输出的功能
package sensorparser

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Sensor 表示一个传感器的信息
type Sensor struct {
	Name     string
	Value    float64
	Unit     string
	Low      *float64
	High     *float64
	Critical *float64
}

// Device 表示一个设备，包含多个传感器
type Device struct {
	Name    string
	Adapter string
	Sensors []Sensor
}

// Result 表示解析结果，包含多个设备
type Result struct {
	Devices []Device
}

// SensorParser 定义传感器解析器接口
type SensorParser interface {
	Parse(r io.Reader) (*Result, error)
}

// 创建适合的解析器
func NewParser(distro string) SensorParser {
	switch strings.ToLower(distro) {
	case "openeuler":
		return &OpenEulerParser{}
	case "ubuntu":
		return &UbuntuParser{}
	default:
		// 默认使用通用解析器
		return &GenericParser{}
	}
}

// GenericParser 通用解析器实现
type GenericParser struct{}

// Parse 实现SensorParser接口
func (p *GenericParser) Parse(r io.Reader) (*Result, error) {
	scanner := bufio.NewScanner(r)
	result := &Result{
		Devices: []Device{},
	}

	var currentDevice Device
	var currentSensor *Sensor

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			// 空行表示设备信息结束
			if currentDevice.Name != "" {
				if currentSensor != nil {
					currentDevice.Sensors = append(currentDevice.Sensors, *currentSensor)
					currentSensor = nil
				}
				result.Devices = append(result.Devices, currentDevice)
				currentDevice = Device{Sensors: []Sensor{}}
			}
			continue
		}

		if currentDevice.Name == "" && !strings.Contains(line, ":") {
			// 新设备开始 - 确保设备名行不包含冒号
			currentDevice = Device{
				Name:    line,
				Sensors: []Sensor{},
			}
		} else if strings.HasPrefix(line, "Adapter:") {
			// 适配器信息
			currentDevice.Adapter = strings.TrimSpace(strings.TrimPrefix(line, "Adapter:"))
		} else {
			// 传感器信息
			if strings.HasPrefix(line, " ") && currentSensor != nil {
				// 这是传感器的续行（比如临界值信息）
				sensor, err := parseSensorWithLimits(currentSensor.Name, currentSensor.Value, line)
				if err == nil && sensor != nil {
					*currentSensor = *sensor
				}
			} else {
				// 如果有未添加的传感器，先添加到设备中
				if currentSensor != nil {
					currentDevice.Sensors = append(currentDevice.Sensors, *currentSensor)
					currentSensor = nil
				}
				
				// 解析新的传感器行
				sensor, err := parseSensorLine(line)
				if err == nil && sensor != nil {
					currentSensor = sensor
				}
			}
		}
	}

	// 添加最后一个设备和传感器（如果有）
	if currentDevice.Name != "" {
		if currentSensor != nil {
			currentDevice.Sensors = append(currentDevice.Sensors, *currentSensor)
		}
		result.Devices = append(result.Devices, currentDevice)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Parse 提供便捷的解析方法，使用通用解析器
func Parse(r io.Reader) (*Result, error) {
	parser := &GenericParser{}
	return parser.Parse(r)
}

// parseSensorLine 解析单行传感器信息
func parseSensorLine(line string) (*Sensor, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, nil // 不是传感器行
	}

	name := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])
	
	// 确保值部分包含数字
	hasDigit := false
	for _, r := range valueStr {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}
	
	if !hasDigit {
		return nil, nil // 不是有效的传感器值
	}

	// 处理多行值的情况
	if strings.Contains(valueStr, "(") {
		return parseSensorWithLimits(name, valueStr)
	}

	// 处理简单值的情况
	return parseSensorSimple(name, valueStr)
}

// parseSensorSimple 解析简单的传感器值
func parseSensorSimple(name, valueStr string) (*Sensor, error) {
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

// parseSensorWithLimits 解析带有限制值的传感器
func parseSensorWithLimits(name string, valueOrSensor interface{}, limitLineOpt ...string) (*Sensor, error) {
	var sensor *Sensor

	// 根据传入参数类型决定如何处理
	switch v := valueOrSensor.(type) {
	case string:
		// 原来的情况：解析包含值和限制的字符串
		valueStr := v

		// 分离主值和限制值
		mainParts := strings.SplitN(valueStr, "(", 2)
		if len(mainParts) < 1 {
			return nil, nil
		}

		// 解析主值
		var err error
		sensor, err = parseSensorSimple(name, mainParts[0])
		if err != nil {
			return nil, err
		}

		// 解析限制值
		if len(mainParts) > 1 {
			limitsPart := mainParts[1]
			// 移除结尾的括号
			limitsPart = strings.TrimSuffix(limitsPart, ")")

			// 处理限制值部分
			processLimits(sensor, limitsPart)
		}

	case float64:
		// 新情况：已有传感器值，解析限制行
		if len(limitLineOpt) == 0 {
			return nil, nil
		}

		// 创建包含基本信息的传感器
		sensor = &Sensor{
			Name:  name,
			Value: v,
		}

		// 处理限制值行
		limitLine := limitLineOpt[0]
		processLimits(sensor, limitLine)
	}

	return sensor, nil
}

// 处理限制值
func processLimits(sensor *Sensor, limitsStr string) {
	// 解析low, high, critical值
	limits := strings.Split(limitsStr, ",")
	for _, limit := range limits {
		limit = strings.TrimSpace(limit)

		if strings.HasPrefix(limit, "low") {
			lowStr := extractValue(limit)
			if lowStr != "" {
				low, err := strconv.ParseFloat(lowStr, 64)
				if err == nil {
					sensor.Low = &low
				}
			}
		} else if strings.HasPrefix(limit, "high") {
			highStr := extractValue(limit)
			if highStr != "" {
				high, err := strconv.ParseFloat(highStr, 64)
				if err == nil {
					sensor.High = &high
				}
			}
		} else if strings.HasPrefix(limit, "crit") {
			critStr := extractValue(limit)
			if critStr != "" {
				crit, err := strconv.ParseFloat(critStr, 64)
				if err == nil {
					sensor.Critical = &crit
				}
			}
		}
	}
}

// extractValue 从"key = value"格式的字符串中提取值
func extractValue(s string) string {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return ""
	}

	valueStr := strings.TrimSpace(parts[1])
	// 移除可能的+号和单位
	valueStr = strings.TrimPrefix(valueStr, "+")
	valueStr = strings.Fields(valueStr)[0]

	return valueStr
}
