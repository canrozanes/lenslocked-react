import { createBrowserRouter } from "react-router-dom";
import ErrorPage from "routes/error-page";
import "App.css";
import Contact from "routes/contact";
import Faq from "routes/faq";
import Root from "routes/root";
import Home from "routes/home";
import SignUp from "routes/signup";
import SignIn from "routes/sign-in";
import ForgotPassword from "routes/forgot-password";
import ResetPassword from "routes/reset-password";
import RequireAuth from "auth/require-auth";

// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call
export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        index: true,
        element: (
          <RequireAuth>
            <Home />
          </RequireAuth>
        ),
      },
      {
        path: "contact",
        element: (
          <RequireAuth>
            <Contact />
          </RequireAuth>
        ),
      },
      {
        path: "faq",
        element: (
          <RequireAuth>
            <Faq />
          </RequireAuth>
        ),
      },
      {
        path: "signup",
        element: <SignUp />,
      },
      {
        path: "signin",
        element: <SignIn />,
      },
      {
        path: "forgot-pw",
        element: <ForgotPassword />,
      },
      {
        path: "reset-pw",
        element: <ResetPassword />,
      },
      {
        path: "*",
        element: <ErrorPage />,
      },
    ],
  },
]);
