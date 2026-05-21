"use client";

import React from "react";

export default function AdminHomePage() {
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
        VoIP Cloud PBX Admin
      </h1>
      <p style={{ color: "#555" }}>
        Admin frontend scaffold ready. Tenant and user management will be
        implemented in later phases.
      </p>
    </main>
  );
}
