import { Html, Head, Main, NextScript } from 'next/document';

export default function Document() {
  const description =
    'Why wait until GitHub Copilot X is here? Generate proper PR descriptions now using AI!';

  return (
    <Html lang="en">
      <Head>
        <link rel="icon" href="/favicon.ico" />

        <meta name="description" content={description} />
        <meta property="og:site_name" content="propr.dev" />
        <meta property="og:description" content={description} />
        <meta property="og:title" content="Propr" />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content="Propr" />
        <meta name="twitter:description" content={description} />

        <meta property="og:image" content="https://propr.dev/og-image.png" />
        <meta name="twitter:image" content="https://propr.dev/og-image.png" />
      </Head>
      <body>
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
