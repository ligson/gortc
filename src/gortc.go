package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"gortc/src/config"
	"gortc/src/util"
)

func main() {
	util.LogConfig()

	appConfig, _ := config.GetConfig()
	app := appConfig.App
	fmt.Println("hello")
	logrus.Debug("ok....")
	// 创建媒体引擎和 API 对象
	m := webrtc.MediaEngine{}
	err := m.RegisterDefaultCodecs()
	if err != nil {
		logrus.Error(err)
		return
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	iceServers := []webrtc.ICEServer{
		{
			URLs:           []string{app.ServerType + ":" + app.ServerHost + ":" + app.ServerPort},
			Username:       app.ServerUser,
			Credential:     app.ServerPassword,
			CredentialType: webrtc.ICECredentialTypePassword,
		},
	}

	cfg := webrtc.Configuration{
		ICEServers: iceServers,
	}

	// 创建两个 PeerConnection 对象
	pc1, err := api.NewPeerConnection(cfg)
	if err != nil {
		logrus.Error(err)
		return
	}
	pc1Sdp := pc1.LocalDescription()
	pc1sdpBuffer, err := json.Marshal(pc1Sdp)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debug(pc1sdpBuffer)

	pc1sdpStr := base64.StdEncoding.EncodeToString(pc1sdpBuffer)
	logrus.Debug("SDP信息:" + pc1sdpStr)

	remoteSdp := webrtc.SessionDescription{}
	remoteSdpStr := "xxxx"
	remoteSdpBuffer, err := base64.StdEncoding.DecodeString(remoteSdpStr)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(remoteSdpBuffer, &remoteSdp)
	if err != nil {
		panic(err)
	}
	err = pc1.SetRemoteDescription(remoteSdp)
	if err != nil {
		panic(err)
	}
	answer, err := pc1.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	err = pc1.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// pc1 创建一个 data channel
	dc1, _ := pc1.CreateDataChannel("data", nil)

	// 设置 dc1 的消息处理函数
	dc1.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Println("Message from dc1: " + string(msg.Data))
	})

	// pc1 创建 offer
	offer, _ := pc1.CreateOffer(nil)

	// pc1 设置本地描述
	err = pc1.SetLocalDescription(offer)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 发送消息给 dc2
	dc1.SendText("Hello from dc1")

	// 等待，让数据有时间发送
	select {}
}
