// shim the websocket library
globalThis.WebSocket = require("websocket").w3cwebsocket;
const { connect, StringCodec } = require("nats.ws");

(async () => {
  try {
    // to create a connection to a nats-server:
    const nc = await connect({ servers: "ws://localhost:4223" });
    const sc = StringCodec();
    console.log("connected to nats, waiting for messages...");
    const sub = nc.subscribe("webhook");
    (async () => {
      for await (const m of sub) {
        console.log(`[${sub.getProcessed()}]: ${sc.decode(m.data)}`);
      }
      console.log("subscription closed");
    })();
  } catch (err) {
    console.log("cannot connect to nats-server. Is it running?");
  }
})();
