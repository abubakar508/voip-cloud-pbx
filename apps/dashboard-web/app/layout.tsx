import React from "react";
import { AppProviders } from "./providers";

export const metadata = {
  title: "VoIP Cloud PBX Dashboard",
  description: "Operator dashboard for VoIP Cloud PBX",
};

export default function RootLayout(props: { children: React.ReactNode }) {
  const { children } = props;
  return (
    <html lang="en">
      <body>
        <AppProviders>{children}</AppProviders>
      </body>
    </html>
  );
}
