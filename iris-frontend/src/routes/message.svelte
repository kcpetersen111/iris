
<script>
	import { get } from "svelte/store";
	import { userStore } from "../resources/login";

    export let message;
    export let same;
    let day =[
        "Sunday",
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday"
    ]
    function formatDate(){
        console.log(message["timestamp"])
        let d = new Date(message["timestamp"]+" UTC")
        return `${day[d.getDay()]}, ${d.getMonth()}/${d.getDate()}/${d.getFullYear()} ${d.getHours()}:${d.getMinutes()}`
    }
</script>
<div class="ml-3 mb-1">
    {#if !same}
        <div class="">
            <span class="text-s font-bold">
                {message["name"]}
            </span>
            <span class="text-xs text-gray-500">
                {formatDate()}
            </span>
        </div>
    {/if}

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