package main

/*
#include <stdio.h>

void printint(int v) {
    printf("printint: %d\n", v);
}
*/
import "C"

// 在 Go 代码中出现了 import "C" 语句，则表示使用了 CGO 特性，紧临这行语句前面的注释是一种 “特殊语法” ，
// 里面包含的是正常的 C 语言代码。

// 在确保 CGO 启用的情况下，还可以在 “当前目录” 中包含 C/C++ 对应的源文件。

// 头文件被 include 之后，里面的所有 C 语言元素都会被加入 “C” 这个虚拟的包中。

// NOTE：import "C" 语句需要单独占一行，不能与其他包一同 import。

func main() {
	v := 42

	// NOTE:
	//     1. Go 是强类型语言，所以 CGO 中传递的参数类型 “必须” 与声明的类型完全一致，而且传递前 “必须” 用 C 中的转换函数转换成
	//       对应的 C 类型。
	//     2. 通过虚拟的 C 包导入的 C 语言符号并不需要已大写字母开头，它们不受 Go 语言的导出规则约束。
	C.printint(C.int(v))
}
