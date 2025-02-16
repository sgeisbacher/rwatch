package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	WEBRTC_STATE_IDLE = webRTCState{
		Name: "IDLE",
		Msg:  "idle",
	}
	WEBRTC_STATE_CREATING_SESSION = webRTCState{
		Name: "CREATE_SESSION",
		Msg:  "creating session",
	}
	WEBRTC_STATE_AWAITING_CLIENT = webRTCState{
		Name: "AWAITNG_CLIENT",
		Msg:  "awaiting client",
	}
	WEBRTC_STATE_CONNECTING = webRTCState{
		Name: "CONNECTING",
		Msg:  "connecting",
	}
	WEBRTC_STATE_TRANSFER = webRTCState{
		Name: "TRANSFER",
		Msg:  "transfering data",
	}
	WEBRTC_STATE_FAILED = webRTCState{
		Name: "FAILED",
		Msg:  "failed connection",
	}
)

type webRTCState struct {
	Name string
	Msg  string
}

type logEvent struct {
	timestamp time.Time
	msg       string
}

func createAppState(signalingSrv string) *appStateManager {
	return &appStateManager{
		mu:           &sync.Mutex{},
		webRTCStates: []webRTCState{WEBRTC_STATE_IDLE},
		logs:         []logEvent{},
		signalingSrv: signalingSrv,
	}
}

type appStateManager struct {
	mu              *sync.Mutex
	webRTCStates    []webRTCState
	logs            []logEvent
	signalingSrv    string
	webrtcSessionId string
}

func (appState *appStateManager) SetState(newState webRTCState) {
	appState.mu.Lock()
	appState.webRTCStates = append(appState.webRTCStates, newState)
	appState.mu.Unlock()
}

func (appState *appStateManager) Current() webRTCState {
	appState.mu.Lock()
	state := appState.webRTCStates[len(appState.webRTCStates)-1]
	appState.mu.Unlock()
	return state
}

func (appState *appStateManager) SetWebRTCSessionId(id string) {
	appState.Log(fmt.Sprintf("got sessionid %q", id))
	appState.webrtcSessionId = id
}

func (appState *appStateManager) GetWebRTCSessionId() string {
	return appState.webrtcSessionId
}

func (appState *appStateManager) Log(msg string) {
	appState.mu.Lock()
	appState.logs = append(appState.logs, logEvent{time.Now(), msg})
	appState.mu.Unlock()
}

func (appState appStateManager) GenUrl(relPath string) string {
	return fmt.Sprintf("http://%s/%s", appState.signalingSrv, relPath)
}

func (appState appStateManager) GenSessionUrl(relPath string) string {
	return fmt.Sprintf("http://%s/%s%s", appState.signalingSrv, appState.webrtcSessionId, relPath)
}
