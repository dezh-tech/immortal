const req = `["REQ","nak",{"ids":["cbba15aff4ed4db6370834c9370436ba20615ffa2d170515058f11e522c8dc02"]}]`;
const close = `["CLOSE", "nak"]`;

let ws = new WebSocket("ws://localhost:3000/ws");

ws.onmessage = (e) => {
  console.log(e.data);
};

ws.send(
  `["REQ","nak",{"ids":["cbba15aff4ed4db6370834c9370436ba20615ffa2d170515058f11e522c8dc02"]}]`,
);

// ws.send(close);
