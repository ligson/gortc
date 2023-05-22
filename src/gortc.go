package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"gortc/src/config"
	"gortc/src/util"
	"strconv"
)

var pc1 *webrtc.PeerConnection

func chat() {
	// pc1 创建一个 data channel
	dc1, _ := pc1.CreateDataChannel("data", nil)

	// 设置 dc1 的消息处理函数
	dc1.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Println("Message from dc1: " + string(msg.Data))
	})
	// 发送消息给 dc2
	dc1.SendText("Hello from dc1")

}
func iceChange() {
	fmt.Println("input other sdb:")
	remoteSdp := webrtc.SessionDescription{}
	var remoteSdpStr string
	n, err := fmt.Scanln(&remoteSdpStr)
	if err != nil {
		panic(err)
	}
	logrus.Debug("读取" + strconv.Itoa(n) + "个字节,远程sdp-base64是:" + remoteSdpStr)
	remoteSdpBuffer, err := base64.StdEncoding.DecodeString(remoteSdpStr)
	if err != nil {
		panic(err)
	}
	logrus.Debug("远程sdp信息是:" + string(remoteSdpBuffer))
	err = json.Unmarshal(remoteSdpBuffer, &remoteSdp)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	err = pc1.SetRemoteDescription(remoteSdp)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	chat()
}

func onICECandidate(cdt *webrtc.ICECandidate) {
	bytes, _ := json.Marshal(cdt.ToJSON())
	logrus.Debug("on ice candidate:" + string(bytes))
	iceChange()
}

func initConnection() {
	appConfig, _ := config.GetConfig()
	app := appConfig.App
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
	for _, server := range iceServers {
		bytes, _ := server.MarshalJSON()
		logrus.Debug("ICE Server:" + string(bytes))
	}

	cfg := webrtc.Configuration{
		ICEServers: iceServers,
	}

	// 创建两个 PeerConnection 对象
	pc1, err = api.NewPeerConnection(cfg)
	if err != nil {
		logrus.Error(err)
		return
	}
	pc1.OnICECandidate(onICECandidate)

	offer, err := pc1.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	err = pc1.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// 获取 Local Description
	localDesc := pc1.LocalDescription()
	bytes, _ := json.Marshal(localDesc)
	logrus.Info("Local Description json :" + string(bytes))
	sdpBase64 := base64.StdEncoding.EncodeToString(bytes)
	logrus.Info("local sdp base64:" + sdpBase64)

}
func main() {
	util.LogConfig()
	initConnection()
	// 等待，让数据有时间发送
	select {}
}
