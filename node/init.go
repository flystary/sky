package node

import (
	"errors"
	"log"
	"net"
	"github.com/flystary/sky/g"
	"github.com/flystary/sky/netIo"
	"github.com/flystary/sky/proto"
	"github.com/flystary/sky/utils"
)

var ERR_UNKNOWN_CMD = errors.New("Unknown command type")

func ServerInitConnection(conn net.Conn) (bool, *Node) {
	// 端口重用模式下发送的一段垃圾数据
	netIo.Write(conn, []byte(g.PROTOCOL_FEATURE))

	var PacketHeader proto.PacketHeader
	netIo.ReadPacket(conn, &PacketHeader)

	if PacketHeader.Separator != g.PROTOCOL_SEPARATOR ||
		PacketHeader.CmdType != proto.INIT {
		log.Println("[-]InitPacket error: separator or cmd type")
		conn.Close()
		return false, nil
	}

	var initPacketCmd proto.InitPacketCmd
	netIo.ReadPacket(conn, &initPacketCmd)

	initPacketRet := proto.InitPacketRet{
		OsType:  utils.GetSystemType(),
		HashID:  utils.UUIDToArray32(CurrentNode.HashID),
		IsAdmin: 0,
	}
	size, _ := utils.PacketSize(initPacketRet)
	PacketHeader = proto.PacketHeader{
		Separator: g.PROTOCOL_SEPARATOR,
		CmdType:   proto.INIT,
		DataLen:   size,
	}
	netIo.WritePacket(conn, PacketHeader)
	netIo.WritePacket(conn, initPacketRet)

	// clientNode := &Node{
	// 	HashID:        utils.Array32ToUUID(initPacketCmd.HashID),
	// 	IsAdmin:       initPacketCmd.IsAdmin,
	// 	Conn:          conn,
	// 	ConnReadLock:  &sync.Mutex{},
	// 	ConnWriteLock: &sync.Mutex{},
	// 	// Socks5SessionIDLock:  &sync.Mutex{},
	// 	// Socks5DataBufferLock: &sync.RWMutex{},
	// 	DirectConnection: true,
	// }
	// clientNode.InitDataBuffer()

	clientNode := NewNode(
		initPacketCmd.IsAdmin,
		utils.Array32ToUUID(initPacketCmd.HashID),
		conn,
		true,
	)

	Nodes[utils.Array32ToUUID(initPacketCmd.HashID)] = clientNode
	clientNodeID := utils.Array32ToUUID(initPacketCmd.HashID)
	GNetworkTopology.AddRoute(clientNodeID, clientNodeID)
	GNetworkTopology.AddNetworkMap(CurrentNode.HashID, clientNodeID)
	GNodeInfo.AddNode(clientNodeID)

	return true, clientNode
}

func ClentInitConnection(conn net.Conn) (bool, *Node) {
	// 端口重用模式下发送的一段垃圾数据
	netIo.Write(conn, []byte(g.PROTOCOL_FEATURE))

	// Node的初始状态为UNINIT，所以首先CurrentNode会向连接的对端发送init packet
	initPacketCmd := proto.InitPacketCmd{
		OsType:  utils.GetSystemType(),
		HashID:  utils.UUIDToArray32(CurrentNode.HashID),
		IsAdmin: 0,
	}
	size, _ := utils.PacketSize(initPacketCmd)
	PacketHeader := proto.PacketHeader{
		Separator: g.PROTOCOL_SEPARATOR,
		CmdType:   proto.INIT,
		DataLen:   size,
	}
	netIo.WritePacket(conn, PacketHeader)
	netIo.WritePacket(conn, initPacketCmd)

	// 读取返回包
	// init阶段可以看做连接建立阶段，双方进行握手后交换信息
	// 所有init阶段无需校验数据包中的DstHashID,因为此时双方还没有获取双方的HashID
	netIo.ReadPacket(conn, &PacketHeader)
	if PacketHeader.Separator != g.PROTOCOL_SEPARATOR ||
		PacketHeader.CmdType != proto.INIT {
		log.Println("[-]InitPacket error: separator or cmd type error")
		conn.Close()
		return false, nil
	}
	var initPacketRet proto.InitPacketRet
	netIo.ReadPacket(conn, &initPacketRet)
	// 新建节点加入map
	// serverNode := &Node{
	// 	HashID:        utils.Array32ToUUID(initPacketRet.HashID),
	// 	IsAdmin:       initPacketRet.IsAdmin,
	// 	Conn:          conn,
	// 	ConnReadLock:  &sync.Mutex{},
	// 	ConnWriteLock: &sync.Mutex{},
	// 	// Socks5SessionIDLock:  &sync.Mutex{},
	// 	// Socks5DataBufferLock: &sync.RWMutex{},
	// 	DirectConnection: true,
	// }
	// serverNode.InitDataBuffer()

	serverNode := NewNode(
		initPacketRet.IsAdmin,
		utils.Array32ToUUID(initPacketRet.HashID),
		conn,
		true,
	)

	Nodes[utils.Array32ToUUID(initPacketRet.HashID)] = serverNode

	serverNodeID := utils.Array32ToUUID(initPacketRet.HashID)
	GNetworkTopology.AddRoute(serverNodeID, serverNodeID)
	GNetworkTopology.AddNetworkMap(CurrentNode.HashID, serverNodeID)
	GNodeInfo.AddNode(serverNodeID)

	return true, serverNode
}
