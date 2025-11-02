#include "textflag.h"

// func add(x, y int64) int64
TEXT ·add(SB), NOSPLIT, $0-24
    MOVQ x+0(FP), AX    // 加载第一个参数
    MOVQ y+8(FP), BX    // 加载第二个参数
    ADDQ BX, AX         // 执行加法
    MOVQ AX, ret+16(FP) // 存储结果
    RET
