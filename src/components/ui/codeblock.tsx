'use client';

import { FC, memo } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { prism } from 'react-syntax-highlighter/dist/cjs/styles/prism';

interface Props {
  language: string;
  value: string;
}

interface languageMap {
  [key: string]: string | undefined;
}

export const programmingLanguages: languageMap = {
  javascript: '.js',
  python: '.py',
  java: '.java',
  c: '.c',
  cpp: '.cpp',
  'c++': '.cpp',
  'c#': '.cs',
  ruby: '.rb',
  php: '.php',
  swift: '.swift',
  'objective-c': '.m',
  kotlin: '.kt',
  typescript: '.ts',
  go: '.go',
  perl: '.pl',
  rust: '.rs',
  scala: '.scala',
  haskell: '.hs',
  lua: '.lua',
  shell: '.sh',
  sql: '.sql',
  html: '.html',
  css: '.css',
  // add more file extensions here, make sure the key is same as language prop in CodeBlock.tsx component
};

const CodeBlock: FC<Props> = memo(({ language, value }) => {
  return (
    <div className="w-full rounded-lg bg-muted">
      <SyntaxHighlighter
        language={language}
        style={prism}
        PreTag="div"
        showLineNumbers
        customStyle={{
          margin: 0,
          width: '100%',
          background: 'transparent',
          padding: '1.5rem 1rem',
        }}
        codeTagProps={{
          style: {
            fontSize: '0.9rem',
            fontFamily: `var(--font-geist-mono)`,
          },
        }}
      >
        {value}
      </SyntaxHighlighter>
    </div>
  );
});

CodeBlock.displayName = 'CodeBlock';

export { CodeBlock };
