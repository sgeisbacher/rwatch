export interface LogEvent {
	level: 'debug' | 'info' | 'warn' | 'error';
	msg: string;
	timestamp: Date;
}

export async function initConnection(
	setData: (data: string[]) => void,
	addLogEvent: (event: LogEvent) => void
) {
	const resp = await fetch(genUrl('./ice-config'));
	const iceServers = await resp.json();
	const pc = new RTCPeerConnection({
		iceServers
	});
	const log = (msg: string) => {
		addLogEvent({ level: 'info', msg, timestamp: new Date() });
		// document.getElementById('logs').innerHTML += msg + '<br>';
	};
	const writeToTerm = (msg: string) => {
		setData([msg]);
		// TODO
		// let termElem = document.getElementById('terminal');
		// if (!termElem) {
		// 	termElem = document.createElement('pre');
		// 	termElem.setAttribute('id', 'terminal');
		// 	const container = document.getElementById('termContainer');
		// 	container.appendChild(termElem);
		// }
		// termElem.innerHTML = msg;
	};

	const sendChannel = pc.createDataChannel('foo');
	sendChannel.onclose = () => console.log('sendChannel has closed');
	sendChannel.onopen = () => {
		console.log('sendChannel has opened');
		clearTimeout(interv);
	};
	sendChannel.onmessage = (e) => writeToTerm(e.data);

	let interv: ReturnType<typeof setInterval>;
	pc.oniceconnectionstatechange = () => log(pc.iceConnectionState);
	pc.onicecandidate = async (event) => {
		if (event.candidate === null) {
			const offerSD = btoa(JSON.stringify(pc.localDescription));
			log('sending webrtc-offer ...');
			await fetch(genUrl('./offer'), {
				method: 'POST',
				body: offerSD
			});
			interv = setInterval(startSession, 3000);
		}
	};

	pc.onnegotiationneeded = () =>
		pc
			.createOffer()
			.then((d) => pc.setLocalDescription(d))
			.catch(log);

	// SEND MESSAGE BACK TO SERVER
	// window.sendMessage = () => {
	// 	const message = document.getElementById('message').value;
	// 	if (message === '') {
	// 		return alert('Message must not be empty');
	// 	}

	// 	sendChannel.send(message);
	// };

	const startSession = async () => {
		log('contacting signaling server ...');
		const resp = await fetch(genUrl('./answer'));
		const answerSD = await resp.text();
		// log("answer:", answerSD);

		if (answerSD === '') {
			return log('no webrtc-answer yet, retrying ...');
		}

		try {
			pc.setRemoteDescription(JSON.parse(atob(answerSD)));
		} catch (e) {
			alert(e);
		}
	};

	let time = 0;
	function logSeconds() {
		time++;
		// TODO
		// document.getElementById("timetrack").innerHTML += time + ";";
	}
	setInterval(logSeconds, 1000);
}

function genUrl(relPath: string) {
	return relPath;
}
