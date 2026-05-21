"use client";

import React from "react";
import type { Tenant, UserRole } from "@shared-ts/index";

const roles: UserRole[] = ["ADMIN", "TENANT_ADMIN", "AGENT", "VIEWER"];

const demoTenants: Tenant[] = [
  {
    id: "tenant-1",
    name: "Acme Corp",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: "tenant-2",
    name: "Globex Ltd",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
];

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
      <p style={{ color: "#555", marginBottom: "1rem" }}>
        Tenants and roles use shared types from <code>@shared-ts</code>.
      </p>
      <div style={{ display: "flex", gap: "2rem" }}>
        <div>
          <h2 style={{ fontSize: "1rem", fontWeight: 600 }}>Tenants</h2>
          <ul>
            {demoTenants.map((t) => (
              <li key={t.id}>{t.name}</li>
            ))}
          </ul>
        </div>
        <div>
          <h2 style={{ fontSize: "1rem", fontWeight: 600 }}>Roles</h2>
          <ul>
            {roles.map((r) => (
              <li key={r}>{r}</li>
            ))}
          </ul>
        </div>
      </div>
    </main>
  );
}
