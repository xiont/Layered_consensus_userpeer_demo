package network

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/logcustom"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

//通过固定格式的地址信息,构建出P2P节点信息对象
func buildPeerInfoByAddr(addrs string) peer.AddrInfo {
	///ip4/0.0.0.0/tcp/9000/p2p/QmUyYpeMSqZp4oNMhANdG6sGeckWiGpBnzfCNvP7Pjgbvg
	p2p := strings.TrimSpace(addrs[strings.Index(addrs, "/p2p")+len("/p2p/"):])
	ipTcp := addrs[:strings.Index(addrs, "/p2p/")]
	//通过ip与端口获得multiAddr
	multiAddr, err := multiaddr.NewMultiaddr(ipTcp)
	if err != nil {
		log.Debug(err)
	}
	//拼接成multiAddr数组
	m := []multiaddr.Multiaddr{multiAddr}
	//获得host.ID
	id, err := peer.IDB58Decode(p2p)
	if err != nil {
		log.Error(err)
	}
	//传入host.ID , multiAddr数组 拼接成P2P节点信息对象
	return peer.AddrInfo{peer.ID(id), m}
}

//默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

//默认前十二位为命令名称
func SplitMessage(message []byte) (cmd string, content []byte) {
	cmdBytes := message[:prefixCMDLength]
	newCMDBytes := make([]byte, 0)
	for _, v := range cmdBytes {
		if v != byte(0) {
			newCMDBytes = append(newCMDBytes, v)
		}
	}
	cmd = string(newCMDBytes)
	content = message[prefixCMDLength:]
	return
}

//gossip 区块头 和 远程主机地址 端口
type AddrMapBlockHeader struct{
	Addr string
	Port string
	BlockHeaderByte []byte
}


//交易组的序列化
func SerializeAddrMapBlockHeader(abh AddrMapBlockHeader ) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(&abh)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}
//交易组的反序列化
func DeserializeAddrMapBlockHeader(d []byte) AddrMapBlockHeader{
	var abh AddrMapBlockHeader
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&abh)
	if err != nil {
		log.Panic(err)
	}
	return abh
}
