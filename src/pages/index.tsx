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
import ReactMarkdown from 'react-markdown';
import Step from 'components/Step';

const DIFF_PLACEHOLDER = `diff --git a/.docker/cassandra/Dockerfile b/.docker/cassandra/Dockerfile
new file mode 100644
index 0000000000..2d20e27312
--- /dev/null
+++ b/.docker/cassandra/Dockerfile
@@ -0,0 +1,6 @@
+FROM cassandra:3.11
+
+COPY entrypoint.sh /entrypoint.sh
+
+ENTRYPOINT ["/entrypoint.sh"]
+CMD ["cassandra", "-f"]`;

const TEMPLATE_PLACEHOLDER = `### Added
### Removed
### Fixed`;

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
        <p className={styles.subtitle}>
          Why wait until GitHub Copilot X is here...
        </p>
        <h1 className={styles.title}>
          Generate your next Pull Request description using ChatGPT
        </h1>

        <form className={styles.form} onSubmit={handleSubmit}>
          <Step step={1}>Navigate to your PR on GitHub</Step>
          <Step step={2}>Add `.diff` to the end of the URL</Step>
          <Step step={3}>Copy paste 🚀</Step>

          <TextArea
            value={diff}
            placeholder={DIFF_PLACEHOLDER}
            onChange={(event) => setDiff(event.target.value)}
            rows={10}
          />

          <Step step={4}>Provide a template</Step>

          <TextArea
            value={template}
            placeholder={TEMPLATE_PLACEHOLDER}
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
