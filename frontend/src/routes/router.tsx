import { createBrowserRouter } from "react-router-dom";
import ErrorPage from "routes/error-page";
import "App.css";
import Contact from "routes/contact";
import Faq from "routes/faq";
import Root from "routes/root";
import Home from "routes/home";
import SignUp from "routes/signup";

// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call
export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Home /> },
      {
        path: "contact",
        element: <Contact />,
      },
      {
        path: "faq",
        element: <Faq />,
      },
      {
        path: "signup",
        element: <SignUp />,
      },
    ],
  },
]);
