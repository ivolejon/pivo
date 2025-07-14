import axios from "axios";

export const httpClient = axios.create({
  baseURL: import.meta.env.VITE_GO_API_URL,
  headers: {
    "Content-type": "application/json",
  },
});
