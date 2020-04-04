package block

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math"
	"math/big"
	"time"
)

//工作量证明(pow)结构体
type proofOfWork struct {
	*BlockHeader
	//难度
	Target *big.Int
}

// TODO 获取POW实例
func NewProofOfWork(blockHeader *BlockHeader) *proofOfWork {
	//_matrix := [10][10]int64{
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//	[10]int64{0,0,0,0,0,0,0,0,0,0},
	//}
	//randomMatrix := RandomMatrix{matrix:_matrix}
	//blockHeader.RandomMatrix = randomMatrix
	target := big.NewInt(1)
	//返回一个大数(1 << 256-TargetBits)
	target.Lsh(target, 256-TargetBits)
	pow := &proofOfWork{blockHeader, target}
	return pow
}

//TODO 当前是否已经出块
var MineFlag = false
var MineReturnStruct struct {
	Nonce    int64
	HashByte []byte
	Ts       Transaction
	Err      error
}

//进行hash运算,获取到当前区块的hash值
func (p *proofOfWork) run(wsend WebsocketSender) (int64, []byte, Transaction, error) {

	//var nonce int64 = 0
	//var hashByte [32]byte
	//var ts Transaction
	wsend.SendBlockHeaderToUser(*p.BlockHeader)

	//注释后可以禁止该节点挖矿
	//go asyncMine(p)

	//TODO 启动一个计时器来检测当前是否已经出块,每秒检测一次
	ticker1 := time.NewTicker(1 * time.Second)
	func(t *time.Ticker) {
		for {
			<-t.C
			if MineFlag == true {
				MineFlag = false
				break
			}
		}
	}(ticker1)

	return MineReturnStruct.Nonce, MineReturnStruct.HashByte[:], MineReturnStruct.Ts, MineReturnStruct.Err
}

//TODO 异步挖矿
func asyncMine(p *proofOfWork) {
	//先给自己添加奖励等
	publicKeyHash := getPublicKeyHashFromAddress(ThisNodeAddr)
	txo := TXOutput{TokenRewardNum, publicKeyHash}
	ts := Transaction{nil, nil, []TXOutput{txo}}
	ts.hash()
	p.BlockHeader.TransactionToUser = ts

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

	for nonce < maxInt {
		//检测网络上其他节点是否已经挖出了区块
		if p.Height <= NewestBlockHeight {
			//结束计数器
			ticker1.Stop()
			MineReturnStruct.Nonce = 0
			MineReturnStruct.HashByte = nil
			MineReturnStruct.Ts = ts
			MineReturnStruct.Err = errors.New("检测到当前节点已接收到最新区块，所以终止此块的挖矿操作")
			MineFlag = true
			return
		}

		//TODO 假设挖出了随机幻方的第一个数字
		randomMatrix := RandomMatrix{[10][10]int64{
			[10]int64{nonce, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}}

		data := p.jointData(randomMatrix)

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

	MineReturnStruct.Nonce = nonce
	MineReturnStruct.HashByte = hashByte[:]
	MineReturnStruct.Ts = ts
	MineReturnStruct.Err = nil
	MineFlag = true
	return
}

//检验区块是否有效
func (p *proofOfWork) Verify() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-TargetBits)
	data := p.jointData(p.BlockHeader.RandomMatrix)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(target) == -1 {
		return true
	}
	return false
}

// TODO 将 上一区块hash、数据、时间戳、难度位数、随机数 拼接成字节数组
func (p *proofOfWork) jointData(randomMatrix RandomMatrix) (data []byte) {
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

	//拼接成交易数组
	//transData := [][]byte{}
	//for _, v := range p.Block.Transactions {
	//	tBytes := v.getTransBytes() //这里为什么要用到自己写的方法，而不是gob序列化，是因为gob同样的数据序列化后的字节数组有可能不一致，无法用于hash验证
	//	transData = append(transData, tBytes)
	//}
	//获取交易数据的根默克尔节点
	//mt := util.NewMerkelTree(transData)

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
