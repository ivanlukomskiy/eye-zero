package main

import (
	"encoding/json"
	"eye-zero/gst"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
	"os"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	peerConnectionConfig = webrtc.Configuration{}

	videoTrack = &webrtc.TrackLocalStaticSample{}
	pipeline   = &gst.Pipeline{}
	index      string
)

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func main() {
	//videoSrc := flag.String("video-src", "videotestsrc", "GStreamer video src")
	videoSrc := "rpicamsrc bitrate=800000 preview=false ! video/x-h264, width=640, height=480, framerate=30/1 "
	//flag.Parse()

	//config := webrtc.Configuration{
	//	ICEServers: []webrtc.ICEServer{
	//		{
	//			URLs: []string{"stun:stun.l.google.com:19302"},
	//		},
	//	},
	//}

	// Create a video track
	videoTrack_, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion2")
	if err != nil {
		panic(err)
	}
	videoTrack = videoTrack_

	pipeline = gst.CreatePipeline("auto", []*webrtc.TrackLocalStaticSample{videoTrack}, videoSrc)
	pipeline.Start()

	bytes, err := os.ReadFile("index.html") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	index = string(bytes)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)

	log.Fatal(http.ListenAndServe("0.0.0.0:8001", nil))
}

func handleWebsocketMessage(pc *webrtc.PeerConnection, ws *websocket.Conn, message *websocketMessage) error {
	switch message.Event {
	case "offer":
		offer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(message.Data), &offer); err != nil {
			return err
		}

		if err := pc.SetRemoteDescription(offer); err != nil {
			return err
		}

		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			return err
		}

		gatherComplete := webrtc.GatheringCompletePromise(pc)

		if err := pc.SetLocalDescription(answer); err != nil {
			return err
		}

		<-gatherComplete

		answerString, err := json.Marshal(pc.LocalDescription())
		if err != nil {
			return err
		}

		if err = ws.WriteJSON(&websocketMessage{
			Event: "answer",
			Data:  string(answerString),
		}); err != nil {
			return err
		}
	}
	return nil
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		log.Print(err)
		return
	} else if _, err = peerConnection.AddTrack(videoTrack); err != nil {
		log.Print(err)
		return
	}

	defer func() {
		if err := peerConnection.Close(); err != nil {
			log.Println(err)
		}
	}()

	message := &websocketMessage{}
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		} else if err := json.Unmarshal(msg, &message); err != nil {
			log.Print(err)
			return
		}

		if err := handleWebsocketMessage(peerConnection, ws, message); err != nil {
			log.Print(err)
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, index)
}
