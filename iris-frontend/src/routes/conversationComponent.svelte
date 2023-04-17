<script lang="ts">
	import { get } from 'svelte/store';
	import { userStore } from '../resources/login';

	// will need name and a conversation to be pulled up when this is clicked on
import { currentConvo } from '../resources/startup';
import {serverAddr, serverPort } from '../resources/variables'
import { messageStore } from '../resources/messages';
	import { onMount } from 'svelte';


export let convo: any;
let audioStream;
let audioURL;
let sourceBuffer;
let speakerReady = false;
let source:MediaSource;
let queue = [];
onMount(()=>{
	source = new MediaSource();
	
	source.addEventListener('sourceopen',(ev:Event)=>{
		console.log("source open")
		let b:MediaSource = ev.target;
		// console.log(b.canPlayType("audio/webm; codecs=opus"))
		sourceBuffer = b.addSourceBuffer("audio/webm;codecs=opus")
		sourceBuffer.addEventListener('update',()=>{
			console.log("source buffer update")
			if (queue.length > 0 && !sourceBuffer.updating) {
				sourceBuffer.appendBuffer(queue.shift());
			}
		})
		speakerReady = true;
	},false)
	// audioURL = URL.createObjectURL(source);
	// const au = new Audio()
	// au.autoplay = true
	// au.src = audioURL
	// au.play()
	// au.src = audioURL
// 	audioURL = new MediaSource();
// 	sourceBuffer = audioURL.addSourceBuffer("audio/webm; codecs='opus'")
})
let inCall = false;
function changeConversation() {
	messageStore.set([])
	//switch what is in focus for the user
	currentConvo.set(convo.platformID);
	console.log(convo.platformName + ' was swaped to');
}

function shorten(text: String) {
	if (text == undefined) {
		return "";
	}
	if (text.length > 14) {
		return text.slice(0, 15);
	}
	return text;
}

function until(conditionFunction) {

	const poll = resolve => {
		if(conditionFunction()) resolve();
		else setTimeout(_ => poll(resolve), 400);
	}
	return new Promise(poll);
}

async function setUpwebSocks(reconnectTries:number,) {
	let connected = false;
	let metaData = {
		"Caller": get(userStore),
		"Callee": convo.platformID
	}
	const socket = new WebSocket('ws://'+serverAddr+":"+serverPort+'/start')
	socket.addEventListener('open', ()=>{
		socket.send(JSON.stringify(metaData))
		connected = true;
		console.log("WebSocket is open")
	})
	socket.addEventListener('connect',()=>{
		console.log("WebSocket is connected")
	})
	socket.addEventListener('error', (err)=>{
		console.log(err)
		if (reconnectTries < 5){
			reconnectTries++;
			setTimeout(function() {
				setUpwebSocks(reconnectTries);
			}, 1000);
		}
	})
	// socket.addEventListener('message',(msg)=>{console.log(msg);playBlob(msg)})
	await until(_=> connected == true)
	return socket
}

async function endCall() {
	console.log("end call")
	inCall = false
}
async function awaitEndCall(ws) {
	await until((_=> inCall == false))
	console.log("closing after end call")
	ws.close()
	source.endOfStream();
}

async function startCall() {
	console.log("starting a call")
	let localStreamRequest = getAudioStream();
	
	// audioURL = new MediaSource();
	// sourceBuffer = audioURL.addSourceBuffer("audio/webm; codecs='opus'")

	let ws;
	if (window['WebSocket']){
		console.log("browser supports ws")
		ws = await setUpwebSocks(0);
	} else {
		console.log("browser does not support websocket protocols")
		return;
	}

	inCall = true;
	
	
	let audioStream = await localStreamRequest;

	
	

	const mr = new MediaRecorder(audioStream,{mimeType:"audio/webm; codecs= opus"});
	mr.start(3)
	mr.ondataavailable = e =>{
		if(!speakerReady){
			return;
		}
		
		// console.log("made it to appendBuffer")
		// e.data.arrayBuffer().then(b=>{
		// 	sourceBuffer.appendBuffer(b)
		// })
		 console.log(e)
		//  if(ws.readyState == ws.open	){

			 ws.send(e.data)
		// } else {
		// }
	}

	ws.addEventListener('message',(msg)=>{
		console.log("got a message")
		if (sourceBuffer.updating || !speakerReady || queue.length > 0){
			queue.push(msg)
		} else{
			
			msg.data.arrayBuffer().then(b=>{
				console.log(sourceBuffer)
				sourceBuffer.appendBuffer(b)
			})
		}
	})
	ws.addEventListener('close',(ev:CloseEvent)=>{
		console.log("MR stop")
		console.log(ev)
		mr.stop()
		audioStream.getAudioTracks()[0].stop()
		// audioURL.endOfStream()
	})
	awaitEndCall(ws);

	playBlob()
}

// au.controls = true
async function playBlob(){
	const au = new Audio()
	au.onclose = ()=>{
		console.log("closing audio")
	}
	au.autoplay = true
	au.src = URL.createObjectURL(source); 
	// au.volume = 0.50;
	au.play()
}

function playStream(s){
	const au = new Audio()
	au.autoplay = true
	au.srcObject = s
}

async function getAudioStream() {
	return await navigator.mediaDevices.getUserMedia({ video: false, audio: true });
}

</script>
<div>
	<!-- <audio>
		<source src="{audioURL}" type="audio/mp3">
	</audio> -->
	
	<div
		class="text-white hover:cursor-pointer pl-5 hover:bg-slate-400 hover:text-black flex justify-between"
		on:click={changeConversation}
	>
		<span> {shorten(convo.platformName)} </span>
		{#if !inCall}
		<button class="pr-5" on:click={startCall}>Call</button>
		{:else}
		<button class="pr-5" on:click={endCall}>End</button>
		{/if}
	
	</div>
</div>
