"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "../store/authStore";
import {
  fetchMediaCalls,
  fetchMediaQoS,
  MediaCall,
  QoSStream,
} from "../lib/api";

export default function HomePage() {
  const router = useRouter();
  const { accessToken, userId, tenantId, clear } = useAuthStore();

  const [calls, setCalls] = useState<MediaCall[]>([]);
  const [qos, setQoS] = useState<QoSStream[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  if (!accessToken) {
    router.push("/login");
    return null;
  }

  useEffect(() => {
    let isMounted = true;

    async function load() {
      try {
        const [callsData, qosData] = await Promise.all([
          fetchMediaCalls(),
          fetchMediaQoS(),
        ]);
        if (!isMounted) return;
        setCalls(callsData);
        setQoS(qosData);
      } catch (e) {
        if (!isMounted) return;
        setError("Failed to load media data");
      } finally {
        if (isMounted) setLoading(false);
      }
    }

    load();

    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <main
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
        padding: "2rem",
        fontFamily:
          "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
        background: "#020617",
        color: "#e5e7eb",
      }}
    >
      <header
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <div>
          <h1 style={{ fontSize: "1.8rem", fontWeight: 600 }}>
            VoIP Cloud PBX Dashboard
          </h1>
          <p style={{ fontSize: "0.9rem", color: "#9ca3af" }}>
            Logged in as {userId} (tenant {tenantId})
          </p>
        </div>
        <button
          onClick={() => clear()}
          style={{
            padding: "0.4rem 0.9rem",
            borderRadius: "0.375rem",
            border: "none",
            background: "#ef4444",
            color: "#fff",
            cursor: "pointer",
          }}
        >
          Logout
        </button>
      </header>

      {loading && <p>Loading calls and QoS...</p>}
      {error && <p style={{ color: "#f87171" }}>{error}</p>}

      {!loading && !error && (
        <>
          <section>
            <h2 style={{ fontSize: "1.2rem", marginBottom: "0.5rem" }}>
              Active Calls
            </h2>
            {calls.length === 0 ? (
              <p style={{ color: "#9ca3af" }}>No calls yet.</p>
            ) : (
              <table
                style={{
                  width: "100%",
                  borderCollapse: "collapse",
                  fontSize: "0.9rem",
                }}
              >
                <thead>
                  <tr
                    style={{
                      textAlign: "left",
                      borderBottom: "1px solid #1f2937",
                    }}
                  >
                    <th>Call ID</th>
                    <th>From</th>
                    <th>To</th>
                    <th>Direction</th>
                    <th>Started</th>
                    <th>Ended</th>
                  </tr>
                </thead>
                <tbody>
                  {calls.map((c) => (
                    <tr
                      key={c.callId}
                      style={{ borderBottom: "1px solid #111827" }}
                    >
                      <td>{c.callId}</td>
                      <td>{c.fromUser}</td>
                      <td>{c.toUser}</td>
                      <td>{c.direction}</td>
                      <td>{c.startedAt}</td>
                      <td>{c.endedAt || "-"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </section>

          <section>
            <h2 style={{ fontSize: "1.2rem", marginBottom: "0.5rem" }}>
              QoS Streams
            </h2>
            {qos.length === 0 ? (
              <p style={{ color: "#9ca3af" }}>No RTP streams yet.</p>
            ) : (
              <table
                style={{
                  width: "100%",
                  borderCollapse: "collapse",
                  fontSize: "0.9rem",
                }}
              >
                <thead>
                  <tr
                    style={{
                      textAlign: "left",
                      borderBottom: "1px solid #1f2937",
                    }}
                  >
                    <th>SSRC</th>
                    <th>Addr</th>
                    <th>Packets</th>
                    <th>Lost</th>
                    <th>Last Seq</th>
                  </tr>
                </thead>
                <tbody>
                  {qos.map((s, idx) => (
                    <tr key={idx} style={{ borderBottom: "1px solid #111827" }}>
                      <td>{s.key.ssrc}</td>
                      <td>{s.key.addr}</td>
                      <td>{s.stats.packets}</td>
                      <td>{s.stats.lost}</td>
                      <td>{s.stats.lastSeq}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </section>
        </>
      )}
    </main>
  );
}
