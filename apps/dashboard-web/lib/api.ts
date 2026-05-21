"use client";

import axios from "axios";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://api.localhost:8080";

export const apiClient = axios.create({
  baseURL: API_BASE,
  withCredentials: false,
});

export interface LoginResponse {
  userId: string;
  tenantId: string;
  accessToken: string;
  refreshToken: string;
}

export interface MediaCall {
  callId: string;
  tenantId: string;
  fromUser: string;
  toUser: string;
  direction: string;
  startedAt: string;
  endedAt?: string;
}

export interface QoSStream {
  key: {
    ssrc: number;
    addr: string;
  };
  stats: {
    packets: number;
    lost: number;
    lastSeq: number;
  };
}

export async function fetchMediaCalls(): Promise<MediaCall[]> {
  const res = await fetch("http://localhost:8082/calls");
  if (!res.ok) {
    throw new Error("Failed to load media calls");
  }
  return res.json();
}

export async function fetchMediaQoS(): Promise<QoSStream[]> {
  const res = await fetch("http://localhost:8082/qos");
  if (!res.ok) {
    throw new Error("Failed to load QoS");
  }
  return res.json();
}
