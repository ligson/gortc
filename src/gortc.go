package main

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"gortc/src/util"
)

func main() {
	util.LogConfig()
	fmt.Println("hello")
	logrus.Debug("ok....")
	// 创建媒体引擎和 API 对象
	m := webrtc.MediaEngine{}
	m.RegisterDefaultCodecs()
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	iceServers := []webrtc.ICEServer{
		{
			URLs:           []string{"turn:mycoturn:53478"},
			Username:       "user",
			Credential:     "password",
			CredentialType: webrtc.ICECredentialTypePassword,
		},
	}

	cfg := webrtc.Configuration{
		ICEServers: iceServers,
	}

	// 创建两个 PeerConnection 对象
	pc1, _ := api.NewPeerConnection(cfg)
	pc2, _ := api.NewPeerConnection(cfg)

	// pc1 创建一个 data channel
	dc1, _ := pc1.CreateDataChannel("data", nil)

	// 设置 dc1 的消息处理函数
	dc1.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Println("Message from dc1: " + string(msg.Data))
	})

	// pc2 设置 data channel 处理函数
	pc2.OnDataChannel(func(dc2 *webrtc.DataChannel) {
		dc2.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Println("Message from dc2: " + string(msg.Data))

			// 将消息回发给 dc1
			dc2.SendText("Hello from dc2")
		})
	})

	// pc1 创建 offer
	offer, _ := pc1.CreateOffer(nil)

	// pc1 设置本地描述
	pc1.SetLocalDescription(offer)

	// pc2 设置远端描述
	pc2.SetRemoteDescription(offer)

	// pc2 创建 answer
	answer, _ := pc2.CreateAnswer(nil)

	// pc2 设置本地描述
	pc2.SetLocalDescription(answer)

	// pc1 设置远端描述
	pc1.SetRemoteDescription(answer)

	// 发送消息给 dc2
	dc1.SendText("Hello from dc1")

	// 等待，让数据有时间发送
	select {}
}
