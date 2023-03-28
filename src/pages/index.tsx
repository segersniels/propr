import Head from 'next/head';
import {
  DetailedHTMLProps,
  FormEvent,
  TextareaHTMLAttributes,
  useCallback,
  useState,
} from 'react';
import styles from 'styles/Home.module.css';
import Footer from 'components/Footer';
import { AiOutlineLoading } from 'react-icons/ai';
import ReactMarkdown from 'react-markdown'

const TextArea = (
  props: DetailedHTMLProps<
    TextareaHTMLAttributes<HTMLTextAreaElement>,
    HTMLTextAreaElement
  >
) => {
  return <textarea className={styles.textarea} {...props} />;
};

const GenerateButton = ({
  diff,
  isGenerating,
}: {
  diff: string;
  isGenerating: boolean;
}) => {
  if (isGenerating) {
    return (
      <button
        type="submit"
        className={styles.button}
        disabled={!diff.trim().length}
      >
        <AiOutlineLoading className="animate-spin font-bold mx-2 stroke-[3rem]" />
      </button>
    );
  }

  return (
    <button
      type="submit"
      className={styles.button}
      disabled={!diff.trim().length}
    >
      Generate
    </button>
  );
};

export default function Home() {
  const [diff, setDiff] = useState('');
  const [template, setTemplate] = useState('');
  const [message, setMessage] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);

  const handleSubmit = useCallback(
    async (event: FormEvent<HTMLFormElement>) => {
      try {
        event.preventDefault();
        setIsGenerating(true);

        const response = await fetch('/api/generate', {
          method: 'POST',
          body: JSON.stringify({ template, diff }),
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) {
          return;
        }

        const data = await response.json();
        setMessage(data.message);
      } finally {
        setIsGenerating(false);
      }
    },
    [diff, template]
  );

  return (
    <div className={styles.container}>
      <Head>
        <title>Propr</title>
      </Head>

      <div className={styles.wrapper}>
        <h1 className={styles.title}>
          Generate proper PR descriptions 📝
        </h1>

        <form onSubmit={handleSubmit}>
          <TextArea
            value={diff}
            placeholder="Paste your diff here"
            onChange={(event) => setDiff(event.target.value)}
            rows={10}
          />

          <TextArea
            value={template}
            placeholder="Provide a template"
            onChange={(event) => setTemplate(event.target.value)}
            rows={4}
          />

          <GenerateButton diff={diff} isGenerating={isGenerating} />
        </form>

        {message && (
          <>
            <hr className="my-4 w-64 mx-auto" />

            <div
              className={styles.message}
              onClick={() => {
                return navigator.clipboard.writeText(message);
              }}
            >
              <ReactMarkdown>{message}</ReactMarkdown>
            </div>
          </>
        )}
      </div>

      <Footer />
    </div>
  );
}
