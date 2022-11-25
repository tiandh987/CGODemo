package main

/*
 #include <string.h>

struct A {
    int age;
    int weight;
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var sA C.struct_A

	sA.age = 10

	fmt.Printf("age: %d, weight: %d\n", sA.age, sA.weight)

	fmt.Printf("sizeof: %d\n", C.sizeof_struct_A)

	fmt.Printf("size_t: %d\n", C.size_t(C.sizeof_struct_A))

	cParam := unsafe.Pointer(C.malloc(C.size_t(C.sizeof_struct_A)))
	//defer C.free(cParam)

	ptrParam := uintptr(cParam)
	cTempParam := (*C.struct_A)(unsafe.Pointer(ptrParam))

	cTempParam.age = 100
	cTempParam.weight = 500
	fmt.Printf("cParam age: %d, weight: %d\n", cTempParam.age, cTempParam.weight)

	C.memset(cParam, 0, C.sizeof_struct_A)

	fmt.Printf("cParam2 age: %d, weight: %d\n", cTempParam.age, cTempParam.weight)

}
