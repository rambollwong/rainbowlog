# RainbowLog

![](https://ramboll.wang/img/RainbowLog_LOGO_SD.png)

[[中文]](./README_CN.md)

Simple, configurable, structured Go logging library.

RainbowLog's API is designed to provide an excellent developer experience and outstanding performance.
Its unique chained API allows RainbowLog to write JSON log events by avoiding allocations and reflection.

To keep the codebase and API simple, RainbowLog focuses only on efficient structured logging.
Pretty log output on the console can be achieved through the provided (but less efficient) `rainbowlog.ConsolePacker`.

![pretty.png](./pretty.png)

## Features

- **High Performance**: Zero memory allocation, no reflection logging
- **Structured Logging**: Supports JSON and text format log output
- **Flexible Configuration**: Supports multiple configuration methods (code configuration, configuration files, etc.)
- **Multi-output Support**: Can output to multiple targets simultaneously (files, standard output, etc.)
- **Log Level Control**: Supports Debug, Info, Warn, Error, Fatal, Panic levels
- **Modular Labels**: Supports adding labels to different modules for easy log classification
- **Hook Mechanism**: Supports custom hook functions to handle log events
- **Sub Logger**: Supports creating sub loggers that inherit parent configuration
- **Error Stack Tracing**: Supports error stack information output
- **Caller Information**: Supports recording the file and line number where the log was generated

## Installation

```shell
go get -u github.com/rambollwong/rainbowlog
```

## Quick Start

### Using Global Logger

#### Global Logger with Default Options

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

> Note: By default, logs are written to `os.Stderr`

#### Global Logger with Rainbow Default Options

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

Output:

![pretty.png](./pretty.png)

#### Global Logger with Configuration File

RainbowLog supports setting logger options based on configuration files.
To use a configuration file, you need to ensure that the configuration file contains RainbowLog configuration items.
RainbowLog supports three formats of configuration files: `.yaml`|`.json`|`.toml`. For specific configuration templates, please refer to the corresponding files in the [`config`](./config) package.

Assuming we have prepared a configuration file `rainbowlog.yaml` and placed it in the same directory as the executable file:

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

To use `.json` or `.toml` type configuration files, just modify `log.DefaultConfigFileName`, for example:

```go
log.DefaultConfigFileName = "rainbowlog.json"
```

Or

```go
log.DefaultConfigFileName = "rainbowlog.toml"
```

If you also want to specify the directory where the configuration file is located, just modify `log.DefaultConfigFilePath`:

```go
log.DefaultConfigFilePath = "/path/of/config/files"
```

> Note: Modifying `log.DefaultConfigFileName` and `log.DefaultConfigFilePath` needs to be executed before `log.UseDefaultConfigFile()`, otherwise it will not take effect.

#### Global Logger with Custom Options

If you want to use custom options for the global logger, we provide the `log.UseCustomOptions(opts ...Option)` API to achieve this.
Supported `Option` details can be found in [option.go](./option.go).

### Custom Logger

If you don't want to use the Global Logger, you can initialize a `Logger` instance through the `New` method.
The `New` method accepts `Option` parameters.
Supported `Option` details can be found in [option.go](./option.go).

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

SubLogger support allows you to create an instance that inherits from the parent logger and reset certain `Option` when needed.
For example, in a submodule scenario where a different `LABEL` is used than the parent Logger.

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

### Modifying Time Output Format

RainbowLog's default time format is `2006-01-02 15:04:05.000`,
you can modify this format through the `WithTimeFormat(timeFormat string)` option,
the format can be a string that conforms to golang time format rules,
or it can be `UNIX` or `UNIXMS` or `UNIXMICRO` or `UNIXNANO`,
which represent the return values of `Unix()` or `UnixMilli()` or `UnixMicro()` or `UnixNano()` respectively, used to output `time.Time`.

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

### Advanced Usage

#### BufferedWriter

RainbowLog provides BufferedWriter functionality, which can significantly improve log writing performance:

```go
bufferedWriter := rainbowlog.BufferedWriter(os.Stdout, 4096)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, bufferedWriter),
)
```

#### SyncWriter

For non-thread-safe Writers, you can use SyncWriter wrapper:

```go
syncWriter := rainbowlog.SyncWriter(os.Stdout)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, syncWriter),
)
```

**Note: POSIX and Windows operating systems are inherently write-safe, no need to wrap with SyncWriter!**

#### MultiWriter

Supports writing to multiple targets simultaneously:

```go
multiWriter := rainbowlog.MultiLevelWriter(os.Stdout, fileWriter)

logger := rainbowlog.New(
    rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, multiWriter),
)
```

### Performance Optimization Recommendations

1. **Use Buffered Writer**: For frequent log writing operations, using a buffered writer can significantly improve performance
2. **Set Log Levels Appropriately**: Appropriately increase log levels in production environments to avoid too many debug logs affecting performance
3. **Avoid Recording Too Much Information in Hot Paths**: Try to reduce the content and frequency of log recording on critical performance paths
4. **Use Object Pool**: RainbowLog uses object pools internally to reduce memory allocation, ensure proper use of the `Done()` method to release record objects

### More

Of course, we provide more features and capabilities, looking forward to your exploration and discovery!

## Contact Us

- Email: `ramboll.wong@hotmail.com`
- Telegram Technical Discussion Group: [[Join Now]](https://t.me/+ovN79ozHG4c1YWNl)
- Blog：[Ramboll's Blog](https://ramboll.wang)

## Support with a Donation

If you like this project, feel free to buy the author a cup of lemonade ☕️. Your support is my motivation for continuous updates!

- WeChat Pay:

<img src="https://ramboll.wang/img/wechat_pay.jpg" alt="WeChat Pay" width="200"/>

- Alipay:

<img src="https://ramboll.wang/img/ali_pay.jpg" alt="Alipay" width="200"/>