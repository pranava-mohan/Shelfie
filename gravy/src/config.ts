const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";
const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8000/ws/";

const API = {
  BASE_URL: API_URL,
  LIBRARIAN_LOGIN_URL: "/login/admin",
};

export { API, API_URL, WS_URL };
