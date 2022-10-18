//go:build debug
// +build debug

// build 标志 是在 Go 或 CGO 环境下的 C/C++ 文件开头的一种 “特殊的注释”。

package main

var buildMode = "debug"

// 可以通过以下命令构建：
//      go build -tags="debug"
//      go build -tags="windows debug"         (可以通过 -tags 命令行参数指定多个 build 标志，它们之间用 “空格” 分隔)

//
//
// 逻辑操作（或（空格）、与（逗号））
// 示例：
//      // +build linux,386 darwin,!cgo
// 表示：
//      只有在 “linux/386” 或 “Darwin平台下非CGO环境” 才进行构建
