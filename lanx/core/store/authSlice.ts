import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface AuthState {
  token: string | null;
  id: string;
  name?: string;
  email?: string;
  avatar?: string;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const initialState: AuthState = {
  token: null,
  id: "",
  isAuthenticated: false,
  isLoading: true,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    loginSuccess: (
      state,
      action: PayloadAction<{
        token: string;
        adminStatus: boolean;
        userID: string;
        avatar: string;
      }>
    ) => {
      state.token = action.payload.token;
      state.isAuthenticated = true;
      state.isLoading = false;
      state.id = action.payload.userID;
      state.avatar = action.payload.avatar;
    },
    setUserInfo: (
      state,
      action: PayloadAction<{ name: string; email: string; avatar: string }>
    ) => {
      state.name = action.payload.name;
      state.email = action.payload.email;
      state.avatar = action.payload.avatar;
    },
    logoutSuccess: (state) => {
      state.token = null;
      state.isAuthenticated = false;
      state.isLoading = false;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
  },
});

export const { loginSuccess, logoutSuccess, setLoading, setUserInfo } =
  authSlice.actions;
export default authSlice.reducer;
