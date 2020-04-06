package cli

import (
	"fmt"
	blc "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/network"
	"time"
)

func (cli *Cli) getBlock(offset string) {
	blockList := network.GetBlock(offset)
	for _,block := range blockList{
		fmt.Println("========================================================================================================")
		fmt.Printf("本块hash         %x\n", block.BBlockHeader.Hash)
		fmt.Println("  	------------------------------交易数据------------------------------")
		tss := block.Transactions
		tss = append(tss, block.BBlockHeader.TransactionToUser)
		for _, v := range tss {
			fmt.Printf("   	 本次交易id:  %x\n", v.TxHash)
			fmt.Println("   	  tx_input：")
			for _, vIn := range v.Vint {
				fmt.Printf("			交易id:  %x\n", vIn.TxHash)
				fmt.Printf("			索引:    %d\n", vIn.Index)
				fmt.Printf("			签名信息:    %x\n", vIn.Signature)
				fmt.Printf("			公钥:    %x\n", vIn.PublicKey)
				fmt.Printf("			地址:    %s\n", blc.GetAddressFromPublicKey(vIn.PublicKey))
			}
			fmt.Println("  	  tx_output：")
			for index, vOut := range v.Vout {
				fmt.Printf("			金额:    %d    \n", vOut.Value)
				fmt.Printf("			公钥Hash:    %x    \n", vOut.PublicKeyHash)
				fmt.Printf("			地址:    %s\n", blc.GetAddressFromPublicKeyHash(vOut.PublicKeyHash))
				if len(v.Vout) != 1 && index != len(v.Vout)-1 {
					fmt.Println("			---------------")
				}
			}
		}
		fmt.Println("  	--------------------------------------------------------------------")
		fmt.Printf("随机数           %d\n", block.BBlockHeader.RandomMatrix.Matrix[0][0])
		fmt.Printf("MerkelRoot      %x\n", block.BBlockHeader.MerkelRootHash)
		fmt.Printf("MerkelRootW     %x\n", block.BBlockHeader.MerkelRootWHash)
		fmt.Printf("CA              %s\n", block.BBlockHeader.CA.Address)
		fmt.Printf("时间戳           %s\n", time.Unix(block.BBlockHeader.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("区块高度         %d\n", block.BBlockHeader.Height)
		fmt.Printf("上一个块hash     %x\n", block.BBlockHeader.PreHash)


	}
	//fmt.Printf("地址:%s的余额为：%s\n", address, balance)
}
