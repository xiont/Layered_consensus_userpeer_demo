package network

import (
	"bytes"
	"encoding/gob"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)


//TODO获取金额
func GetBalance(address string) string{
	result,err := HttpChangeData(CGetBalance,[]byte(address))
	if err !=nil{
		panic(err)
	}
	return string(result)
}

//获取附近区块
func GetBlock(offset string) []block.Block{
	offsetNum, err := strconv.Atoi(offset)
	if err !=nil{
		panic(err)
	}
	result,err := HttpChangeData(CGetBlock,util.IntToBytes(offsetNum))
	if err !=nil{
		panic(err)
	}
	blockList := DeserializeBlockList(result)
	return blockList
}




//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func HttpChangeDataByAddr(command string,data []byte,addr string ,port string )([]byte,error){
	//fmt.Printf("http://"+addr+":"+port+"/"+command)
	result,err := HttpPost("http://"+addr+":"+port+"/"+command,data,"application/octet-stream")
	return result,err
}

func HttpChangeData(command string,data []byte)([]byte,error){
	result,err := HttpPost("http://"+RemoteHost+":"+RemotePort+"/"+command,data,"application/octet-stream")
	return result,err
}

func HttpPost(url string, data []byte, contentType string) ([]byte,error) {
	//jsonStr, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	//content = string(result)
	return result,err
}



//交易组的序列化
func SerializeTransactions(tss []block.Transaction ) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(&tss)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}
//交易组的反序列化
func DeserializeTransactions(d []byte) []block.Transaction{
	var tss []block.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&tss)
	if err != nil {
		log.Panic(err)
	}
	return tss
}


//BlockList反序列化
func DeserializeBlockList(d []byte) []block.Block {
	var blockList []block.Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&blockList)
	if err != nil {
		log.Panic(err)
	}
	return blockList
}

// 将BlockList序列化成[]byte
func SerializeBlockList(blockList []block.Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(&blockList)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}
