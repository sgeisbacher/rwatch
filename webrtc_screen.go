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
	. "github.com/sgeisbacher/rwatch/utils"
)

const SERVER_TOKEN_API = "262ac525-97e6-4d9c-94b8-78781c192014"

func createWebRTCScreen(appState *appStateManager) *WebRTCScreen {
	return &WebRTCScreen{
		session: &WebRTCSession{
			appState: appState,
		},
	}
}

type WebRTCScreen struct {
	session *WebRTCSession
}

type logFn = func(format string, a ...any)

func createLogger(appState *appStateManager) logFn {
	return func(format string, a ...any) {
		appState.Log(fmt.Sprintf(format, a...))
	}
}
func (screen *WebRTCScreen) InitScreen() {
	screen.session.Init()
}

func (screen *WebRTCScreen) Run(runnerDone chan bool) {}

func (screen *WebRTCScreen) SetOutput(info ExecutionInfo) {
	screen.session.latestExecution = info
}

func (screen *WebRTCScreen) SetError(err error) {
	// TODO
	screen.session.appState.Log(fmt.Sprintf("SetError not yet implemented, got: %v", err))
}

func (screen *WebRTCScreen) Done() {}

type WebRTCSession struct {
	appState        *appStateManager
	latestExecution ExecutionInfo
}

func (sess *WebRTCSession) Init() {
	log := createLogger(sess.appState)

	err := sess.CreateSession()
	if err != nil {
		// TODO handle err
		panic(err)
	}

	ICEServers, err := sess.getICEServersFromServer()
	if err != nil {
		log("E: while getting ICEServers-config: %v", err)
		// TODO avoid panic, handle this
		panic("check your internet connection")
	}
	if len(ICEServers) == 0 {
		// TODO handle err
		log("E: got empty ICEServers-config: %v", ICEServers)
	}
	config := webrtc.Configuration{
		ICEServers: ICEServers,
	}

	reset := make(chan bool, 1)
	for {
		sess.appState.SetState(WEBRTC_STATE_CREATING_SESSION)

		// Create a new RTCPeerConnection
		peerConnection, err := webrtc.NewPeerConnection(config)
		if err != nil {
			panic(err)
		}
		defer func() {
			if cErr := peerConnection.Close(); cErr != nil {
				log("cannot close peerConnection: %v", cErr)
			}
		}()

		peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
			log("Peer Connection State has changed: %s", s.String())

			if s == webrtc.PeerConnectionStateConnecting {
				sess.appState.SetState(WEBRTC_STATE_CONNECTING)
			}
			if s == webrtc.PeerConnectionStateConnected {
				sess.appState.SetState(WEBRTC_STATE_TRANSFER)
			}
			if s == webrtc.PeerConnectionStateClosed {
				log("Peer Connection has gone to closed exiting")
				sess.appState.SetState(WEBRTC_STATE_FAILED)
				time.Sleep(5 * time.Second)
				sess.Reset(log)
				reset <- true
			}
		})

		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			log("New DataChannel %s %d", d.Label(), d.ID())

			d.OnOpen(func() {
				log("Data channel '%s'-'%d' open.", d.Label(), d.ID())

				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					// //log fmt.Printf("Sending '%s'\n", screen.text)
					data, err := json.Marshal(sess.latestExecution)
					if err != nil {
						log("E: while marshalling data for transport: %v", err)
					}
					err = d.SendText(string(data))
					if err != nil {
						log("E: while sending: %v", err)
						break
					}
				}
			})

			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				log("Message from DataChannel '%s': '%s'", d.Label(), string(msg.Data))
			})
		})

		// Waiting for and Import client-offer from signaling-server
		sess.appState.SetState(WEBRTC_STATE_AWAITING_CLIENT)
		offerData := repeat(func() string { return sess.getOfferFromServer(log) }, 2*time.Second)
		offer := webrtc.SessionDescription{}
		decode(offerData, &offer)
		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			// TODO avoid panic, handle
			fmt.Printf("got setremotedescr-err %v\n", err)
			panic(err)
		}

		// Create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			fmt.Printf("got createanswer-err %v\n", err)
			// TODO avoid panic, handle
			panic(err)
		}

		// Create channel that is blocked until ICE Gathering is complete
		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		// Sets the LocalDescription, and start our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			fmt.Printf("got setlocaldescr-err %v\n", err)
			// TODO avoid panic, handle
			panic(err)
		}

		//log fmt.Println("waiting for gathering")
		// Block until ICE Gathering is complete, disabling trickle ICE
		// we do this because we only can exchange one signaling message
		// TODO: in a production application you should exchange ICE Candidates via OnICECandidate
		fmt.Println("waiting for ICE gathering ...")
		<-gatherComplete
		fmt.Println("done ICE gathering")

		// Push answer to signaling-server
		answerSessionDescr := encode(peerConnection.LocalDescription())
		log("sending answer ...")
		//log fmt.Printf("%s\n", answerSessionDescr)
		_, err = http.Post(sess.appState.GenSessionUrl("/answer"), "text/plain", strings.NewReader(answerSessionDescr))
		if err != nil {
			log("E: while posting answer to signaling server: %v", err)
		}

		<-reset
	}
}

func (sess *WebRTCSession) CreateSession() error {
	url := sess.appState.GenUrl("/create-session")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", SERVER_TOKEN_API)
	if err != nil {
		return fmt.Errorf("error while creating req %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("E: while getting offer from server: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("E: while reading offer from server: %w", err)
	}
	serverSessionId := string(respBody)
	sess.appState.SetWebRTCSessionId(serverSessionId)
	return nil
}

func (sess WebRTCSession) getOfferFromServer(log logFn) string {
	resp, err := http.Get(sess.appState.GenSessionUrl("/offer"))
	if err != nil {
		log("E: while getting offer from server: %v", err)
		return ""
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log("E: while reading offer from server: %v", err)
		return ""
	}
	return string(respBody)
}
func (sess WebRTCSession) getICEServersFromServer() ([]webrtc.ICEServer, error) {
	resp, err := http.Get(sess.appState.GenSessionUrl("/ice-config"))
	if err != nil {
		return nil, fmt.Errorf("E: while getting ice-config from server: %v", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("E: while reading offer from server: %v", err)
	}
	var iceServers []webrtc.ICEServer
	err = json.Unmarshal(respBody, &iceServers)
	return iceServers, nil
}

func repeat(fn func() string, delay time.Duration) string {
	// start := time.Now()
	for {
		result := fn()
		// repeatingSince := time.Since(start)
		// fmt.Printf("debug: got offer from server %s (%0.0fs): %s\n", time.Now().Format("15:04:05"), repeatingSince.Seconds(), result)
		if len(result) > 0 {
			return result
		}
		//log fmt.Printf("did not get a result, rerun in %v ...\n", delay)
		time.Sleep(delay)
	}
}

func (sess WebRTCSession) Reset(log logFn) {
	log("reset session ...")
	_, err := http.Post(sess.appState.GenSessionUrl("/answer"), "text/plain", strings.NewReader(""))
	if err != nil {
		log("E: while resetting answer on signaling server: %v", err)
		// TODO avoid panic, handle
		os.Exit(1)
	}
	_, err = http.Post(sess.appState.GenSessionUrl("/offer"), "text/plain", strings.NewReader(""))
	if err != nil {
		log("E: while resetting offer on signaling server: %v", err)
		// TODO avoid panic, handle
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
