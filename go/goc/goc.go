package goc

/*
#include <stdio.h>

// 方法1: 直接在注释块中定义C函数
void cHello() {
    printf("Hello World from C!\n");
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func goc() {
	// go env -w CGO_ENABLED=1
	// 方法1: 调用在注释块中定义的C函数
	C.cHello()

	// 方法2: 直接调用C标准库函数
	C.printf(C.CString("Hello from C printf directly!\n"))

	// 方法3: 传递参数给C函数
	name := C.CString("Go Developer")
	defer C.free(unsafe.Pointer(name)) // 记得释放C分配的内存
	C.printf(C.CString("Hello, %s!\n"), name)

	fmt.Println("Hello from Go!")
}
