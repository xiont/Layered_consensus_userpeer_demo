package network

import (
	"flag"
	log "github.com/corgi-kx/logcustom"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"time"
)


//接收远程主机的推送入口


func MonitorBlockHeader(clier Clier) {
	//监听云计算节点的websocket区块头推送服务
	var addr = flag.String("addr", RemoteHost+":"+RemotePort, "http service address")
	flag.Parse()
	//log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Info("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				//log.Println("read:", err)
				return
			}
			//取信息的前十二位得到命令
			cmd, content := SplitMessage(message)
			log.Tracef("本节点已接收到websocket推送命令：%s", cmd)
			switch command(cmd) {

			case cBHeader:
				log.Tracef("本节点已接收到websocket推送命令：%s", cmd)
				//直接拆解，后面还要组合
				go transferBlockHeader(content)
			case cVersion:
				log.Tracef("本节点已接收到websocket推送命令：%s", cmd)
				//发送完全信息，transfer之后进行拆解
				go transferBlockVersion(message)
			case cGMessage:
				log.Tracef("本节点已接收到websocket推送命令：%s", cmd)
				log.Info(string(content))

				//log.Printf("recv: %s", message)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				//log.Println("write:", err)
				return
			}
		case <-interrupt:
			//log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				//log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}


func transferBlockHeader(content []byte){
	//TODO 此处需要将content 通过p2p发送给其他节点，同时要附带该区块云计算节点的位置 即RemoteHost和RemotePort
	//组装远程主机和区块头
	abh := AddrMapBlockHeader{
		Addr:            RemoteHost,
		Port:            RemotePort,
		BlockHeaderByte: content,
	}

	abhByte := SerializeAddrMapBlockHeader(abh)

	data := jointMessage(cABHeader,abhByte)
	_ = gossip.Publish(pubsubTopic, data)

}




func transferBlockVersion(content []byte){
	_ = gossip.Publish(pubsubTopic, content)
}