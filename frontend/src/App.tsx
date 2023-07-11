import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import ErrorPage from './routes/error-page';
import './App.css'

const router = createBrowserRouter([
  {
    path: "/",
    element: <div>This is a React App</div>,
    errorElement: <ErrorPage />,
  },
]);

function App() {

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <RouterProvider router={router}/>
    </>
  )
}

export default App
