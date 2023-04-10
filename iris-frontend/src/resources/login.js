import { serverAddr,serverPort } from "../resources/variables";
import {writable} from "svelte/store"

const userStore = writable("")

export async function login(Email, Password){
    // console.log(Email)
    let msg = {
        "email": Email,
        "password": Password    
    }

    let response  = await fetch("http://"+serverAddr+":"+serverPort+"/user",{
            method:"POST",
            headers: {
                "Content-Type":"application/json",
            },
            body: JSON.stringify(msg)
        },
    )
    let body;
    try {
        body = await response.json();
        console.log(body);
            
    } catch (error) {
        console.log("response body was not json");
    }
        

    if(response.status ==201){
        console.log("Successful login attempt")
        userStore.set(body.userID)
        console.log(userStore)

    }else if(response.status ==401){
        console.log("Unsuccessful login attempt")
        this.loginPassword ="";
        alert("log in failed")
    } else {
        console.log("There was some sort of error. ", response.status, response);
    }
}

export async function createUser(Email, Password){
    let msg = {
        "name": Email,
        "email": Email,
        "password": Password,
        "role": "user"
    }

    let response  = await fetch("http://"+serverAddr+":"+serverPort+"/createUser",{
            method: "POST",
            headers: {
                "Content-Type":"application/json",
            },
            body: JSON.stringify(msg)
        },
    )
    let body;
    try {
        body = await response.json();
        console.log(body);
        console.log("break")
        
    } catch (error) {
        console.log("response body was not json");
    }
        
    if(response.status ==201 || response.status==200){
        console.log("Successful user creation")
        userStore.set(body.userID)


    }else if(response.status ==401){
        console.log("Unsuccessful create attempt")
        this.loginPassword ="";
        return

    } else if(response.status == 409){
        console.log("User email is taken")
        this.loginPassword ="";
        return
    }else {
        console.log("There was some sort of error. ", response.status, response);
        return
    }
}

export { userStore }