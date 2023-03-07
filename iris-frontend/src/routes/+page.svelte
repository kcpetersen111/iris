<script>
	import { userStore } from "../resources/login";
	import Login from "./login.svelte";

 


</script>
<!-- {$writableStore} -->
<!-- {#if userId ==""} -->
<Login></Login>
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