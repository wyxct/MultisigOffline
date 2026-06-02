package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("        冷钱包离线签名工具        ")
	fmt.Println("   【私钥永不触网 · 断网使用】")
	fmt.Println("=====================================")

	var unsignedTxHex string
	var privateKeyHex string
	var chainID int64

	fmt.Print("请输入 未签名交易Hex: ")
	fmt.Scanln(&unsignedTxHex)
	fmt.Print("请输入 私钥: ")
	fmt.Scanln(&privateKeyHex)
	fmt.Print("请输入 ChainID: ")
	fmt.Scanln(&chainID)

	txBytes, err := hex.DecodeString(strings.TrimPrefix(unsignedTxHex, "0x"))
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}
	var tx types.Transaction
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		fmt.Println("解析交易失败:", err)
		return
	}

	signedTx, err := offlineSign(privateKeyHex, &tx, big.NewInt(chainID))
	if err != nil {
		fmt.Println("签名失败:", err)
		return
	}

	signedBin, _ := signedTx.MarshalBinary()
	signedHex := "0x" + hex.EncodeToString(signedBin)

	fmt.Println("\n=====================================")
	fmt.Println("签名成功！复制以下内容到热端：")
	fmt.Println("=====================================")
	fmt.Println(signedHex)
	fmt.Println("=====================================")
}

func offlineSign(privateKeyHex string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
