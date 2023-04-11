<script lang="ts">
	import { get } from 'svelte/store';
	import { userStore } from '../resources/login';

	// will need name and a conversation to be pulled up when this is clicked on
import { currentConvo } from '../resources/startup';
import {serverAddr, serverPort } from '../resources/variables'


export let convo: any;
let audioStream;
function changeConversation() {
	//fetch the conversation
	// change the conversation I will just set it in a store and then either do a force reload or have the state change take care of it
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
	socket.addEventListener('message',(msg)=>{console.log(msg)})
	await until(_=> connected == true)
	return socket
}

async function startCall() {
	console.log("starting a call")
	let localStreamRequest = getAudioStream();
	

	let ws;
	if (window['WebSocket']){
		console.log("browser supports ws")
		ws = await setUpwebSocks(0)
	} else {
		console.log("browser does not support websocket protocols")
		return;
	}

	let audioStream = await localStreamRequest;
	const mr = new MediaRecorder(audioStream);
	mr.start(3)
	mr.ondataavailable = (e) =>{
		ws.send(e)
	}
	ws.addEventListener('close',()=>{
		mr.stop()
	})
	
	// ws.addEventListener('open',(event)=>{
	// 	ws.send()
	// })
	// playStream(audioStream)

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

<audio>
	<source src="{audioStream}" type="audio/mp3">
</audio>

<div
	class="text-white hover:cursor-pointer pl-5 hover:bg-slate-400 hover:text-black flex justify-between"
	on:click={changeConversation}
>
	<span> {shorten(convo.platformName)} </span>
	<span class="pr-5" on:click={startCall}>Call</span>
</div>
