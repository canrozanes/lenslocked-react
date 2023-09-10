import axios from "axios";

let csrfToken = "";

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

instance.interceptors.request.use(
  function (config) {
    const token = getCookie("csrf");

    if (token) {
      config.headers["x-csrf-token"] = token;
    }
    return config;
  },
  function (error) {
    return Promise.reject(error);
  },
);

export default instance;

function getCookie(name: string) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);

  if (!parts) {
    return;
  }

  if (parts.length === 2) {
    return parts.pop()?.split(";").shift() ?? "";
  }
}
