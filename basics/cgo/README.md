# cgo 语句

在 import "C" 语句前的注释中，可以通过 #cgo 语句设置 “编译”、“链接” 阶段的相关参数。

编译阶段参数主要用于：
1. 定义相关宏
2. 指定头文件检索路径

链接阶段参数主要用于：
1. 指定库文件检索路径
2. 指定要链接的库文件

示例：
```go
// #cgo CFLAGS: -DPNG_DEBUG=1 -I ./include
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"
```

CFLAGS部分：
1. -D，定义了宏 PNG_DEBUG，值为 1；
2. -I，定义了头文件包含的检索目录；

LDFLAGS部分：
1. -L，指定了链接时库文件检索目录；
2. -l，指定了链接时需要链接的 png 库；

**NOTE：**
由于 C/C++ 遗留的问题，
1. C “头文件” 检索目录可以是 “相对路径”，
2. 但是 “库文件” 检索目录则需要 “绝对路径”。
3. 在 “库文件” 的检索目录中，可以通过 “${SRCDIR}” 变量表示当前包目录的绝对路径。

示例：
```go
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
上面的代码在链接时将被展开为：
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
```

## 编译选项
#cgo 语句主要影响 CFLAGS、CPPFLAGS、CXXFLAGS、FFLAGS、LDFLAGS 几个编译器环境变量。

LDFLAGS：用于设置 “链接” 时的参数
(在链接阶段，C 和 C++ 的链接选项时通用的)

CFLAGS：用于 C 编译器的选项

CXXFLAGS：用于 C++ 编译器的选项

CPPFLAGS：用于 C 和 C++ 共有的编译选项

FFLAGS:

## 条件选择
#cgo 语句还支持 条件选择，当满足 某个操作系统 或 某个 CPU 架构类型 时，
后面的编译或链接选项生效。

示例：
```go
// #cgo windows CFLAGS: -DX86=1
// #cgo !windows LDFLAGS: -lm
```
在 windows 平台下，编译前会预定义 X86 宏为 1；

在非 windows 平台下，在链接阶段会要求链接 math 数学库。

限制：
    只能是 基于 Go 语言支持的 Windows、Darwin、Linux 等已经支持的操作系统。