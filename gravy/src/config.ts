const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const API = {
  BASE_URL: API_URL,
  LIBRARIAN_LOGIN_URL: "/login/admin",
};

export { API, API_URL };
