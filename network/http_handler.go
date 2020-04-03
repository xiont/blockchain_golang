package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	log "github.com/corgi-kx/logcustom"
	"io/ioutil"
	"net/http"
	"strconv"
)

//TODO 用户提交交易主要方法
func httpGenerateTransactions(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body) //读取服务器返回的信息
	if err != nil {
		fmt.Println("read err")
	}
	//fmt.Println( body)
	tss := DeserializeTransactions(body)

	_, _ = w.Write([]byte("这是交易提交返回信息"))

	var s Send
	//var tss []block.Transaction
	////向其他节点包括自己发送交易，会自己处理的
	//fmt.Printf("%s",tss)
	log.Info("接收到%s发送过来的交易", r.Host)
	s.SendTransToPeers(tss)

}

//TODO 用户节点查找UTXOs的方法
func httpFindUTXOFromAddress(w http.ResponseWriter, r *http.Request) {

	addressbyte, err := ioutil.ReadAll(r.Body) //读取服务器返回的信息
	if err != nil {
		fmt.Println("read err")
	}
	//fmt.Println( addressbyte)
	//fmt.Println(string(addressbyte))

	u := block.UTXOHandle{}
	//获取数据库中的未消费的utxo
	utxos := u.FindUTXOFromAddress(string(addressbyte))
	_, _ = w.Write(serializeUTXOs(utxos))
}

//TODO 用户节点查找交易id对应的交易信息
func httpFindTransaction(w http.ResponseWriter, r *http.Request) {
	//交易id
	tsidbyte, err := ioutil.ReadAll(r.Body) //读取服务器返回的信息
	if err != nil {
		fmt.Println("read err")
	}
	bc := block.NewBlockchain()
	ts, _ := bc.FindTransaction(nil, tsidbyte)
	//fmt.Printf("%s",ts)
	_, _ = w.Write(SerializeTransaction(ts))
}

//TODO 处理提交上来的已经证明的区块 var CMineStruct  = "push_mine_struct"  //user_net 向云节点发送已证明的数据
func httpPushMineStruct(w http.ResponseWriter, r *http.Request) {
	//交易id
	mineStructBytes, err := ioutil.ReadAll(r.Body) //读取服务器返回的信息
	if err != nil {
		fmt.Println("read err")
	}
	mine := DeserializeMineStruct(mineStructBytes)
	//TODO 此处应该要做一次验证(要取得刚刚的区块头，用现在的数据做一次验证)
	block.MineReturnStruct.Nonce = mine.Nonce
	block.MineReturnStruct.Ts = mine.Ts
	block.MineReturnStruct.HashByte = mine.HashByte
	block.MineReturnStruct.Err = nil
	block.MineFlag = true
	//fmt.Printf("%s",ts)
	//TODO 接收到信息，应该要向用户节点反馈
	data := jointMessage(cGMessage, []byte("云节点已接收提交区块(但不一定上链)！"))

	_, _ = w.Write(data)
}

//TODO 用户节点获取对应地址的金额
func httpGetBalance(w http.ResponseWriter, r *http.Request) {
	//交易id
	addressbyte, err := ioutil.ReadAll(r.Body) //读取服务器返回的信息
	if err != nil {
		fmt.Println("read err")
	}
	address := string(addressbyte)
	bc := block.NewBlockchain()
	balance := bc.GetBalance(address)

	//int->string->byte
	_, _ = w.Write([]byte(strconv.Itoa(balance)))
}

type MineStruct struct {
	Nonce    int64
	HashByte []byte
	Ts       block.Transaction
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

// 将transaction序列化成[]byte
func SerializeTransaction(ts block.Transaction) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(&ts)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

//交易组的序列化
func SerializeTransactions(tss []block.Transaction) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(&tss)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

//交易组的反序列化
func DeserializeTransactions(d []byte) []block.Transaction {
	var tss []block.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&tss)
	if err != nil {
		log.Panic(err)
	}
	return tss
}

func serializeUTXOs(utxos []*block.UTXO) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(utxos)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func dserializeUTXOs(d []byte) []*block.UTXO {
	var model []*block.UTXO
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&model)
	if err != nil {
		log.Panic(err)
	}
	return model
}

//type myHandler struct{}
//
//func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	_, _ = w.Write([]byte("this is version 3"))
//}
//
//func sayBye(w http.ResponseWriter, r *http.Request) {
//	// 睡眠4秒  上面配置了3秒写超时，所以访问 “/bye“路由会出现没有响应的现象
//	time.Sleep(4 * time.Second)
//	_, _ = w.Write([]byte("bye bye ,this is v3 httpServer"))
//}
