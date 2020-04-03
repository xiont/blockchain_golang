package block

import (
	"bytes"
	"encoding/gob"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math/big"
	"time"
)

type BlockHeader struct {
	//0 上一个区块的hash
	PreHash []byte
	//上一个区块的幻方
	PreRandomMatrix RandomMatrix
	//1 FIXME 当前Block的Merkel根的W重计算值,首先设置为 []byte
	MerkelRootHash  []byte
	MerkelRootWHash []byte
	//2 FIXME W重计算值的数字签名（由挖矿人签名）
	MerkelRootWSignature []byte
	//3 FIXME CA证书
	CA CACertificate
	//4 FIXME 对用户节点的激励（激励金额由计算节点决定）
	TransactionToUser Transaction
	//时间戳
	TimeStamp int64
	//区块高度
	Height int
	//FIXME 随机幻方
	RandomMatrix RandomMatrix
	//Nonce int64
	//本区块hash = Hash（PreHash + MerkelRootHash + MerkelRootSignature + CA + Transaction(计算节点会对激励验证) +
	// TimeStamp + Height + NBits + RandomMatrix）< 难度
	Hash []byte
}

type Block struct {
	//区块头
	BBlockHeader BlockHeader
	//数据data 当云计算节点收集交易时，
	//第1笔交易是输入为空，输入为给打包该区块的云计算节点的激励
	Transactions []Transaction
}

////进行挖矿来生成区块
//func mineBlock(transaction []Transaction, preHash []byte, height int) (*Block, error) {
//	timeStamp := time.Now().Unix()
//	//hash数据+时间戳+上一个区块hash
//	block := Block{preHash, transaction, timeStamp, height, 0, nil}
//	pow := NewProofOfWork(&block)
//	nonce, hash, err := pow.run()
//	if err != nil {
//		return nil, err
//	}
//	block.Nonce = nonce
//	block.Hash = hash[:]
//	log.Info("pow verify : ", pow.Verify())
//	log.Infof("已生成新的区块,区块高度为%d", block.Height)
//	return &block, nil
//}

//TODO 进行挖矿来生成区块
//blockheader 中preHash 已知 merkelTreeRoot 已知 preRandomMatrix
//CA 已知 Height 已知
//未知：MerkelRootWHash MerkelRootWSignature TransactionToUser TimeStamp RandomMatrix Hash
func mineBlock(transactions []Transaction, preHash []byte, height int, preRandomMatrix RandomMatrix, cA CACertificate, wsend WebsocketSender) (*Block, error) {

	//生成一笔奖励
	timeStamp := time.Now().Unix()
	//生成交易的merkelRoot
	merkelRootHash := generateMerkelRoot(transactions)

	//txo := TXOutput{nil, nil}
	ts := Transaction{nil, nil, nil}
	randomMatrix_ := RandomMatrix{Matrix: [10][10]int64{}}
	//transactionToUser := ts

	blockHeader := BlockHeader{
		PreHash:              preHash,
		PreRandomMatrix:      preRandomMatrix,
		MerkelRootHash:       merkelRootHash,
		MerkelRootWHash:      []byte(""),
		MerkelRootWSignature: []byte(""),
		CA:                   cA,
		TransactionToUser:    ts,
		TimeStamp:            timeStamp,
		Height:               height,
		RandomMatrix:         randomMatrix_,
		Hash:                 nil,
	}

	block := Block{
		BBlockHeader: blockHeader,
		Transactions: transactions,
	}

	//TODO POW
	//传递Blockheader进行PoW
	pow := NewProofOfWork(&blockHeader)
	nonce, hash, new_ts, err := pow.run(wsend)
	if err != nil {
		return nil, err
	}
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

	blockHeader.RandomMatrix = randomMatrix
	blockHeader.TransactionToUser = new_ts
	blockHeader.Hash = hash[:]

	block.BBlockHeader.RandomMatrix = randomMatrix
	block.BBlockHeader.TransactionToUser = new_ts
	block.BBlockHeader.Hash = hash[:]

	log.Info("pow verify : ", pow.Verify())
	log.Infof("已生成新的区块,区块高度为%d", block.BBlockHeader.Height)
	return &block, nil
}

// TODO 生成交易组的merkel根
func generateMerkelRoot(transactions []Transaction) []byte {
	txs := [][]byte{}
	for _, transaction := range transactions {
		txs = append(txs, transaction.getTransBytes())
	}
	merkelTree := util.NewMerkelTree(txs)
	return merkelTree.MerkelRootNode.Data
}

////生成创世区块
//func newGenesisBlock(transaction []Transaction) *Block {
//	//创世区块的上一个块hash默认设置成下面的样子
//	preHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	//生成创世区块
//	genesisBlock, err := mineBlock(transaction, preHash, 1)
//	if err != nil {
//		log.Error(err)
//	}
//	return genesisBlock
//}

//TODO 生成创世区块
func newGenesisBlock(transaction []Transaction, wsend WebsocketSender) *Block {
	//创世区块的上一个块hash默认设置成下面的样子
	preHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	//前一个区块的随机幻方
	randomMatrix := RandomMatrix{[10][10]int64{
		[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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
	//cA证书
	cA := CACertificate{Address: ThisNodeAddr}
	//生成创世区块
	genesisBlock, err := mineBlock(transaction, preHash, 1, randomMatrix, cA, wsend)
	if err != nil {
		log.Error(err)
	}

	return genesisBlock
}

//// 将Block对象序列化成[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *Block) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}

//blockHeader序列化
func SerializeBlockHeader(bh *BlockHeader) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(bh)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

//blockHeader反序列化
func DeserializeBlockHeader(d []byte) *BlockHeader {
	var bh BlockHeader
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&bh)
	if err != nil {
		log.Panic(err)
	}
	return &bh
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//
func isGenesisBlock(block *Block) bool {
	var hashInt big.Int
	//hashInt.SetBytes(block.PreHash)
	hashInt.SetBytes(block.BBlockHeader.PreHash)
	if big.NewInt(0).Cmp(&hashInt) == 0 {
		return true
	}
	return false
}
