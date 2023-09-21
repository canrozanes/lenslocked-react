import { AxiosError, AxiosResponse } from "axios";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import http from "utils/api/http";

const AxiosInterceptor = ({ children }: { children: JSX.Element }) => {
  const [isSet, setIsSet] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    function responseInterceptor(response: AxiosResponse) {
      // Any status code that lie within the range of 2xx cause this function to trigger
      // Do something with response data
      return response.data;
    }

    function errorInterceptor(error: AxiosError) {
      if (
        // /api/user/me is called on page load, we don't want to redirect if it fails
        // because otherwise all page loads with unauthenticated users would lead to /signin
        !error.request.responseURL.endsWith("/api/users/me") &&
        error.response?.status === 401
      ) {
        navigate("/signin");
      }
      // Any status codes that falls outside the range of 2xx cause this function to trigger
      // Do something with response error
      return Promise.reject(error);
    }

    const interceptor = http.interceptors.response.use(
      responseInterceptor,
      errorInterceptor,
    );

    setIsSet(true);

    return () => http.interceptors.response.eject(interceptor);
  }, [navigate]);

  return isSet && children;
};

export default AxiosInterceptor;
