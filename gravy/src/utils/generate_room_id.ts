import crpto from "crypto";

export default function generateRoomID(length = 32) {
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);

  let str = "";
  const chars =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_";

  for (let i = 0; i < length; i++) {
    str += chars[array[i] % chars.length];
  }

  return str;
}
