import axios from "axios";

export function axiosErrorFactory(err: any) {
  let errorMessage = "Network error. Please check your internet connection.";

  if (axios.isAxiosError(err) && err.response) {
    errorMessage = `[HTTP ${err.response.status}]: ${err.response.data.message}`;
  }
  return errorMessage;
}
