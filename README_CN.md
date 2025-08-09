# RainbowLog

![](https://ramboll.wang/img/RainbowLog_LOGO_SD.png)

[[English](./README.md)

简单、可配置、结构化的 Go 日志库。

RainbowLog 的 API 设计旨在提供出色的开发者体验和卓越的性能。
其独特的链式 API 允许 RainbowLog 通过避免分配和反射来写入 JSON 日志事件。

为了保持代码库和 API 的简洁性，RainbowLog 仅专注于高效的结构化日志记录。
通过提供的（但效率较低的）`rainbowlog.ConsolePacker` 可以实现控制台上的漂亮日志输出。

![pretty.png](./pretty.png)

## 特性

- **高性能**: 零内存分配，无反射的日志记录
- **结构化日志**: 支持 JSON 和文本格式的日志输出
- **灵活配置**: 支持多种配置方式（代码配置、配置文件等）
- **多输出支持**: 可同时输出到多个目标（文件、标准输出等）
- **日志级别控制**: 支持 Debug、Info、Warn、Error、Fatal、Panic 等级别
- **模块化标签**: 支持为不同模块添加标签便于日志分类
- **钩子机制**: 支持自定义钩子函数处理日志事件
- **子日志记录器**: 支持创建继承父级配置的子日志记录器
- **错误堆栈跟踪**: 支持错误堆栈信息输出
- **Caller 信息**: 支持记录日志产生的文件和行号

## 安装

```shell
go get -u github.com/rambollwong/rainbowlog
```

## 快速开始

### 使用Global Logger

#### 使用默认选项的Global Logger

```go
package main

import (
	"github.com/rambollwong/rainbowlog/log"
)

func main() {
	log.UseDefault()

	log.Logger.Info().Msg("Hello world!").Done()
}

// Output: {"_TIME_":"2024-02-19 19:50:09.008","_LEVEL_":"INFO","_CALLER_":"/path/to/main.go:10","message":"Hello world!"}
```

> 注意: 默认情况下，日志写入到 `os.Stderr`

#### 使用 Rainbow 默认选项的Global Logger

```go
package main

import (
	"errors"

	"github.com/rambollwong/rainbowlog/log"
)

func main() {
	log.UseRainbowDefault()

	log.Info().Msg("Hello world!").Done()
	log.Debug().WithLabels("MODEL1").Msg("Something debugging...").Done()
	log.Warn().WithLabels("MODEL2", "SERVICE1").Msg("Something warning!").Int("IntegerValue", 888).Done()
	log.Error().Msg("failed to do something").Err(errors.New("something wrong")).Done()
	log.Fatal().Msg("fatal to do something").Done()
}
```

输出:

![pretty.png](./pretty.png)

#### 使用配置文件的Global Logger

RainbowLog 支持基于配置文件设置日志记录器选项。
如果要使用配置文件，您需要确保配置文件包含 RainbowLog 配置项。
RainbowLog 支持三种格式的配置文件：`.yaml`|`.json`|`.toml`。有关具体配置模板，请参见 [`config`](./config) 包中的相应文件。

假设我们已准备了一个配置文件 `rainbowlog.yaml` 并将其放置在与执行文件相同的目录中，那么：

```go
package main

import (
	"errors"

	"github.com/rambollwong/rainbowlog/log"
)

func main() {
	log.UseDefaultConfigFile()

	log.Info().Msg("Hello world!").Done()
	log.Debug().WithLabels("MODEL1").Msg("Something debugging...").Done()
	log.Warn().WithLabels("MODEL2", "SERVICE1").Msg("Something warning!").Int("IntegerValue", 888).Done()
	log.Error().Msg("failed to do something").Err(errors.New("something wrong")).Done()
	log.Fatal().Msg("fatal to do something").Done()
}
```

如果要使用 `.json` 或 `.toml` 类型的配置文件，只需修改 `log.DefaultConfigFileName`，例如：

```go
log.DefaultConfigFileName = "rainbowlog.json"
```

或者

```go
log.DefaultConfigFileName = "rainbowlog.toml"
```

如果您还想指定配置文件所在的目录，只需修改 `log.DefaultConfigFilePath`：

```go
log.DefaultConfigFilePath = "/path/of/config/files"
```

> 注意: 修改 `log.DefaultConfigFileName` 和 `log.DefaultConfigFilePath` 需要在 `log.UseDefaultConfigFile()` 之前执行，否则不会生效。

#### 使用自定义选项的Global Logger

如果您想为全局日志记录器使用自定义选项，我们保留了 `log.UseCustomOptions(opts ...Option)` API 来实现。
支持的 `Option` 详见 [option.go](./option.go)。

### 自定义Logger

如果您不想使用Global Logger，可以通过 `New` 方法初始化一个 `Logger` 实例。
`New` 方法接收 `Option` 参数。
支持的 `Option` 详见 [option.go](./option.go)。

```go
package main

import (
	"errors"
	"path/filepath"

	"github.com/rambollwong/rainbowlog"
)

func main() {
	DefaultConfigFileName := "rainbowlog.yaml"
	DefaultConfigFilePath := "/path/of/config/files"

	logger := rainbowlog.New(
		rainbowlog.WithDefault(),
		rainbowlog.WithConfigFile(filepath.Join(DefaultConfigFilePath, DefaultConfigFileName)),
	)

	logger.Info().Msg("Hello world!").Done()
	logger.Debug().WithLabels("MODEL1").Msg("Something debugging...").Done()
	logger.Warn().WithLabels("MODEL2", "SERVICE1").Msg("Something warning!").Int("IntegerValue", 888).Done()
	logger.Error().Msg("failed to do something").Err(errors.New("something wrong")).Done()
	logger.Fatal().Msg("fatal to do something").Done()
}
```

### SubLogger

SubLogger支持允许您创建一个继承自父日志记录器的实例，并在需要时重置某些 `Option`。
例如，在子模块中使用与父Logger不同的 `LABEL` 的场景。

```go
package main

import (
	"os"

	"github.com/rambollwong/rainbowlog"
	"github.com/rambollwong/rainbowlog/level"
)

func main() {
	logger := rainbowlog.New(
		rainbowlog.WithDefault(),
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stderr),
		rainbowlog.WithCallerMarshalFunc(nil),
		rainbowlog.WithLevel(level.Info),
		rainbowlog.WithLabels("ROOT"),
	)

	logger.Debug().Msg("Hello world!").Done()
	logger.Info().Msg("Hello world!").Done()

	subLogger := logger.SubLogger(
		rainbowlog.WithLevel(level.Debug),
		rainbowlog.WithLabels("SUBMODULE"),
	)

	subLogger.Debug().Msg("Hello world!").Done()
	subLogger.Info().Msg("Hello world!").Done()
}

// Output:
// {"_TIME_":"2024-02-21 11:28:02.150","_LEVEL_":"INFO","_LABEL_":"ROOT","message": "Hello world!"}
// {"_TIME_":"2024-02-21 11:28:02.150","_LEVEL_":"DEBUG","_LABEL_":"SUBMODULE","message":"Hello world!"}
// {"_TIME_":"2024-02-21 11:28:02.150","_LEVEL_":"INFO","_LABEL_":"SUBMODULE","message":"Hello world!"}

```

### 修改时间输出格式

RainbowLog 的默认时间格式是 `2006-01-02 15:04:05.000`，
您可以通过 `WithTimeFormat(timeFormat string)` 选项修改此格式，
格式可以是符合 golang 时间格式规则的字符串，
也可以是 `UNIX` 或 `UNIXMS` 或 `UNIXMICRO` 或 `UNIXNANO`，
它们分别代表 `Unix()` 或 `UnixMilli()` 或 `UnixMicro()` 或 `UnixNano()` 的返回值，用于输出 `time.Time`。

```go
package main

import (
	"os"

	"github.com/rambollwong/rainbowlog"
)

func main() {
	logger := rainbowlog.New(
		rainbowlog.WithDefault(),
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stderr),
		rainbowlog.WithTimeFormat(rainbowlog.TimeFormatUnix),
	)

	logger.Info().Msg("Hello world!").Done()
}

// Output:{"_TIME_":1708346689,"_LEVEL_":"INFO","_CALLER_":"main.go:16","message":"Hello world!"}
```

### Hook

```go
package main

import (
	"fmt"
	"os"

	"github.com/rambollwong/rainbowlog"
	"github.com/rambollwong/rainbowlog/level"
)

func main() {
	var hook rainbowlog.HookFunc = func(r rainbowlog.Record, level level.Level, message string) {
		fmt.Printf("hook: %s, %s\n", level.String(), message)
	}
	logger := rainbowlog.New(
		rainbowlog.WithDefault(),
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stderr),
		rainbowlog.WithCallerMarshalFunc(nil),
		rainbowlog.AppendsHooks(hook),
		rainbowlog.WithLevel(level.Info),
	)

	logger.Debug().Msg("Hello world!").Done()
	logger.Info().Msg("Hello world!").Done()
}

// Output: 
// hook: info, Hello world!
// {"_TIME_":"2024-02-21 11:42:17.592","_LEVEL_":"INFO","message":"Hello world!"}
```

### 高级用法

#### BufferedWriter

RainbowLog 提供了BufferedWriter功能，可以显著提高日志写入性能：

```go
bufferedWriter := rainbowlog.BufferedWriter(os.Stdout, 4096)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, bufferedWriter),
)
```

#### SyncWriter

对于非线程安全的Writer，可以使用SyncWriter包装：

```go
syncWriter := rainbowlog.SyncWriter(os.Stdout)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, syncWriter),
)
```

**注意：POSIX 和 Windows 操作系统本身就是写入安全的，无需使用SyncWriter来包装！**

#### MultiWriter

支持同时写入多个目标：

```go
multiWriter := rainbowlog.MultiLevelWriter(os.Stdout, fileWriter)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, multiWriter),
)
```

### 性能优化建议

1. **使用缓冲写入器**: 对于频繁的日志写入操作，使用缓冲写入器可以显著提高性能
2. **合理设置日志级别**: 在生产环境中适当提高日志级别，避免过多的调试日志影响性能
3. **避免在热路径中记录过多信息**: 在关键性能路径上尽量减少日志记录的内容和频率
4. **使用对象池**: RainbowLog 内部使用对象池来减少内存分配，确保正确使用 `Done()` 方法释放记录对象

### 更多

当然，我们还提供更多功能和特性，期待您的探索和发现！

## 联系我们

- 邮箱：`ramboll.wong@hotmail.com`
- Telegram 技术交流群：[[点击加入]](https://t.me/+ovN79ozHG4c1YWNl)
- 博客：[Ramboll's Blog](https://ramboll.wang)

## 打赏支持

如果你喜欢这个项目，欢迎请作者喝杯柠檬水 ☕️，你的支持是我持续更新的动力！

- 微信支付：

<img src="https://ramboll.wang/img/wechat_pay.jpg" alt="WeChat Pay" width="200"/>

- 支付宝：

<img src="https://ramboll.wang/img/ali_pay.jpg" alt="Alipay" width="200"/>