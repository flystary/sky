package node

import (
	"net"
	"sync"

)

type Node struct {
	IsAdmin 		uint16
	HashID			string
	Conn 			net.Conn

	// Conn的锁，因为Conn读写Packet的时候如果不加锁，多个routine会出现乱序的情况
	ConnReadLock	*sync.Mutex
	ConnWriteLock	*sync.Mutex

	// 控制信道缓冲区
	CommandBuffers map[uint16]*Buffer

	// 数据信道缓冲区
	DataBuffers map[uint16]*DataBuffer

	// 是否与本节点直接连接
	DirectConnection bool
}