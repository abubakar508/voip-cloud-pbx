import React from "react";

export const metadata = {
  title: "VoIP Cloud PBX Admin",
  description: "Admin console for VoIP Cloud PBX",
};

export default function RootLayout(props: { children: React.ReactNode }) {
  const { children } = props;
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
