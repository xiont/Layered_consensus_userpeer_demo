package block

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math"
	"math/big"
	"time"
)

//工作量证明(pow)结构体
type ProofOfWork struct {
	*BlockHeader
	//难度
	Target *big.Int
}

// TODO 获取POW实例
func NewProofOfWork(blockHeader *BlockHeader) *ProofOfWork {
	target := big.NewInt(1)
	//返回一个大数(1 << 256-TargetBits)
	target.Lsh(target, 256-TargetBits)
	pow := &ProofOfWork{blockHeader, target}
	return pow
}

//进行hash运算,获取到当前区块的hash值
func (p *ProofOfWork) run(wsend WebsocketSender) (int64, []byte, Transaction, error) {
	var nonce int64 = 0
	var hashByte [32]byte
	var hashInt big.Int
	log.Info("准备挖矿...")
	//开启一个计数器,每隔五秒打印一下当前挖矿,用来直观展现挖矿情况

	times := 0
	ticker1 := time.NewTicker(5 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			times += 5
			log.Infof("正在挖矿,挖矿区块高度为%d,已经运行%ds,nonce值:%d,当前hash:%x", p.Height, times, nonce, hashByte)
		}
	}(ticker1)

	//先给自己添加奖励等
	publicKeyHash := GetPublicKeyHashFromAddress(ThisNodeAddr)
	txo := TXOutput{TokenRewardNum, publicKeyHash}
	ts := Transaction{nil, nil, []TXOutput{txo}}
	ts.Hash()
	p.BlockHeader.TransactionToUser = ts

	wsend.SendBlockHeaderToUser(*p.BlockHeader)

	for nonce < MaxInt {
		//检测网络上其他节点是否已经挖出了区块
		if p.Height <= NewestBlockHeight {
			//结束计数器
			ticker1.Stop()
			return 0, nil, ts, errors.New("检测到当前节点已接收到最新区块，所以终止此块的挖矿操作")
		}

		//TODO 假设挖出了随机幻方的第一个数字
		randomMatrix := RandomMatrix{[10][10]int64{
			{nonce, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}}

		data := p.JointData(randomMatrix)

		hashByte = sha256.Sum256(data)
		//fmt.Printf("\r current hash : %x", hashByte)
		//将hash值转换为大数字
		hashInt.SetBytes(hashByte[:])
		//如果hash后的data值小于设置的挖矿难度大数字,则代表挖矿成功!
		if hashInt.Cmp(p.Target) == -1 {
			//TODO
			p.RandomMatrix = randomMatrix
			break
		} else {
			//nonce++
			bigInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				log.Panic("随机数错误:", err)
			}
			nonce = bigInt.Int64()
		}
	}
	//结束计数器
	ticker1.Stop()
	log.Infof("本节点已成功挖到区块!!!,高度为:%d,nonce值为:%d,区块hash为: %x", p.Height, nonce, hashByte)
	return nonce, hashByte[:], ts, nil
}

//TODO 当前是否已经出块
var MineFlag = false

//0表示当前没有在挖矿 >0 表示当前挖矿的高度
var NowPowHeight = -1

var MineReturnStruct struct {
	nonce    int64
	hashByte []byte
	ts       Transaction
	err      error
}

type MineStruct struct {
	Nonce    int64
	HashByte []byte
	Ts       Transaction
}

//MineStruct序列化
func SerializeMineStruct(bh *MineStruct) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(bh)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

//MineStruct反序列化
func DeserializeMineStruct(d []byte) *MineStruct {
	var bh MineStruct
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&bh)
	if err != nil {
		log.Panic(err)
	}
	return &bh
}



//TODO 异步挖矿
func AsyncMine(p *ProofOfWork, send Sender, addr string, port string) {
	//未出快标志
	MineFlag = false
	NowPowHeight = p.Height

	//先给自己添加奖励等
	publicKeyHash := GetPublicKeyHashFromAddress(ThisNodeAddr)
	txo := TXOutput{TokenRewardNum, publicKeyHash}
	ts := Transaction{nil, nil, []TXOutput{txo}}
	ts.Hash()
	p.BlockHeader.TransactionToUser = ts

	//TODO 计算 MerkelRootWHash = Hash|W(ts.getTransBytes()+MerkelRootHash)
	//MerkelRootHash + ts.Hash()
	merkelRootWHash := sha256.Sum256( bytes.Join([][]byte{ts.getTransBytes(), p.MerkelRootHash},[]byte("")))
	//WNum 为hash重数
	for i := 0;i < WNum;i++{
		merkelRootWHash = sha256.Sum256(merkelRootWHash[:])
	}
	p.BlockHeader.MerkelRootWHash = merkelRootWHash[:]

	//TODO MerkelRootWHash 签名
	privKey,err := getThisAddrPrivKey()
	if err != nil{
		log.Warn(err)
	}else{
		p.BlockHeader.MerkelRootWSignature = ellipticCurveSign(privKey,merkelRootWHash[:])
	}


	//privKey := wallets.Wallets[string(address)].PrivateKey
	//			//进行签名操作
	//			tss[i].Vint[index].Signature = ellipticCurveSign(privKey, copyTs.TxHash)

	var nonce int64 = 0
	var hashByte [32]byte
	var hashInt big.Int
	log.Info("准备挖矿...")

	//开启一个计数器,每隔五秒打印一下当前挖矿,用来直观展现挖矿情况
	times := 0
	ticker1 := time.NewTicker(5 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			times += 5
			log.Infof("正在挖矿,挖矿区块高度为%d,已经运行%ds,nonce值:%d,当前hash:%x", p.Height, times, nonce, hashByte)
		}
	}(ticker1)

	for nonce < MaxInt {
		//检测网络上其他节点是否已经挖出了区块
		if p.Height <= NewestBlockHeight {
			//结束计数器
			ticker1.Stop()
			MineReturnStruct.nonce = 0
			MineReturnStruct.hashByte = nil
			MineReturnStruct.ts = ts
			MineReturnStruct.err = errors.New("检测到当前节点已接收到最新区块，所以终止此块的挖矿操作")
			log.Info("检测到当前节点已接收到更高的高度信息，所以终止此块的挖矿操作")
			MineFlag = true
			NowPowHeight = 0
			return
		}

		//TODO 假设挖出了随机幻方的第一个数字
		randomMatrix := RandomMatrix{[10][10]int64{
			{nonce, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}}


		data := p.JointData(randomMatrix)

		hashByte = sha256.Sum256(data)
		//fmt.Printf("\r current hash : %x", hashByte)
		//将hash值转换为大数字{
		hashInt.SetBytes(hashByte[:])
		//如果hash后的data值小于设置的挖矿难度大数字,则代表挖矿成功!
		if hashInt.Cmp(p.Target) == -1 {
			//TODO
			p.RandomMatrix = randomMatrix
			break
		} else {
			//nonce++
			bigInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				log.Panic("随机数错误:", err)
			}
			nonce = bigInt.Int64()
		}
	}
	//结束计数器
	ticker1.Stop()
	log.Infof("本用户节点已成功挖到区块!!!,高度为:%d,nonce值为:%d,区块hash为: %x", p.Height, nonce, hashByte)

	p.BlockHeader.Hash = hashByte[:]

	send.SendMinedBlockHeader(*p.BlockHeader, addr, port)
	//此时发送版本信息给其余节点
	send.SendVersionToPeers(p.Height)

	MineReturnStruct.err = nil
	MineFlag = true
	NowPowHeight = 0
	return
}

//检验区块是否有效
func (p *ProofOfWork) Verify() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-TargetBits)
	data := p.JointData(p.BlockHeader.RandomMatrix)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(target) == -1 {
		return true
	}
	return false
}

// TODO 将 上一区块hash、数据、时间戳、难度位数、随机数 拼接成字节数组
func (p *ProofOfWork) JointData(randomMatrix RandomMatrix) (data []byte) {
	preHash := p.BlockHeader.PreHash
	preRandomMatrixByte := RandomMatrixToBytes(p.BlockHeader.PreRandomMatrix)
	merkelRootHash := p.BlockHeader.MerkelRootHash
	merkelRootWHash := p.BlockHeader.MerkelRootWHash
	merkelRootWSignature := p.BlockHeader.MerkelRootWSignature
	cAByte := CAToBytes(p.BlockHeader.CA)

	transactionToUserByte := p.BlockHeader.TransactionToUser.getTransBytes()

	timeStampByte := util.Int64ToBytes(p.BlockHeader.TimeStamp)
	heightByte := util.Int64ToBytes(int64(p.BlockHeader.Height))
	randomMatrixByte := RandomMatrixToBytes(randomMatrix)
	targetBitsByte := util.Int64ToBytes(int64(TargetBits))


	data = bytes.Join([][]byte{
		preHash,
		preRandomMatrixByte,
		merkelRootHash,
		merkelRootWHash,
		merkelRootWSignature,
		cAByte,
		transactionToUserByte,
		timeStampByte,
		heightByte,
		randomMatrixByte,
		targetBitsByte},
		[]byte(""))
	return data
}
