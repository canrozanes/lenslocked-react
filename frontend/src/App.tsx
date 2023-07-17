import { RouterProvider } from "react-router-dom";
import { router } from "routes/router";
import "App.css";

function App() {
  return (
    <>
      {/* eslint-disable-next-line @typescript-eslint/no-unsafe-assignment */}
      <RouterProvider router={router} />
    </>
  );
}

export default App;
