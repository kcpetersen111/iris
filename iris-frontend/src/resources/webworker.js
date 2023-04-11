import { get } from 'svelte/store';
import { jwtStore, userStore } from './login';
import { pulling } from './messages';
import { currentConvo, platformList } from './startup';
import {messageStore} from "./messages"

platformList.subscribe(()=>{
    let msg = {
        "platform":get(platformList)
    }
    // let x = JSON.stringify(msg)
    // console.log("Json: "+ x)
    postMessage(msg)
})

messageStore.subscribe(()=>{
    let msg = {
        "message":get(messageStore)
    }
    // let x = JSON.stringify(msg)
    // console.log("Json: "+ x)
    postMessage(msg)
})


onmessage = ({ data: { msg } }) => {
    let arr = msg.split(" ")
    switch (arr[0]) {
        case 'start':
            startTimer();
            break;
        case 'stop':
            stopTimer();
            break
        case 'id':
            userStore.set(arr[1])
            break
        case 'jwt':
            jwtStore.set(arr[1])
            break
        case 'c':
            currentConvo.set(arr[1])
  }
};

let timer = undefined;

const startTimer = () => (timer = setInterval(pulling, 1000));

const stopTimer = () => {
  if (!timer) {
    return;
  }

  clearInterval(timer);
  timer = undefined;
};

export {};