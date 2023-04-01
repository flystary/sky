package crypto

import "github.com/flystary/sky/g"



func InitEnCryption(passwd string) {
	if passwd != "" {
		g.SECRET_KEY = Md5Raw(passwd)
		g.PROTOCOL_SEPARATOR = string(Md5Raw(passwd + g.PROTOCOL_SEPARATOR)[:4])
		g.PROTOCOL_FEATURE = string(Md5Raw(passwd + g.PROTOCOL_FEATURE)[:8])
	} else {
		// 加密算法导致的缓冲区额外开销
		OVERHEAD = 0
	}
}
