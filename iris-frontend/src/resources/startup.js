import { writable } from "svelte/store"

export const currentConvo = writable("0");


export const socket = writable()
