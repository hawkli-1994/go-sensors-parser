# go-sensors-parser

![Go Version](https://img.shields.io/badge/Go-1.16+-blue)
![License](https://img.shields.io/badge/License-MIT-green)

这个开源库使用 Golang 解析 Linux `sensors` 命令的输出。由于原始传感器数据难以集成到应用程序中，本项目将其简化，让开发者可以轻松地在自己的项目中使用 Linux 系统的硬件传感器数据。

项目充分利用 Golang 的性能和并发特性，提供了简单易用的 API，适用于系统监控和硬件相关项目。

## 特性

- ✅ 解析多种 Linux 发行版（Ubuntu, openEuler 等）的 sensors 输出格式
- ✅ 自动识别当前系统类型并选择合适的解析器
- ✅ 结构化的传感器数据输出，易于在应用中使用
- ✅ 支持所有常见的传感器类型（温度、风扇速度、电压等）
- ✅ 支持临界值、高/低阈值等传感器限制信息

## 安装

```bash
go get github.com/hawkli-1994/go-sensors-parser
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "os/exec"
    
    "github.com/hawkli-1994/go-sensors-parser"
)

func main() {
    // 执行 sensors 命令
    cmd := exec.Command("sensors")
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("执行 sensors 命令失败:", err)
        return
    }
    
    // 自动检测当前系统并选择合适的解析器
    distro := sensorparser.GetSystemDistro()
    parser := sensorparser.NewParser(distro)
    
    // 解析传感器数据
    result, err := parser.Parse(bytes.NewReader(output))
    if err != nil {
        fmt.Println("解析传感器数据失败:", err)
        return
    }
    
    // 打印所有设备和传感器信息
    for _, device := range result.Devices {
        fmt.Printf("设备: %s (适配器: %s)\n", device.Name, device.Adapter)
        for _, sensor := range device.Sensors {
            fmt.Printf("  %s: %.1f %s\n", sensor.Name, sensor.Value, sensor.Unit)
            
            // 打印临界值（如果有）
            if sensor.Critical != nil {
                fmt.Printf("    临界值: %.1f %s\n", *sensor.Critical, sensor.Unit)
            }
        }
        fmt.Println()
    }
}
```

### 使用特定解析器

```go
// 如果您知道当前系统的类型，可以直接选择解析器
ubuntuParser := &sensorparser.UbuntuParser{}
// 或者
openEulerParser := &sensorparser.OpenEulerParser{}
```

## API 文档

### 主要结构

```go
// Sensor 表示一个传感器的信息
type Sensor struct {
    Name     string
    Value    float64
    Unit     string
    Low      *float64    // 可选的低阈值
    High     *float64    // 可选的高阈值
    Critical *float64    // 可选的临界值
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
```

### 主要方法

- `GetSystemDistro()` - 检测当前 Linux 发行版类型
- `NewParser(distro string)` - 根据发行版类型创建合适的解析器
- `Parse(r io.Reader) (*Result, error)` - 通用解析方法，使用默认解析器
- `parser.Parse(r io.Reader) (*Result, error)` - 使用特定解析器解析数据

## 支持的系统

目前支持以下 Linux 发行版的 sensors 输出格式：

- Ubuntu/Debian
- openEuler
- 通用 Linux (Generic)

其他 Linux 发行版也可能被支持，但未经过专门测试。

## 贡献

欢迎贡献代码、报告问题或提出新功能建议！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

感谢所有 Linux 系统监控和传感器管理工具的开发者，特别是 lm-sensors 项目。

---

欢迎加入我们，一起改进这个开源项目！如果您有任何问题或反馈，请在 GitHub 上开 Issue 或直接联系我们。
