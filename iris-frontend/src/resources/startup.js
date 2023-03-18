import { writable } from "svelte/store"
import { serverAddr, serverPort } from "../resources/variables"

export const currentConvo = writable("0");

export function startUp() {

  // const soc = new WebSocket("ws://localhost:4444/start")
  //    const soc = new WebSocket("ws://localhost:4444/start")
  console.log("asdf")
  //    return soc
}
export const socket = writable()
