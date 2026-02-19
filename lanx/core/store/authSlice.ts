import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface AuthState {
  token: string | null;
  id: string;
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

export const { loginSuccess, logoutSuccess, setLoading } = authSlice.actions;
export default authSlice.reducer;
