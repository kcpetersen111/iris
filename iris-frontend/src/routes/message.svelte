
<script>
	import { get } from "svelte/store";
	import { userStore } from "../resources/login";

    export let message;
</script>
<div class="ml-3">
    {#if get(userStore) == message["sender"]}
        <div class=" bg-blue-800 rounded-md w-fit py-1 px-2 text-gray-200">
           {message["message"]}
        </div>
    {:else}
    <div class=" bg-blue-500 rounded-md w-fit py-1 px-2 text-gray-200">
        {message["message"]}
     </div>
    {/if}
</div>