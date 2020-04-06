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
	cVersion     command = "version"  //p2p and usernet
	cGetHash     command = "getHash"  //p2p
	cHashMap     command = "hashMap"	//p2p
	cGetBlock    command = "getBlock"	//p2p
	cBlock       command = "block"		//p2p
	cTransaction command = "transaction"	//p2p
	cMyError     command = "myError"	//p2p

	cBHeader	 command = "blockHeader" //云节点向用户节点推送未证明的区块头
	cABHeader    command = "aBlockHeader" //用户节点向用户节点推送区块头，在区块头基础上附带了远程主机地址和端口，用于提交
	cGMessage    command = "generalMsg"  //向用户节点发送通用信息

)

var	CUTXOs   = "find_utxo_from_address"
var CFindTs   = "find_transaction"
var CMinedBH  = "push_mined_blockheader"  //user_net 向云节点发送已证明的数据
var CGetBalance = "get_balance" //user_net向云节点获取地址的金额
var CGetBlock = "get_block" //user_net向云节点获取区块数据

//默认远程地址
var RemoteHost = "127.0.0.1"

//默认远程访问端口（http请求）
var RemotePort ="7004"

