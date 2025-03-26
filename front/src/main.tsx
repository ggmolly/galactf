import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { AuthProvider } from "./providers/auth.provider.tsx";
import { Toaster } from "./components/ui/sonner.tsx";
import { ChallengesProvider } from "./providers/challenges.provider.tsx";
import { WsProvider } from "./providers/ws.provider.tsx";
import { ThemeProvider } from "./providers/theme.provider.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="dark">
      <AuthProvider>
        <ChallengesProvider>
          <WsProvider>
            <App />
          </WsProvider>
        </ChallengesProvider>
      </AuthProvider>
      <Toaster position="top-right" richColors className="pointer-events-auto" />
    </ThemeProvider>
  </StrictMode>
);
