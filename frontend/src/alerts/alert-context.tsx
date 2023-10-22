// AuthContext.js
import { createContext, useContext, useState, ReactNode } from "react";

const ALERT_TIME = 5000;
const initialState = {
  text: "",
  type: "",
};

const AlertContext = createContext({
  ...initialState,
  setAlert: (_text: string) => {},
});

export const AlertProvider = ({ children }: { children: ReactNode }) => {
  const [text, setText] = useState("");
  const [type, setType] = useState("");

  const setAlert = (text: string) => {
    setText(text);
    setType(type);

    setTimeout(() => {
      setText("");
      setType("");
    }, ALERT_TIME);
  };

  return (
    <AlertContext.Provider
      value={{
        text,
        type,
        setAlert,
      }}
    >
      {children}
    </AlertContext.Provider>
  );
};

export default function useAlert() {
  return useContext(AlertContext);
}
