package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

	data := bytes.Join([][]byte{
		morphdom,
		axios,
		wasmexec,
		load,
	}, nil)

	name := fmt.Sprintf("%s", runtime.Version())
	if err = ioutil.WriteFile(name+".js", data, 0777); err != nil {
		log.Fatalf("failed to write loader javascript: %v", err)
	}

	if err = ioutil.WriteFile("latest.js", data, 0777); err != nil {
		log.Fatalf("failed to write latest loader javascript: %v", err)
	}

	// the boot file does minimal js, small enough to be in a data url executed in the header
	boot, err := ioutil.ReadFile("boot.js")
	if err != nil {
		log.Fatalf("failed to read boot javascript: %v", err)
	}

	cdn := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wirebase/loader/%s.min.js", name)
	boot = bytes.Replace(boot, []byte(`./latest.js`), []byte(cdn), 1)

	// use the terser binary to minify
	exe, err := exec.LookPath("terser")
	if err != nil {
		log.Fatalf("failed to find terser (minifier), make sure it installed")
	}

	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(boot)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to minify boot script: %v", err)
	}

	// output an example script tag that will do it all
	fmt.Println(buf.String())
	fmt.Printf(`<script src="data:text/javascript;base64,` + base64.URLEncoding.EncodeToString(buf.Bytes()) + `" />` + "\n")
}
