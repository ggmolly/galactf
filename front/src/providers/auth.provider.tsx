import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { User } from "@/interfaces/user.interface";
import { apiClient } from "@/lib/axios";
import { toast } from "sonner";

interface AuthProviderProps {
  children: ReactNode;
}

interface AuthContextType {
  user: User;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

function redirectToForest() {
  localStorage.setItem("toastOnLogin", "true");

  const redirectCount = Number(localStorage.getItem("redirectCount") || "0");
  if (redirectCount >= 5) {
    alert("Redirect loop detected, try clearing your cookies and the local storage.");
    return;
  }

  localStorage.setItem("redirectCount", String(redirectCount + 1));
  window.location.href =
    "https://forest.galadrim.fr/login?redirect=" +
    encodeURIComponent(window.location.origin) +
    "&nonce=" +
    Math.random();
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User>({} as User);
  const [authStatus, setAuthStatus] = useState<"logging-in" | "logged-in" | "logged-out">(
    "logging-in"
  );

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await apiClient.get("/auth/me");
        if (response.status === 200) {
          setUser(response.data.data);
          setAuthStatus("logged-in");
          if (localStorage.getItem("toastOnLogin") === "true") {
            toast.success(`Logged in as ${response.data.data.name} !`, {
              dismissible: true,
              duration: 7000,
            });
            localStorage.removeItem("toastOnLogin");
            localStorage.removeItem("redirectCount");
          }
        } else {
          setAuthStatus("logged-out");
          redirectToForest();
        }
      } catch (err) {
        setAuthStatus("logged-out");
        redirectToForest();
      }
    };
    fetchUser();
  }, []);

  // Nothing must be rendered unless the user is logged in
  if (authStatus === "logging-in") {
    return <div>Logging in...</div>;
  }

  return <AuthContext.Provider value={{ user }}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
