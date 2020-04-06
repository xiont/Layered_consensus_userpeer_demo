package block

//用于network包向对等节点发送信息
type Sender interface {
	//用于network包向对等节点发送信息
	SendVersionToPeers(height int)
	SendTransToPeers(tss []Transaction)

	//用于用户节点向云节点请求数据
	GetUTXOsBytes(address string) []byte
	//根据交易ID获取交易
	GetTrans(ID []byte) []byte
	//发送挖到的区块头
	SendMinedBlockHeader(minedBH BlockHeader, addr string, port string)
}

//TODO 用于network包向用户节点发送信息
type WebsocketSender interface {
	SendBlockHeaderToUser(bh BlockHeader)
	SendVersionToUser()
}
