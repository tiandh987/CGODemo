要使用 CGO 特性：
1. 需要安装 C/C++ 构建工具链
   1. GCC（macOS/Linux）
   2. MinGW（Windows）
2. 保证环境变量 CGO_ENABLED 被设置为 1

在本地构建时 CGO 默认是启用的，在 “交叉编译” 时 CGO 默认是禁止的。

例如要构建 ARM 环境运行的 Go 程序，需要手工设置好 C/C++ 交叉构建的工具链，
同时开启 CGO_ENABLED 环境变量。然后通过 import "C" 语句启用 CGO 特性。