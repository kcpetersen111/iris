<script lang="ts">
	// will need name and a conversation to be pulled up when this is clicked on

	function changeConversation() {
		//fetch the conversation

		//switch what is in focus for the user
		console.log('testing');
	}
	function shorten(text: String) {
		if (text.length > 14) {
			return text.slice(0, 15);
		}
		return text;
	}
</script>

<div class="text-white hover:cursor-pointer pl-5" on:click={changeConversation}>
	<span> {shorten('people in the conversation')} </span>
</div>
