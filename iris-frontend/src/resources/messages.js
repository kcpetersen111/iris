import { serverAddr,serverPort } from "../resources/variables";
import { platformList } from '../resources/startup';
import { get } from "svelte/store";
export async function sendMessage(target, message){

    console.log(message)
    
}

export async function pulling(){
    //store the most recent message and only get ones newer than that
    //fetch request for getting platform
    let platform = fetch("http://"+serverAddr+":"+serverPort+"/platform")
    //fetch request for getting messages

    let newPlatform = await platform;
    if(newPlatform.ok){
        let plat = await newPlatform.json()
        console.log(plat)
        if (plat){
            platformList.set(plat.platform);
        }
    }
    // get(platformList).map((item)=>{
    //     if 
    // })
    
}