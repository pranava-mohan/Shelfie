import { API } from "@/config";
import axios from "axios";
import { useAppSelector } from "../store/store";

export const useClient = () => {
  const token = useAppSelector((state) => state.auth.token);

  const authHeaders = token
    ? {
        Authorization: `Bearer ${token}`,
      }
    : {};

  const client = axios.create({
    baseURL: API.BASE_URL,
    headers: {
      ...authHeaders,
    },
  });

  return { client };
};
