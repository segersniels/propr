import './globals.css';

import { Inter } from 'next/font/google';
import { Analytics } from '@vercel/analytics/react';
import Footer from 'components/footer';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Propr',
  description:
    'Why wait until copilot for pull requests is here? Generate proper PR descriptions based on your git diff now with the power of AI!',
  icons: '/favicon.ico',
  twitter: {
    card: 'summary_large_image',
    title: 'Propr',
    description:
      'Why wait until copilot for pull requests is here? Generate proper PR descriptions based on your git diff now with the power of AI!',
    images: ['https://propr.dev/og-image.png'],
  },
  openGraph: {
    title: 'Propr',
    description:
      'Why wait until copilot for pull requests is here? Generate proper PR descriptions based on your git diff now with the power of AI!',
    images: ['https://propr.dev/og-image.png'],
    url: 'https://propr.dev',
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
        <main className="flex w-full min-h-screen items-center justify-center">
          {children}
        </main>

        <Footer />
        <Analytics />
      </body>
    </html>
  );
}
