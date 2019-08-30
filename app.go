// +build wasm

package main

import "syscall/js"

func main() {
	js.Global().Get("document").Get("documentElement").Get("classList").Call("add", `launched`)
}
