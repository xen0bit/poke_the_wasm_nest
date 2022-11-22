package main

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	gojs "github.com/tetratelabs/wazero/imports/go"
	"github.com/tetratelabs/wazero/sys"
)

//go:embed wasm_mod/nested.wasm
var addWasm []byte

// main invokes Wasm compiled via `GOARCH=wasm GOOS=js`, which reports the star
// count of wazero.
//
// This shows how to integrate an HTTP client with wasm using gojs.
func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()
	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Combine the above into our baseline config, overriding defaults.
	config := wazero.NewModuleConfig().
		// By default, I/O streams are discarded, so you won't see output.
		WithStdout(os.Stdout).WithStderr(os.Stderr)

	// Compile the WebAssembly module using the default configuration.
	//start := time.Now()
	compiled, err := r.CompileModule(ctx, addWasm)
	if err != nil {
		log.Panicln(err)
	}
	// compilationTime := time.Since(start).Milliseconds()
	// log.Printf("CompileModule took %dms with %dKB cache", compilationTime, dirSize(compilationCacheDir)/1024)

	// Instead of making real HTTP calls, return fake data.
	//ctx = gojs.WithRoundTripper(ctx, &fakeGitHub{})

	// Execute the "run" function, which corresponds to "main" in stars/main.go.
	//start = time.Now()
	err = gojs.Run(ctx, r, compiled, config)
	//runTime := time.Since(start).Milliseconds()
	//log.Printf("gojs.Run took %dms", runTime)
	if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
		log.Panicln(err)
	} else if !ok {
		log.Panicln(err)
	}
}
