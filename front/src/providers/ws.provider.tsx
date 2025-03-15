import {
  callHandler,
  clearEventHandlers,
  registerEventHandler,
  WS_CHALLENGE_ATTEMPT,
} from "@/proto/handlers";
import { createContext, ReactNode, useContext, useEffect, useRef, useState } from "react";
import { useChallenges } from "./challenges.provider";
import { ChallengeAttempt } from "@/proto/challenge_attempt";
import { handleChalAttempt } from "@/events/challenge_attempt.event";
import { useAuth } from "./auth.provider";

interface WsProviderProps {
  children: ReactNode;
}

interface WsContextType {
  ws: WebSocket | null;
}

const WsContext = createContext<WsContextType | undefined>(undefined);

export const WsProvider: React.FC<WsProviderProps> = ({ children }) => {
  const wsRef = useRef<WebSocket | null>(null);
  const [reconnectDelay, setReconnectDelay] = useState(1000);
  const { challenges, setChallenges } = useChallenges();
  const { user } = useAuth();

  const challengesRef = useRef(challenges);
  useEffect(() => {
    challengesRef.current = challenges;
  }, [challenges]);

  useEffect(() => {
    if (!user.id) {
      return;
    }
    let reconnectTimeout: NodeJS.Timeout;

    const connectWebSocket = () => {
      console.log("[ws] attempting to establish connection...");
      const scheme = import.meta.env.VITE_API_URL.startsWith("https") ? "wss" : "ws";
      const ws = new WebSocket(`${scheme}://${import.meta.env.VITE_API_URL.split("://")[1]}/ws`);
      ws.binaryType = "arraybuffer";

      ws.onopen = () => {
        console.log("[ws] connected");
        setReconnectDelay(1000);
        console.log("[ws] registering event handlers...");
        clearEventHandlers();
        registerEventHandler(WS_CHALLENGE_ATTEMPT, handleChalAttempt, ChallengeAttempt.decode, {
          challenges: challengesRef,
          setChallenges,
          user,
        });
      };

      ws.onclose = (event) => {
        console.warn("[ws] disconnected:", event);
        const nextDelay = event.wasClean ? 1000 : Math.min(reconnectDelay + 1000, 10000);
        setReconnectDelay(nextDelay);
        console.log(`[ws] reconnecting in ${nextDelay / 1000} seconds...`);
        reconnectTimeout = setTimeout(connectWebSocket, nextDelay);
      };

      ws.onerror = (err) => {
        console.error("[ws] error:", err);
      };

      ws.onmessage = (event) => {
        const data = new Uint8Array(event.data);
        const eventId = data[0];
        console.info(`[ws] event id: 0x${eventId.toString(16)}`);
        callHandler(eventId, data);
      };

      wsRef.current = ws;
    };

    connectWebSocket();

    return () => {
      if (wsRef.current) {
        console.log("[ws] cleaning up connection...");
        wsRef.current.close();
        wsRef.current = null;
      }
      console.log("[ws] cleaning up handlers...");
      clearEventHandlers();
      clearTimeout(reconnectTimeout);
    };
  }, [user]);

  return <WsContext.Provider value={{ ws: wsRef.current }}>{children}</WsContext.Provider>;
};

export const useWs = (): WsContextType => {
  const context = useContext(WsContext);
  if (context === undefined) {
    throw new Error("useWs must be used within a WsProvider");
  }
  return context;
};
