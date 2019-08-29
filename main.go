package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

func main() {
	morphdom, err := ioutil.ReadFile("morphdom.js")
	if err != nil {
		log.Fatalf("failed to read morphdom javascript: %v", err)
	}

	axios, err := ioutil.ReadFile("axios.js")
	if err != nil {
		log.Fatalf("failed to read axios javascript: %v", err)
	}

	wasmexec, err := ioutil.ReadFile(filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
	if err != nil {
		log.Fatalf("failed to read the wasm_exec.js from the Go root directory: %w", err)
	}

	load, err := ioutil.ReadFile("load.js")
	if err != nil {
		log.Fatalf("failed to read loader javascript: %v", err)
	}

	if err = ioutil.WriteFile(runtime.Version()+".js", bytes.Join([][]byte{
		morphdom,
		axios,
		wasmexec,
		load,
	}, nil), 0777); err != nil {
		log.Fatalf("failed to write loader javascript: %v", err)
	}

	// commit a release
	// update the bootstrap js
}
