<script>
    import {sendMessage} from "./../resources/messages"
    import { currentConvo } from '../resources/startup';
    import {get, writable} from "svelte/store"
    import {messageStore} from "../resources/messages"
	import Message from "./message.svelte";
	import { afterUpdate, onMount } from "svelte";
   

    function send(){
        if(message==""){
            return;
        }
        sendMessage(message)
        message = "";
    }

    let message = "";
    let inConversation = false;
    currentConvo.subscribe((conver)=>{
        if (conver){
            inConversation = true
        }
    })

    let scrollBind;
    let first = true;
    const scrollToBottom = async (node) => {
        first = false;
        node.scroll({ top: node.scrollHeight, behavior: 'smooth' });
    }; 

    afterUpdate(()=>{scrollToBottom(scrollBind)})

    let displayMessages = [];
    messageStore.subscribe(()=>{
            displayMessages = get(messageStore)
            console.log(displayMessages)
            // scrollToBottom(scrollBind)
    })

    // gonna need something to send messages
</script>

<div class="col-start-2 col-span-5 row-span-full row-start-1  max-h-full">
    <div class=" bg-white  grid grid-rows-6 max-h-screen grid-row h-full">
        <div bind:this={scrollBind} class="row-start-1 row-span-5 overflow-y-scroll">
            <div class=" pt-1">
                {#each displayMessages as m }
            
                    <Message message={m}/>
                    
                {/each}
                <!-- {#if first}
                {scrollToBottom(scrollBind)}
                {/if} -->
            </div>
        </div>

        <div class="grid grid-cols-12 grid-rows-2 h-full p-6 row-start-6">
            {#if inConversation}
            <input class="col-span-12 row-span-1 row-start-2" placeholder="Send a message" bind:value={message} on:keydown={(key)=>{if(key.code=="Enter"){send()}}}/>
            {/if}
            <!-- <div class=" bg-blue-700 h-full col-span-2 col-start-12 row-span-1 row-start-2 ml-2 mr-2"></div> -->
        </div>
    </div>
</div>