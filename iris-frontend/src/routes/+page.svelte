<script>
	import Login from './login.svelte';
	import ConversationBar from './conversationBar.svelte';
	import Messageing from './messageing.svelte';
	import { onDestroy } from 'svelte';
	
	import { jwtStore, userStore } from '../resources/login';
	import { get } from 'svelte/store';
	import { currentConvo, platformList } from '../resources/startup';
	import { messageStore } from '../resources/messages';
	
	let syncWorker = undefined;
	let userId = "";
	let first = true;
	userStore.subscribe((id)=>{
		userId = id;
		if(id && first){
			first = false;
			startPulling()
		}
	})
	const loadWorker = async () => {
		const SyncWorker = await import('../resources/webworker?worker');
		
		syncWorker = new SyncWorker.default();

		syncWorker.onmessage = ( msg ) => {
			// if(msg === undefined){
			// 	console.log("Got message undefined from worker")
			// 	return
			// }
			try {
				// console.log( msg.data["platform"])
				// let pulled = JSON.parse(msg)
				if(msg.data["platform"]){
					platformList.set(msg.data["platform"])
				}
				if(msg.data["message"]){
					messageStore.set(msg.data["message"])
				}
			}catch{
			}
			
		}
		currentConvo.subscribe((c)=>{
			syncWorker.postMessage({msg:"c "+get(currentConvo)})
		})

		syncWorker.postMessage({ msg: "id "+get(userStore)})
		syncWorker.postMessage({ msg: "jwt "+get(jwtStore)})
		syncWorker.postMessage({ msg: 'start' });
	};
	

async function startPulling(){
	loadWorker()
}

onDestroy(() => syncWorker?.postMessage({ msg: 'stop' }));



</script>

{#if userId !=""}
	<!-- Home page -->
		<div class="grid grid-cols-6 grid-rows-6 h-screen">
			<ConversationBar />
			<Messageing />
		</div>
{:else}
    <!-- login page -->

		<Login />
{/if} 

