"use client";

import { useEffect, useRef, useState } from "react";

interface WebSocketState {
  status: "disconnected" | "connecting" | "connected" | "error";
  send: (data: string) => void;
  lastMessage: string | null;
}

export function useWebSocket(url: string | null): WebSocketState {
  const [status, setStatus] =
    useState<WebSocketState["status"]>("disconnected");
  const [lastMessage, setLastMessage] = useState<string | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!url) {
      setStatus("disconnected");
      return;
    }

    setStatus("connecting");
    const ws = new WebSocket(url);
    socketRef.current = ws;

    ws.onopen = () => {
      setStatus("connected");
    };

    ws.onmessage = (event) => {
      setLastMessage(event.data);
    };

    ws.onerror = () => {
      setStatus("error");
    };

    ws.onclose = () => {
      setStatus("disconnected");
    };

    return () => {
      ws.close();
      socketRef.current = null;
    };
  }, [url]);

  const send = (data: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(data);
    }
  };

  return { status, send, lastMessage };
}
