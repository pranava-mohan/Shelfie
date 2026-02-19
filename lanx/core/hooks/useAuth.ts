import { loginSuccess, logoutSuccess } from "@/core/store/authSlice";
import { useAppDispatch } from "@/core/store/store";
import { getToken, removeToken, saveToken } from "@/core/tokenStorage";
import { jwtDecode } from "jwt-decode";
export const useAuth = () => {
  const dispatch = useAppDispatch();

  async function login(token: string) {
    await saveToken(token);
    const payload: any = jwtDecode(token);
    const adminStatus = (payload.admin as boolean) || false;
    const userID = payload.id as string;
    const avatar = payload.avatar as string;
    dispatch(loginSuccess({ token, adminStatus, userID, avatar }));
  }

  async function logout() {
    await removeToken();
    dispatch(logoutSuccess());
  }

  async function check() {
    const token = await getToken();

    if (token == null) {
      return [false, null];
    } else {
      const payload: any = jwtDecode(token);
      const data = {
        exp: payload.exp as number,
        id: payload.id as string,
        admin: (payload.admin as boolean) || false,
      };

      // {"exp": 234234234, "id": "lsdkjflasdjf", "type": "normal"}
      const expirationDate = new Date(data.exp * 1000);

      if (expirationDate > new Date()) {
        return [true, token];
      } else {
        await removeToken();
        return [false, null];
      }
    }
  }

  return {
    check,
    login,
    logout,
  };
};
