export type TenantId = string;
export type UserId = string;
export type CallId = string;
export type UserRole = "ADMIN" | "TENANT_ADMIN" | "AGENT" | "VIEWER";
export interface Tenant {
    id: TenantId;
    name: string;
    createdAt: string;
    updatedAt: string;
}
export interface User {
    id: UserId;
    tenantId: TenantId;
    email: string;
    displayName: string;
    role: UserRole;
    createdAt: string;
    updatedAt: string;
}
export type CallDirection = "INBOUND" | "OUTBOUND";
export interface CallRecord {
    id: CallId;
    tenantId: TenantId;
    fromExtension: string;
    toExtension: string;
    direction: CallDirection;
    startedAt: string;
    endedAt?: string;
    durationSeconds?: number;
    status: "RINGING" | "ANSWERED" | "FAILED" | "CANCELLED";
}
export interface QoSMetrics {
    callId: CallId;
    packetLossPercent: number;
    jitterMs: number;
    mosScore?: number;
    updatedAt: string;
}
export interface AISummary {
    callId: CallId;
    summary: string;
    keywords: string[];
    sentiment?: "positive" | "neutral" | "negative";
    createdAt: string;
}
export interface ApiError {
    message: string;
    code?: string;
}
export declare const SHARED_TS_VERSION = "0.1.0";
