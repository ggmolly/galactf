import axios from "axios";

export function axiosErrorFactory(err: any) {
    let errorMessage = "Network error. Please check your internet connection.";
    console.log(err)
    if (!axios.isAxiosError(err)) {
        errorMessage = `HTTP ${err.response.status}, failed to load challenges`;
    }
    return errorMessage;
}