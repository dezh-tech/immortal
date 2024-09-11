const eventMsg = `["EVENT",{"kind":1,"id":"cbba15aff4ed4db6370834c9370436ba20615ffa2d170515058f11e522c8dc02","pubkey":"79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798","created_at":1726056813,"tags":[],"content":"test","sig":"282db08a865b5fe97d2351a3e7321ee279b9e6bd16e4e8d9f746c4ea148337e01bae7bc10f7465bdffe740bf6682d7aaadd4777919891ba845e66bf52ee7b6f8"}]`;

let ws = new WebSocket("ws://localhost:3000/ws");

ws.onmessage = (e) => {
  console.log(e.data);
};

ws.send(eventMsg);
