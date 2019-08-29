package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
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

	name := runtime.Version()
	if err = ioutil.WriteFile(name+".js", bytes.Join([][]byte{
		morphdom,
		axios,
		wasmexec,
		load,
	}, nil), 0777); err != nil {
		log.Fatalf("failed to write loader javascript: %v", err)
	}

	cdn := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wirebase/loader/%s.min.js", name)
	bootstrap := fmt.Sprintf(`document.documentElement.className+=" wb-js","object"==typeof WebAssembly&&(document.documentElement.className+=" wb-wasm",window.addEventListener("DOMContentLoaded",function(e){el=document.createElement("script"),el.src="%s",document.body.appendChild(el)}));`, cdn)
	fmt.Printf(`<script src="data:text/javascript;base64,` + base64.URLEncoding.EncodeToString([]byte(bootstrap)) + `" />` + "\n")
}
