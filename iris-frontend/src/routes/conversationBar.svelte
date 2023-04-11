<script>
	// console.log('made it here');
	import { each, onDestroy } from 'svelte/internal';
	import ConversationComponent from './conversationComponent.svelte';
	import { platformList } from '../resources/startup';
	import { serverAddr,serverPort } from "../resources/variables";
	import { emailStore, jwtStore, userStore } from '../resources/login';
	import { get } from 'svelte/store';
	
	const eyeLogo = new URL('../../static/eyeLogo.svg', import.meta.url).href;
	

	let creating = false;
	let addingConversation = [];
	let conversationName = "";
	// need to get an array of the conversations
	// going to be a pair of names and the uuid
	let conversations = [];
	const unsub = platformList.subscribe(list =>{
		conversations = list;
	})
	onDestroy(unsub);

	function createConversation(){
		creating = !creating;
		if (creating){
			addingConversation.push("")
		} else {
			addingConversation = [];
			conversationName = "";
		}
	}
	function addRecipient(){
		addingConversation.push("")
		addingConversation = addingConversation;
	}
	async function conversationSetup(){
		creating = false;
		let privatConversationName = conversationName;
		let privateConversationRecipients = addingConversation;
		conversationName = "";
		addingConversation = [];
		let body = {
			"platformName": privatConversationName,
			"UserId": privateConversationRecipients,
			"email": get(emailStore)
		}
		let res = await fetch("http://"+serverAddr+":"+serverPort+"/platform",{
				method: "POST",
				credentials:"include",
				headers: {
					"Content-Type":"application/json",
					"Authorization": get(jwtStore)
				},
				body: JSON.stringify(body)
			},
		)
		console.log(res.status)
	}
</script>

<div class="bg-slate-600 row-start-1 row-span-6 text-lg	">
	<span class="flex items-center justify-center">
		<img src={eyeLogo} alt="Logo" class="invert" />
	</span>
	<span class="text-white flex justify-around">
		Conversations
		
			{#if !creating}
		<button on:click={createConversation}>
			+
		</button>
			{:else}
		<button on:click={addRecipient}>
			+
		</button>
		<button on:click={createConversation}>
			-
		</button>
			{/if}
		
	</span>
	{#if creating}
	<span class="flex justify-center">
		<input placeholder="Conversation Name"
			class=" m-1 w-5/6"
			bind:value={conversationName}>
	</span>
		{#each addingConversation as field, index}
			<span class="flex justify-center">
				<input placeholder="Recipient"
					class=" m-1 w-5/6"
					bind:value={addingConversation[index]}>
			</span>
		{/each}
		<span class="flex justify-center h-8">
			<button class="bg-green-500 rounded-lg hover:cursor-pointer p-1 w-16" on:click={conversationSetup}>
				Submit
			</button>
		</span>
	{/if}
	{#each conversations as c}
		<ConversationComponent convo={c} />
	{/each}
</div>
