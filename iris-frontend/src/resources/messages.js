import { serverAddr,serverPort } from "../resources/variables";
import { platformList } from '../resources/startup';
import {get, writable} from "svelte/store"
import { jwtStore, userStore } from '../resources/login';
import { currentConvo } from '../resources/startup';

export const messageStore = writable([])

export async function sendMessage(message){

    let body = {
        "message":message,
        "sender": get(userStore),
        "platform": get(currentConvo)
    }
    let msg = await fetch ("http://"+serverAddr+":"+serverPort+"/message",{
        method:"POST",
        credentials:"include",
        headers:{
            "Content-Type":"application/json",
            "Authorization": get(jwtStore)
        },
        body:JSON.stringify(body)
    })
    console.log(message)
    console.log(msg)
    
}

export async function pulling(){
    //store the most recent message and only get ones newer than that
    //fetch request for getting platform
 
    let body = {
        "platformID":get(userStore)
    }
    let platform = fetch("http://"+serverAddr+":"+serverPort+"/getPlatform", {
        method:"POST",
        credentials:"include",
        headers: {
            "Content-Type":"application/json",
            "Authorization": get(jwtStore)
        },
        body:JSON.stringify(body)
    }
        
    )
    //fetch request for getting messages
    
    let msgBody = {
        "platformID":get(currentConvo)
    }
    let msg = fetch ("http://"+serverAddr+":"+serverPort+"/getMessages",{
        method:"POST",
        credentials:"include",
        headers:{
            "Content-Type":"application/json",
            "Authorization": get(jwtStore)
        },
        body:JSON.stringify(msgBody)
    })
   

    //handle resolutions
    let newPlatform = await platform;
    if(newPlatform.ok){
        let plat = await newPlatform.json()
        
        if (plat && get(platformList) != plat){
            platformList.set(plat);
        }
    }
 

    let convoMsg = await msg;
    if(convoMsg.ok){
        let m = await convoMsg.json()
        if (m && get(messageStore) != m){
            // platformList.set(plat);
            messageStore.set(m)
        }
    }
    
}