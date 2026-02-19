import type { Metadata } from "next";
import "./globals.css";
import "react-responsive-modal/styles.css";
import { Toaster } from "react-hot-toast";

export const metadata: Metadata = {
  title: "Library Automation",
  description: "Library Automation Software",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
        <Toaster />
      </body>
    </html>
  );
}
