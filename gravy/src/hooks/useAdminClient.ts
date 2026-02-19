import { API } from "@/config";
import axios from "axios";

export const useAdminClient = () => {
  let token: string | undefined | null;

  if (typeof window !== "undefined") {
    token = localStorage.getItem("admin_token");
  }

  const client = axios.create({
    baseURL: API.BASE_URL,
    headers: {
      Authorization: `Bearer ${token}` || "",
    },
  });

  return client;
};
