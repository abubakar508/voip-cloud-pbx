"use client";

import React from "react";

export default function HomePage() {
  return (
    <main
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        fontFamily:
          "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
      }}
    >
      <h1 style={{ fontSize: "2rem", fontWeight: 600, marginBottom: "0.5rem" }}>
        VoIP Cloud PBX Dashboard
      </h1>
      <p style={{ color: "#555" }}>
        Frontend scaffold ready. SIP softphone, analytics, and live dashboards
        will be added in later phases.
      </p>
    </main>
  );
}
