import { Html, Head, Main, NextScript } from 'next/document';

export default function Document() {
  return (
    <Html lang="en">
      <Head>
        <link
          rel="icon"
          href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>📝</text></svg>"
        />

        <meta
          name="description"
          content="Generate proper pull request descriptions"
        />
        <meta property="og:site_name" content="propr.dev" />
        <meta
          property="og:description"
          content="Generate proper pull request descriptions"
        />
        <meta property="og:title" content="Propr" />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content="Propr" />
        <meta
          name="twitter:description"
          content="Generate proper pull request descriptions"
        />
      </Head>
      <body>
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
