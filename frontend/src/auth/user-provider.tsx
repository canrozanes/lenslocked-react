import React, { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { User } from "models/user";
import { getMe, signOut } from "api/user";

type UserContext = {
  user: User | null;
  setUser: (user: User) => void;
  isUserLoading: boolean;
  setIsUserLoading: (isLoading: boolean) => void;
  handleSignOut: () => Promise<void>;
  isSigningOut: boolean;
};

export const UserContext = React.createContext<UserContext | null>(null);

export default function useUserContext() {
  const value = useContext(UserContext);

  if (!value) {
    throw new Error("User context has not been set yet");
  }
  return value;
}

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isUserLoading, setIsUserLoading] = useState(true);
  const [isSigningOut, setIsSigningOut] = useState(false);

  let navigate = useNavigate();

  const handleSignOut = async () => {
    setIsSigningOut(true);

    try {
      await signOut();
      setUser(null);
      navigate("/signin", { replace: true });
    } catch (e) {
      console.log(e);
    } finally {
      setIsSigningOut(false);
    }
  };

  useEffect(() => {
    getMe()
      .then((res) => {
        setUser(res.user);
        debugger;
        if (res.user) {
          navigate("/", { replace: true });
        }
      })
      .catch((e) => {
        console.log(e);
      })
      .finally(() => {
        setIsUserLoading(false);
      });
  }, []);

  return (
    <UserContext.Provider
      value={{
        user,
        setUser,
        setIsUserLoading,
        isUserLoading: isUserLoading,
        handleSignOut: handleSignOut,
        isSigningOut,
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
