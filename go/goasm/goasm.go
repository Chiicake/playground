package goasm

// use asm in go to add two int64 numbers
// noescape to tell the compiler use the asm function instead of the go function
// the assembly function add is defined in assembly.s

//go:noescape
func add(x, y int64) int64
