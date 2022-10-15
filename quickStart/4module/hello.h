// 抽象一个 hello 模块，模块的全部接口函数都在 hello.h 头文件中

// 作为 hello 模块的用户，可以放心的使用 SayHello() 函数，而无需关心函数的具体实现。
void SayHello(const char* s);

