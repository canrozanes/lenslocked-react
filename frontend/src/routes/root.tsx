import { Outlet } from "react-router-dom";

import Nav from "components/nav";
import { UserProvider } from "auth/user-provider";
import { getCsrf } from "api/csrf";
import { useEffect, useState } from "react";
import AxiosInterceptor from "hooks/axios-interceptor";

export default function Root() {
  const [isLoading, setIsLoading] = useState(true);
  useEffect(() => {
    getCsrf()
      .then(() => {
        setIsLoading(false);
      })
      .catch((e) => {
        console.log(e);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  if (isLoading) {
    // TODO
    return <></>;
  }

  return (
    <UserProvider>
      <Nav />
      <AxiosInterceptor>
        <Outlet />
      </AxiosInterceptor>
    </UserProvider>
  );
}
