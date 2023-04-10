<script>
    import {sendMessage} from "./../resources/messages"
    import { currentConvo } from '../resources/startup';
    function send(){
        if(message==""){
            return;
        }
        
        sendMessage("",message)
        message = "";

    }
    let message = "";
    let inConversation = false;
    currentConvo.subscribe((conver)=>{
        if (conver){
            inConversation = true
        }
    })


    // gonna need something to send messages
</script>

<div class="col-start-2 col-span-5 row-span-full row-start-2">
    
    <div class="grid grid-rows-6 h-full">

        <div class="grid grid-cols-12 grid-rows-2 h-full p-6 row-start-6">
            {#if inConversation}
            <input class="col-span-12 row-span-1 row-start-2" placeholder="Send a message" bind:value={message} on:keydown={(key)=>{if(key.code=="Enter"){send()}}}/>
            {/if}
            <!-- <div class=" bg-blue-700 h-full col-span-2 col-start-12 row-span-1 row-start-2 ml-2 mr-2"></div> -->
        </div>
    </div>
</div>