package block

import (
	"bytes"
	"github.com/corgi-kx/blockchain_golang/util"
)

type RandomMatrix struct {
	Matrix [10][10]int64
}

//刚开始用 [0][0]号元素代替nonce

func RandomMatrixToBytes(randomMatrix RandomMatrix) []byte {
	results := make([][]byte, 10)
	for _, array := range randomMatrix.Matrix {
		for _, num := range array {
			results = append(results, util.Int64ToBytes(num))
		}
	}
	sep := []byte("")

	return bytes.Join(results, sep)
}
