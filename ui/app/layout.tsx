import "./globals.css"
import { Space_Mono } from "next/font/google"
import { ThemeProvider } from "@/components/ThemeProvider"
import type React from "react"

const spaceMono = Space_Mono({
  weight: ["400", "700"],
  subsets: ["latin"],
  variable: "--font-space-mono",
})

export const metadata = {
  title: "typi",
  description: "A gallery of YouTube videos with infinite scrolling and dark mode",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${spaceMono.variable} font-mono`}>
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}

