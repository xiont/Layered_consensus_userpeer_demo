package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/network"
	log "github.com/corgi-kx/logcustom"
)

func (cli Cli) transfer(from, to, amount string) {

	var tss []block.Transaction
	blc := block.NewBlockchain()
	tss = blc.CreateTransaction(from, to, amount, network.Send{})
	tssbytes := network.SerializeTransactions(tss)
	_,err := network.HttpChangeData("post_transactions",tssbytes)

	//fmt.Printf("%s",tssbytes)
	if err != nil{
		log.Panic(err)
	}else {
		fmt.Println("已执行转帐命令")
	}
}
