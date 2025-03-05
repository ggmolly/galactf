import axios from "axios";

export const apiClient = axios.create({
  baseURL: "http://localhost:7777/api/v1",
});