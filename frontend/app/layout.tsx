import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { QueryProvider } from "@/lib/query-provider";
import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: {
    default: "The Asheville Setlist",
    template: "%s | The Asheville Setlist",
  },
  description:
    "Discover concerts and live music in Asheville, NC. Find upcoming shows, explore venues, and stay connected to the local music scene.",
  keywords: [
    "Asheville concerts",
    "live music",
    "music venues",
    "Asheville NC",
    "shows",
    "events",
  ],
  openGraph: {
    type: "website",
    locale: "en_US",
    siteName: "The Asheville Setlist",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <QueryProvider>
          <div className="flex min-h-screen flex-col">
            <Header />
            <main className="flex-1">{children}</main>
            <Footer />
          </div>
        </QueryProvider>
      </body>
    </html>
  );
}
