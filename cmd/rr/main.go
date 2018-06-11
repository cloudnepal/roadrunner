// MIT License
//
// Copyright (c) 2018 SpiralScout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"

	// services (plugins)
	rrhtp "github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/spiral/roadrunner/service/static"

	// cli plugins
	_ "github.com/spiral/roadrunner/cmd/rr/http"
	"github.com/spiral/roadrunner/cmd/rr/debug"

	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"log"
	"runtime/pprof"
	"net/http"
)

var debugMode bool

func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// provides ability to make local connection to services
	rr.Container.Register(rpc.Name, &rpc.Service{})

	// http serving
	rr.Container.Register(rrhtp.Name, &rrhtp.Service{})

	// serving static files
	rr.Container.Register(static.Name, &static.Service{})

	// provides additional verbosity

	// debug mode
	rr.CLI.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "debug mode", )
	cobra.OnInitialize(func() {
		if debugMode {
			service, _ := rr.Container.Get(rrhtp.Name)
			service.(*rrhtp.Service).AddListener(debug.NewListener(rr.Logger).Listener)
		}
	})

	// you can register additional commands using cmd.CLI
	rr.Execute()
}
