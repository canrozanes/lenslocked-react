import { Outlet } from "react-router-dom";

import Nav from "components/nav";
import { UserProvider } from "auth/user-provider";

export default function Root() {
  return (
    <UserProvider>
      <Nav />
      <Outlet />
    </UserProvider>
  );
}
