package main

import (
	"runtime"
	"github.com/flystary/sky/cmd"
	"github.com/flystary/sky/crypto"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// args
	cmd.ParseArgs()
	// logo
	cmd.ShowBanner()

	//
	crypto.InitEnCryption(cmd.Args.Password)

}