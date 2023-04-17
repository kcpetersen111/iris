<script>
    import {login, createUser} from "../resources/login"
	import Page from "./+page.svelte";
    let Email =  "";
    let Password = "";

let flag = true;
let userSuccessfullyCreated = true;
    
</script>
<div class="flex  h-screen w-screen  items-center justify-center">

    <div class=" items-center w-1/3 justify-center  h-2/3">
        <div class="flex flex-row-reverse font-bold text-4xl items-center pb-12">
            <span class="ml-2 p-0 font-extrabold text-transparent  bg-clip-text bg-gradient-to-t from-red-600 via-yellow-500 to-purple-600">
                Iris
            </span>
            Welcome to 
        </div>

        {#if flag}
        <input placeholder="Email" 
            class=" h-10 w-full mb-10" 
            bind:value={Email} 
            on:keydown={(key)=>{if(key.code=="Enter"){
                flag=false;
                }
            }}>
        <span class="flex justify-between h-12">
            <button on:click={()=>{flag=false}} class="bg-blue-500 rounded-lg hover:cursor-pointer p-1 w-16">
                login
            </button>
            <button on:click={()=>{flag=false}} class="bg-green-500 rounded-lg hover:cursor-pointer p-1 w-16">
                Create
            </button>

        </span>

        {:else}
        <input placeholder="{userSuccessfullyCreated?"Email":"Invalid Email"}" 
            class="  h-10 w-full mb-10 {userSuccessfullyCreated?"":" placeholder-red-600"}" 
           
            bind:value={Email} 
            on:click={()=>{
                flag = true;
            }}>
        <input placeholder="Password" 
            class="h-10 w-full mb-10" 
            type="password"
            bind:value={Password} 
            on:keydown={(key)=>{
                if(key.code=="Enter"){
                    login(Email,Password);
                    Password="";
                    Email=""
                }
            }}>
    <span class="flex justify-between h-12">
        <button 
            on:click={()=>{
                login(Email, Password);
                Password="";
                Email="";
            }} 
            class="bg-blue-500 rounded-lg hover:cursor-pointer p-1 w-16"
        >
            login
        </button>
        <button 
            on:click={()=>{
                createUser(Email, Password).then((success)=>{
                    userSuccessfullyCreated = success;
                    console.log(userSuccessfullyCreated)
                })
                Password="";
                Email="";
            }} 
            class="bg-green-500 rounded-lg hover:cursor-pointer p-1 w-16"
        >
            Create
        </button>
    </span>
        {/if}
    </div>
</div>