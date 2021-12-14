import { Client } from "./generated";

let ws = new WebSocket("ws://localhost:8080");

let client = new Client(async (msg: string): Promise<string> => {
  ws.send(msg);
  let res = await new Promise<string>((resolve) => {
    ws.onmessage = (ev) => {
      resolve(ev.data);
    }
  })
  return res;
})

async function main() {
  let res = await client.Receive({
    Val1: 1,
    Val2: ["Hi"],
  });
  console.log("received", res);
}

ws.onopen = main;

