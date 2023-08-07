import './globals.css';

import { Inter } from 'next/font/google';
import { Analytics } from '@vercel/analytics/react';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Propr',
  description:
    'Why wait until GitHub Copilot X is here? Generate proper PR descriptions now using AI!',
  icons: '/favicon.ico',
  twitter: {
    card: 'summary_large_image',
    title: 'Propr',
    description:
      'Why wait until GitHub Copilot X is here? Generate proper PR descriptions now using AI!',
    images: ['https://propr.dev/og-image.png'],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        {children}
        <Analytics />
      </body>
    </html>
  );
}
