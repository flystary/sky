package protocol

import (
	"sky/crypto"
	"sky/g"
)

// 协议类型
const (
	// 初始化，在node对象建立之前
	INIT = iota
	// 控制协议
	SYNC
	LISTEN
	CONNECT
	SHELL
	UPLOAD
	DOWNLOAD
	SOCKS
	LFORWARD
	RFORWARD
	SSHCONNECT
	// 数据传输协议
	SOCKSDATA
	LFORWARDDATA
	RFORWARDDATA
)

type Packet struct {
	Separator 	string
	CmdType		uint16
	SrcHashID 	[32]byte	//源节点ID
	DstHshID	[32]byte	//目的节点ID
	DataLen		uint64
	Data		[]byte
}

func (packet *Packet) ResolvData(cmdPacket interface{}) {
	KEY := g.SECRET_KEY
	if KEY != nil {
		packet.Data, _ = crypto.Decrypt(packet.Data, KEY)
		packet.DataLen = uint64(len(packet.Data))
	}
}