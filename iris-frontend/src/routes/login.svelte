<script>
    import {login, createUser} from "../resources/login"
    let Email =  "";
    let Password = "";


    
</script>
<div class="flex bg-gray-600 h-10 justify-end items-center">

    <div class="items-center">
        <input placeholder="Email" bind:value={Email} class="h-5" on:keydown={(key)=>{if(key.code=="Enter"){login()}}}>
        <input placeholder="Password" bind:value={Password} class="h-5" on:keydown={(key)=>{if(key.code=="Enter"){login()}}}>
        <span on:click={()=>{login(Email, Password);Password="";Email=""}} class="bg-blue-500 rounded-lg hover:cursor-pointer p-1">
            login
        </span>
        <span on:click={()=>{createUser(Email, Password);Password="";Email=""}} class="bg-green-500 rounded-lg hover:cursor-pointer p-1">
            Create
        </span>
    </div>
</div>