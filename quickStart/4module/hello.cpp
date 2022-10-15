// 接口文件 hello.h 是 hello 模块的 实现者 和 使用者 共同的约定，
// 但是该约定并没有要求必须使用 C 语言来实现 SayHello() 函数。

// 也可以用 C++ 语言来重新实现这个函数

//#include <iostream>
//
//extern "C" {
//    #include "hello.c"
//}
//
//void SayHello(const char* s) {
//    std::cout << s;
//}

// Note：在 go build 编译提示 hello.cpp:12:6: error: redefinition of 'void SayHello(const char*)'，
//       注释掉函数即可。