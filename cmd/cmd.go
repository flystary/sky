package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)


const (
	NORMAL 	= 0
	LISTEN 	= 1
	CONNECT = 2
)

type Option struct {
	SrcPort 	int
	DstIP		string
	DstPort 	int
	Mode		int
	Password	string
}

var Args Option

func init() {
	flag.IntVar(&Args.SrcPort, "s", 0, "Listen a local PORT.")
	flag.StringVar(&Args.DstIP, "addr", "", "Remote ip ADDRESS.")
	flag.IntVar(&Args.DstPort, "d", 0, "The PORT on remote host.")
	flag.StringVar(&Args.Password, "pass", "", "The PASSWORD used in encrypted communication. (optional)")
	// change default Usage
	flag.Usage = usage
}

func usage() {
	ShowBanner()
	fmt.Fprintf(os.Stderr, `Sky version: 1.0
Usage:
	# sky-admin
	# sky-admin -s <src-port>
	# sky-admin -d <dst-port> -addr <dst-ip>
Options:
`)
	flag.PrintDefaults()
}

func ParseArgs() {
	flag.Parse()

	if Args.SrcPort == 0 && Args.DstIP != "" && Args.DstPort != 0 {
		// connect to remote port
		Args.Mode = CONNECT
		return
	}

	if Args.SrcPort != 0 && Args.DstIP == "" && Args.DstPort == 0 {
		// listen a local port
		Args.Mode = LISTEN
		return
	}

	if Args.DstIP == "" && Args.DstPort == 0 && Args.SrcPort == 0 {
		Args.Mode = NORMAL
		return
	}

	// error
	flag.Usage()
	os.Exit(0)
}

func ShowBanner() {
	if runtime.GOOS == "windows" {
		fmt.Println()
	} else {
		fmt.Println()
	}
	fmt.Println()
}

// ShowUsage
// func ShowUsage() {
// 	fmt.Println(`
//   help                                     Help information.
//   exit                                     Exit.
//   show                                     Display network topology.
//   getdes                                   View description of the target node.
//   setdes     [info]                        Add a description to the target node.
//   goto       [id]                          Select id as the target node.
//   listen     [lport]                       Listen on a port on the target node.
//   connect    [rhost] [rport]               Connect to a new node through the target node.
//   sshconnect [user@ip:port] [dport]        Connect to a new node through ssh tunnel.
//   shell                                    Start an interactive shell on the target node.
//   upload     [local_file]  [remote_file]   Upload files to the target node.
//   download   [remote_file]  [local_file]   Download files from the target node.
//   socks      [lport]                       Start a socks5 server.
//   lforward   [lhost] [sport] [dport]       Forward a local sport to a remote dport.
//   rforward   [rhost] [sport] [dport]       Forward a remote sport to a local dport.
// `)
// }