<script>
	import Login from './login.svelte';
	import ConversationBar from './conversationBar.svelte';
	import Messageing from './messageing.svelte';
	import { socket, startUp } from '../resources/startup';

	console.log('sadf');

	//let z = () => {
	//	console.log('sadf');
	//	socket.set(startUp());
	//};
	//    z();
</script>

<!-- {$writableStore} -->
<!-- {#if userId ==""} -->
<div class="grid grid-cols-6 grid-rows-6 h-screen">
	<Login />
	<ConversationBar />
	<Messageing />
</div>

<!-- {:else}
    <p>Welcome {userName}</p>
{/if} -->
<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
<p style="h-10px">Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
<p class="font-bold text-green-500">aasdfasdfa</p>
<button on:click={()=>{count++}} class="bg-green-500 rounded-md"> Clicked {count} Times</button>
<br/>
<Test></Test> -->

