const API = {
  BASE_URL: process.env.EXPO_PUBLIC_API_URL,
  LIBRARIAN_LOGIN_URL: "/login/admin",
};

const AUTH = {
  GOOGLE_URL: API.BASE_URL + "/login/google",
};

export { API, AUTH };
