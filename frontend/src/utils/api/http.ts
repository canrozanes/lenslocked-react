import axios from "axios";

let csrfToken = (
  document.getElementsByName("gorilla.csrf.Token")[0] as HTMLMetaElement
).content as string;

const instance = axios.create({
  baseURL: "/api",
  headers: { "X-CSRF-Token": csrfToken },
});

instance.interceptors.response.use(
  function (response) {
    // Any status code that lie within the range of 2xx cause this function to trigger
    // Do something with response data
    return response.data;
  },
  function (error) {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    // Do something with response error
    return Promise.reject(error);
  },
);

export default instance;
