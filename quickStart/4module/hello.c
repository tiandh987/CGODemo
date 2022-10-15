#include "hello.h"
#include <stdio.h>

// 作为 SayHello() 函数的实现者，函数的实现只需要满足头文件中函数的声明的规范即可。
void SayHello(const char* s) {
    puts(s);
}