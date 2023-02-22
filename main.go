package main

import (
	"runtime"
	"sky/cmd"
	"sky/crypto"
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