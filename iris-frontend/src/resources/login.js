import { serverAddr,serverPort } from "../resources/variables";
import {writable} from "svelte/store"

const userStore = writable("")
const jwtStore = writable("")
const emailStore = writable("")

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
        

    if(response.status ==201 ){
        console.log("Successful login attempt")
        console.log(body)
        userStore.set(body.userID)
        jwtStore.set(body.name)
        emailStore.set(Email)

    }else if(response.status ==401){
        console.log("Unsuccessful login attempt")
        this.loginPassword ="";
        alert("log in failed")
    } else {
        console.log("There was some sort of error. ", response.status, response);
    }
}

function ValidateEmail(email) 
{
    const emailRegex = new RegExp(/^[A-Za-z0-9_!#$%&'*+\/=?`{|}~^.-]+@[A-Za-z0-9.-]+$/, "gm");
    return emailRegex.test(email)
    
}

export async function createUser(Email, Password){
    if(!ValidateEmail(Email)){
        return false;
    }
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
        // console.log(body);
        // console.log("break")
        
    } catch (error) {
        console.log("response body was not json");
    }
        
    if(response.status ==201 || response.status==200){
        console.log("Successful user creation")
        userStore.set(body.userID)
        jwtStore.set(body.name)
        emailStore.set(Email)
        return true;

    }else if(response.status ==401){
        console.log("Unsuccessful create attempt")
        this.loginPassword ="";
        return false

    } else if(response.status == 409){
        console.log("User email is taken")
        this.loginPassword ="";
        return false
    }else {
        console.log("There was some sort of error. ", response.status, response);
        return false
    }
}

export { userStore, jwtStore, emailStore }