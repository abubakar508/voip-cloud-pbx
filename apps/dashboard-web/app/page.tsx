"use client";

import React from "react";
import type { User, Tenant, CallRecord } from "@shared-ts/index";

const exampleTenant: Tenant = {
  id: "tenant-demo",
  name: "Demo Tenant",
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

const exampleUser: User = {
  id: "user-demo",
  tenantId: exampleTenant.id,
  email: "demo@tenant.local",
  displayName: "Demo User",
  role: "TENANT_ADMIN",
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

const exampleCall: CallRecord = {
  id: "call-demo",
  tenantId: exampleTenant.id,
  fromExtension: "1001",
  toExtension: "1002",
  direction: "OUTBOUND",
  startedAt: new Date().toISOString(),
  status: "RINGING",
};

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
      <p style={{ color: "#555", marginBottom: "1rem" }}>
        Shared TypeScript types loaded from <code>@shared-ts</code>.
      </p>
      <pre
        style={{
          fontSize: "0.8rem",
          background: "#111",
          color: "#0f0",
          padding: "1rem",
          borderRadius: "0.5rem",
          maxWidth: "600px",
          overflowX: "auto",
        }}
      >
        {JSON.stringify(
          {
            tenant: exampleTenant,
            user: exampleUser,
            call: exampleCall,
          },
          null,
          2,
        )}
      </pre>
    </main>
  );
}
