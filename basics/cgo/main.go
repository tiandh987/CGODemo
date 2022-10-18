package main

// 如果在不同的系统下 CGO 对应着不同的 C 代码，那么可以使用 #cgo 语句定义不同的 C 语言的宏，然后通过宏来区分不同的代码，
// 这样就可用用 C 语言中常用的技术来处理不同平台之间的差异代码。

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    static const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    static const char* os = "linux";
#else
#    error(unknown os)
#endif
*/
import "C"

func main() {
	print(C.GoString(C.os))
}
