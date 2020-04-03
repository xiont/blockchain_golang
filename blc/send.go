package block

//用于network包向对等节点发送信息
type Sender interface {
	SendVersionToPeers(height int)
	SendTransToPeers(tss []Transaction)
}

//TODO 用于network包向用户节点发送信息
type WebsocketSender interface {
	SendBlockHeaderToUser(bh BlockHeader)
	SendVersionToUser()
}
