import { RouterProvider } from "react-router-dom";
import { router } from "routes/router";
import "App.css";

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      {/* eslint-disable-next-line @typescript-eslint/no-unsafe-assignment */}
      <RouterProvider router={router} />
    </div>
  );
}

export default App;
