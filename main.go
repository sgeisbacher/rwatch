package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pion/webrtc/v4"
)

func main() {
	resetSession()
	ICEServers, err := getICEServersFromServer()
	if err != nil {
		fmt.Printf("E: while getting ICEServers-config: %v\n", err)
		panic("check your internet connection")
	}
	config := webrtc.Configuration{
		ICEServers: ICEServers,
	}

	reset := make(chan bool, 1)
	screen := Screen{}

	go run(&screen, "/bin/bash", []string{"./simple-counter.sh", "5"})

	for {
		// Create a new RTCPeerConnection
		peerConnection, err := webrtc.NewPeerConnection(config)
		if err != nil {
			panic(err)
		}
		defer func() {
			if cErr := peerConnection.Close(); cErr != nil {
				fmt.Printf("cannot close peerConnection: %v\n", cErr)
			}
		}()

		peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
			fmt.Printf("Peer Connection State has changed: %s\n", s.String())

			if s == webrtc.PeerConnectionStateFailed {
				fmt.Println("Peer Connection has gone to failed exiting")
				resetSession()
				reset <- true
			}

			if s == webrtc.PeerConnectionStateClosed {
				fmt.Println("Peer Connection has gone to closed exiting")
				resetSession()
				reset <- true
			}
		})

		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

			d.OnOpen(func() {
				fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					// fmt.Printf("Sending '%s'\n", screen.text)
					err = d.SendText(screen.text)
					if err != nil {
						fmt.Printf("E: while sending: %v\n", err)
						break
					}
				}
			})

			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
			})
		})

		// Waiting for and Import client-offer from signaling-server
		offerData := repeat(getOfferFromServer, 2*time.Second)
		offer := webrtc.SessionDescription{}
		decode(offerData, &offer)
		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			panic(err)
		}

		// Create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		// Create channel that is blocked until ICE Gathering is complete
		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		// Sets the LocalDescription, and start our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}

		fmt.Println("waiting for gathering")
		// Block until ICE Gathering is complete, disabling trickle ICE
		// we do this because we only can exchange one signaling message
		// TODO: in a production application you should exchange ICE Candidates via OnICECandidate
		<-gatherComplete

		// Push answer to signaling-server
		answerSessionDescr := encode(peerConnection.LocalDescription())
		fmt.Println("putting answer ...")
		fmt.Printf("%s\n", answerSessionDescr)
		_, err = http.Post(genUrl("/answer"), "text/plain", strings.NewReader(answerSessionDescr))
		if err != nil {
			fmt.Printf("E: while posting answer to signaling server: %v\n", err)
		}

		<-reset
	}
}

func getOfferFromServer() string {
	resp, err := http.Get(genUrl("/offer"))
	if err != nil {
		fmt.Printf("E: while getting offer from server: %v\n", err)
		return ""
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("E: while reading offer from server: %v\n", err)
		return ""
	}
	return string(respBody)
}
func getICEServersFromServer() ([]webrtc.ICEServer, error) {
	resp, err := http.Get(genUrl("/ice-config"))
	if err != nil {
		return nil, fmt.Errorf("E: while getting ice-config from server: %v\n", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("E: while reading offer from server: %v\n", err)
	}
	var iceServers []webrtc.ICEServer
	err = json.Unmarshal(respBody, &iceServers)
	return iceServers, nil
}

func repeat(fn func() string, delay time.Duration) string {
	for {
		result := fn()
		if len(result) > 0 {
			return result
		}
		fmt.Printf("did not get a result, rerun in %v ...\n", delay)
		time.Sleep(delay)
	}
}

func resetSession() {
	fmt.Println("reset session ...")
	_, err := http.Post(genUrl("/answer"), "text/plain", strings.NewReader(""))
	if err != nil {
		fmt.Printf("E: while resetting answer on signaling server: %v\n", err)
		os.Exit(1)
	}
	_, err = http.Post(genUrl("/offer"), "text/plain", strings.NewReader(""))
	if err != nil {
		fmt.Printf("E: while resetting offer on signaling server: %v\n", err)
		os.Exit(1)
	}
}

// JSON encode + base64 a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode a base64 and unmarshal JSON into a SessionDescription
func decode(in string, obj *webrtc.SessionDescription) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, obj); err != nil {
		panic(err)
	}
}

func genUrl(relPath string) string {
	return fmt.Sprintf("http://165.22.91.102:8080/d43981bd-3822-4127-8cec-662f9a4d54f0%s", relPath)
}
