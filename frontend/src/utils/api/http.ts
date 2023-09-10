import axios from "axios";

const instance = axios.create({
  baseURL: "/api",
});

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
