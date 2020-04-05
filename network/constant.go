package network

import "github.com/libp2p/go-libp2p-core/host"

//p2p相关,程序启动时,会被配置文件所替换
var (
	RendezvousString = "meetme"
	ProtocolID       = "/chain/1.1.0"
	ListenHost       = "0.0.0.0"
	ListenPort       = "3001"
	localHost        host.Host
	localAddr        string

	PubsubTopic   = "/libp2p/cloud_blockchain/1.0.0"
	BootstrapAddr = ""
	BootstrapHost = ""
	BootstrapPort = ""
)

//websocket推送监听端口(默认7001)
var (
	WebsocketAddr = "0.0.0.0"
	WebsocketPort = "7001"
)

//交易池
var tradePool = Transactions{}

//交易池默认大小
var TradePoolLength = 2

//版本信息 默认0
const versionInfo = byte(0x00)

//发送数据的头部多少位为命令
const prefixCMDLength = 12

type command string

//网络通讯互相发送的命令
const (
	cVersion     command = "version"     //p2p and usernet
	cGetHash     command = "getHash"     //p2p
	cHashMap     command = "hashMap"     //p2p
	cGetBlock    command = "getBlock"    //p2p
	cBlock       command = "block"       //p2p
	cTransaction command = "transaction" //p2p
	cMyError     command = "myError"     //p2p

	cBHeader  command = "blockHeader" //user_net 向用户节点推送未证明的区块头
	cGMessage command = "generalMsg"  //向用户节点发送通用信息
)

//websocket推送信息格式
