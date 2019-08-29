package main

import (
	"bytes"
	"crypto/sha256"
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

	name := fmt.Sprintf("%.4x-%s", sha256.Sum256(data), runtime.Version())

	if err = ioutil.WriteFile(name+".js", data, 0777); err != nil {
		log.Fatalf("failed to write loader javascript: %v", err)
	}

	boot, err := ioutil.ReadFile("boot.js")
	if err != nil {
		log.Fatalf("failed to read boot javascript: %v", err)
	}

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

	cdn := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wirebase/loader/%s.min.js", name)
	boot = bytes.Replace(boot, []byte(`{{src}}`), []byte(cdn), 1)
	fmt.Println(buf.String())
	fmt.Printf(`<script src="data:text/javascript;base64,` + base64.URLEncoding.EncodeToString([]byte(boot)) + `" />` + "\n")
}
