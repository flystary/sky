package node

import (
	"fmt"
	"log"
	"net"
	"github.com/flystary/sky/g"
	"github.com/flystary/sky/netIo"
	"github.com/flystary/sky/proto"
	"github.com/flystary/sky/utils"
	"sync"
)

// Node 节点
type Node struct {
	IsAdmin uint16   // Node是否是Admin
	HashID  string   // Node的HashID
	Conn    net.Conn // 与Node的TCP连接

	// Conn的锁，因为Conn读写Packet的时候如果不加锁，多个routine会出现乱序的情况
	ConnReadLock  *sync.Mutex
	ConnWriteLock *sync.Mutex

	// 控制信道缓冲区
	CommandBuffers map[uint16]*Buffer

	// 数据信道缓冲区
	DataBuffers map[uint16]*DataBuffer

	// 是否与本节点直接连接
	DirectConnection bool

	// Socks5Running bool // 防止admin node在一个agent上开启多个连接
}

func NewNode(isAdmin uint16, hashID string, conn net.Conn, directConnection bool) *Node {
	newNode := &Node{
		HashID:           hashID,
		IsAdmin:          isAdmin,
		Conn:             conn,
		ConnReadLock:     &sync.Mutex{},
		ConnWriteLock:    &sync.Mutex{},
		DirectConnection: directConnection,
	}
	newNode.InitDataBuffer()
	return newNode
}

// CommandHandler 协议数据包，将协议数据包分类写入Buffer
func (node *Node) CommandHandler(peerNode *Node) {
	defer peerNode.Disconnect()
	for {
		var lowLevelPacket proto.Packet
		err := peerNode.ReadLowLevelPacket(&lowLevelPacket)
		if err != nil {
			fmt.Println("node disconnect: ", err)
			return
		}
		switch utils.Array32ToUUID(lowLevelPacket.DstHashID) {
		case node.HashID:
			if lowLevelPacket.Separator == g.PROTOCOL_SEPARATOR {
				switch lowLevelPacket.CmdType {
				case proto.SYNC:
					fallthrough
				case proto.LISTEN:
					fallthrough
				case proto.CONNECT:
					fallthrough
				case proto.SHELL:
					fallthrough
				case proto.UPLOAD:
					fallthrough
				case proto.DOWNLOAD:
					fallthrough
				case proto.SOCKS:
					fallthrough
				case proto.LFORWARD:
					fallthrough
				case proto.RFORWARD:
					fallthrough
				case proto.SSHCONNECT:
					node.CommandBuffers[lowLevelPacket.CmdType].WriteLowLevelPacket(lowLevelPacket)
				case proto.SOCKSDATA:
					fallthrough
				case proto.LFORWARDDATA:
					fallthrough
				case proto.RFORWARDDATA:
					var data proto.NetDataPacket
					lowLevelPacket.ResolveData(&data)
					peerNodeID := utils.Array32ToUUID(lowLevelPacket.SrcHashID)
					if Nodes[peerNodeID].DataBuffers[lowLevelPacket.CmdType].GetDataBuffer(data.SessionID) != nil {
						if data.Close == 1 {
							Nodes[peerNodeID].DataBuffers[lowLevelPacket.CmdType].GetDataBuffer(data.SessionID).WriteCloseMessage()
						} else {
							// 只将数据写入数据buffer，不写入整个packet
							Nodes[peerNodeID].DataBuffers[lowLevelPacket.CmdType].GetDataBuffer(data.SessionID).WriteBytes(data.Data)
						}
					}
				default:
					log.Println(fmt.Sprintf("[-]%s", ERR_UNKNOWN_CMD))
				}
			} else {
				log.Println("[-]Separator error")
			}
		default:
			// 如果节点为Agent节点转发
			if node.IsAdmin == 0 {
				nextNode := GNetworkTopology.RouteTable[utils.Array32ToUUID(lowLevelPacket.DstHashID)]
				targetNode := Nodes[nextNode]
				if targetNode != nil {
					targetNode.WriteLowLevelPacket(lowLevelPacket)
				} else {
					log.Println("[-]Can not find target node")
				}
			} else {
				// fmt.Println("src id:", utils.Array32ToUUID(lowLevelPacket.SrcHashID))
				// fmt.Println("dst id:", utils.Array32ToUUID(lowLevelPacket.DstHashID))
				// fmt.Println("dst cmd type:", lowLevelPacket.CmdType)
				fmt.Println("[-]Target node error")
			}
		}
	}
}

func (node *Node) InitCommandBuffer() {
	node.CommandBuffers = make(map[uint16]*Buffer)

	node.CommandBuffers[proto.SYNC] = NewBuffer()
	node.CommandBuffers[proto.LISTEN] = NewBuffer()
	node.CommandBuffers[proto.CONNECT] = NewBuffer()
	node.CommandBuffers[proto.SOCKS] = NewBuffer()
	node.CommandBuffers[proto.UPLOAD] = NewBuffer()
	node.CommandBuffers[proto.DOWNLOAD] = NewBuffer()
	node.CommandBuffers[proto.SHELL] = NewBuffer()
	node.CommandBuffers[proto.LFORWARD] = NewBuffer()
	node.CommandBuffers[proto.RFORWARD] = NewBuffer()
	node.CommandBuffers[proto.SSHCONNECT] = NewBuffer()
}

func (node *Node) InitDataBuffer() {
	node.DataBuffers = make(map[uint16]*DataBuffer)

	node.DataBuffers[proto.SOCKSDATA] = NewDataBuffer()
	node.DataBuffers[proto.LFORWARDDATA] = NewDataBuffer()
	node.DataBuffers[proto.RFORWARDDATA] = NewDataBuffer()
}

// TODO 只有与断掉节点之间相连的节点才会清理路由表/网络拓扑表/节点标号等
// 暂无法做到对全网所有节点的如下信息进行清理，这样有些麻烦，暂时也不是刚需
func (node *Node) Disconnect() {
	node.Conn.Close()
	// 删除网络拓扑
	GNetworkTopology.DeleteNode(node)
	// 删除节点
	delete(Nodes, node.HashID)
	// 删除结构体
	node = nil
}

func (node *Node) ReadLowLevelPacket(packet interface{}) error {
	node.ConnReadLock.Lock()
	defer node.ConnReadLock.Unlock()
	err := netIo.ReadPacket(node.Conn, packet)
	if err != nil {
		return err
	}
	return nil
}

func (node *Node) WriteLowLevelPacket(packet interface{}) error {
	node.ConnWriteLock.Lock()
	defer node.ConnWriteLock.Unlock()
	err := netIo.WritePacket(node.Conn, packet)
	if err != nil {
		return err
	}
	return nil
}

// func (node *Node) ReadPacket(header *protocol.PacketHeader, packet interface{}) error {
// 	node.ConnReadLock.Lock()
// 	defer node.ConnReadLock.Unlock()

// 	// 读数据包的头部字段
// 	err := netIo.ReadPacket(node.Conn, header)
// 	if err != nil {
// 		return err
// 	}
// 	// 读数据包的数据字段
// 	err = netIo.ReadPacket(node.Conn, packet)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (node *Node) WritePacket(header proto.PacketHeader, packet interface{}) error {

	node.ConnWriteLock.Lock()
	defer node.ConnWriteLock.Unlock()

	if g.SECRET_KEY != nil {
		// 加密, 将Packet.Data部分整个加密
		cryptPacket := proto.Packet{}
		cryptPacket.PackHeader(header)
		cryptPacket.PackData(packet)
		err := netIo.WritePacket(node.Conn, cryptPacket)
		if err != nil {
			return err
		}
		return nil
	}

	// 写数据包的头部字段
	header.DataLen, _ = utils.PacketSize(packet)
	err := netIo.WritePacket(node.Conn, header)
	if err != nil {
		return err
	}
	// 写数据包的数据字段
	err = netIo.WritePacket(node.Conn, packet)
	if err != nil {
		return err
	}
	return nil
}

type NodeInfo struct {
	// 节点编号，已被分配的节点编号不会在节点断开后分给新加入网络的节点
	NodeNumber2UUID map[int]string
	NodeUUID2Number map[string]int
	// 节点描述
	NodeDescription map[string]string
}

// NodeExist 节点是否存在
func (info *NodeInfo) NodeExist(nodeID string) bool {
	if _, ok := info.NodeUUID2Number[nodeID]; ok {
		return true
	}
	return false
}

// AddNode 添加一个节点并为节点编号
func (info *NodeInfo) AddNode(nodeID string) {
	number := len(info.NodeNumber2UUID) + 1
	info.NodeNumber2UUID[number] = nodeID
	info.NodeUUID2Number[nodeID] = number
}

// UpdateNoteInfo 根据路由表信息给节点编号
func (info *NodeInfo) UpdateNoteInfo() {
	for key := range GNetworkTopology.RouteTable {
		if !info.NodeExist(key) {
			info.AddNode(key)
		}
	}
}
