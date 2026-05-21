export type TenantId = string;
export type UserId = string;

export interface ApiError {
   message: string;
   code?: string;
}

export const SHARED_TS_VERSION = "0.1.0";
