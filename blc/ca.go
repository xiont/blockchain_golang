package block

type CACertificate struct {
	//FIXME 生成该区块的计算节点的钱包地址
	Address string
}

func CAToBytes(cA CACertificate) []byte {
	return []byte(cA.Address)
}
