<script lang="ts">
	// will need name and a conversation to be pulled up when this is clicked on
	import { currentConvo } from '../resources/startup';
	export let convo: any;
	function changeConversation() {
		//fetch the conversation
		// change the conversation I will just set it in a store and then either do a force reload or have the state change take care of it
		//switch what is in focus for the user
		currentConvo.set(convo.uuid);
		console.log(convo.name + ' was swaped to');
	}
	function shorten(text: String) {
		if (text.length > 14) {
			return text.slice(0, 15);
		}
		return text;
	}
</script>

<div
	class="text-white hover:cursor-pointer pl-5 hover:bg-slate-400 hover:text-black"
	on:click={changeConversation}
>
	<span> {shorten(convo.name)} </span>
</div>
