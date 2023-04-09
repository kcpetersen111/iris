<script>
	// console.log('made it here');
	import { each, onDestroy } from 'svelte/internal';
	import ConversationComponent from './conversationComponent.svelte';
	import { platformList } from '../resources/startup';
	const eyeLogo = new URL('../../static/eyeLogo.svg', import.meta.url).href;
	//	import Image from './Image.svelte';
	//    import {conversationComponent} from ConversationComponent
	// need to get an array of the conversations
	// going to be a pair of names and the uuid
	let conversations = [];
	const unsub = platformList.subscribe(list =>{
		conversations = list;
	})
	onDestroy(unsub);
</script>

<div class="bg-slate-600 row-start-1 row-span-6">
	<span class="flex items-center justify-center">
		<img src={eyeLogo} alt="Logo" class="invert" />
	</span>
	{#each conversations as c}
		<ConversationComponent convo={c} />
	{/each}
</div>
