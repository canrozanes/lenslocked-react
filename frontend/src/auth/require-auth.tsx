import useUserContext from "auth/user-provider";
import { Navigate, useLocation } from "react-router-dom";

function RequireAuth({ children }: { children: JSX.Element }) {
  let { user, isUserLoading } = useUserContext();
  let location = useLocation();

  if (isUserLoading) {
    return <p>Loading...</p>;
  }

  if (!user) {
    // Redirect them to the /login page, but save the current location they were
    // trying to go to when they were redirected. This allows us to send them
    // along to that page after they login, which is a nicer user experience
    // than dropping them off on the home page.
    return <Navigate to="/signin" state={{ from: location }} replace />;
  }

  return children;
}

export default RequireAuth;

// Any page that we land on should call /api/me.
//  If that request returns a 401, we should go to /signin.
//  If that request succeeds, we should stay in that page.

// All other calls, if they return a 400 error axios interceptor should kick in
// if it's a logged in page, 401 should return us to login
// if it's a public page, 401 shouldn't return us to login
