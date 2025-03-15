import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";

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

  useEffect(() => {
    let reconnectTimeout: NodeJS.Timeout;

    const connectWebSocket = () => {
      console.log("[ws] attempting to establish connection...");
      const scheme = import.meta.env.VITE_API_URL.startsWith("https")
        ? "wss"
        : "ws";
      const ws = new WebSocket(
        `${scheme}://${import.meta.env.VITE_API_URL.split("://")[1]}/ws`
      );

      ws.onopen = () => {
        console.log("[ws] connected");
        setReconnectDelay(1000);
      };

      ws.onclose = (event) => {
        console.warn("[ws] disconnected:", event);
        let nextDelay = reconnectDelay;
        if (event.wasClean) {
          nextDelay = 1000; // Reduce delay for clean closes
        } else {
          nextDelay = Math.min(reconnectDelay + 1000, 10000);
        }
        setReconnectDelay(nextDelay);
        console.log(`[ws] reconnecting in ${nextDelay / 1000} seconds...`);
        reconnectTimeout = setTimeout(connectWebSocket, nextDelay);
      };

      ws.onerror = (err) => {
        console.error("[ws] error:", err);
      };

      ws.onmessage = (event) => {
        console.log("[ws] message received:", event.data);
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
      clearTimeout(reconnectTimeout);
    };
  }, []);

  return (
    <WsContext.Provider value={{ ws: wsRef.current }}>
      {children}
    </WsContext.Provider>
  );
};

export const useWs = (): WsContextType => {
  const context = useContext(WsContext);
  if (context === undefined) {
    throw new Error("useWs must be used within a WsProvider");
  }
  return context;
};
