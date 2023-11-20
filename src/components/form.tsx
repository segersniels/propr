'use client';

import { useChat } from 'ai/react';
import { useEffect, useRef, useState } from 'react';
import Step from './step';
import { Textarea } from './ui/textarea';
import { Button } from './ui/button';
import { AiOutlineLoading } from 'react-icons/ai';
import Message from './message';
import Warning from './warning';

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

export default function Form() {
  const [diff, setDiff] = useState('');
  const [template, setTemplate] = useState('');
  const ref = useRef<null | HTMLDivElement>(null);

  const { messages, setInput, handleSubmit, isLoading, error } = useChat({
    body: {
      diff,
      template,
    },
  });

  useEffect(() => {
    if (!diff) {
      return;
    }

    setInput(diff);
  }, [diff, setInput]);

  const scrollIntoView = () => {
    if (ref.current !== null) {
      ref.current.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const lastMessage = messages[messages.length - 1];
  const message =
    lastMessage?.role === 'assistant' ? lastMessage.content : null;

  return (
    <div className="flex flex-col w-full">
      <form
        onSubmit={(e) => {
          handleSubmit(e);
          scrollIntoView();

          /**
           * We reset the input to the original diff since the default behavior of `useChat`
           * is to clear the input after submitting
           */
          setInput(diff);
        }}
      >
        <Step step={1}>Navigate to your PR on GitHub</Step>
        <Step step={2}>Add `.diff` to the end of the URL</Step>
        <Step step={3}>Copy paste 🚀</Step>

        <Textarea
          className="min-h-[60px] w-full resize-none bg-transparent p-2 my-4 focus-within:outline-none sm:text-sm font-mono"
          value={diff}
          placeholder={DIFF_PLACEHOLDER}
          onChange={(event) => setDiff(event.target.value)}
          tabIndex={0}
          rows={10}
          spellCheck={false}
        />

        <Step step={4}>Provide a template</Step>

        <Textarea
          className="min-h-[60px] w-full resize-none bg-transparent p-2 my-4 focus-within:outline-none sm:text-sm font-mono"
          value={template}
          placeholder={TEMPLATE_PLACEHOLDER}
          onChange={(event) => setTemplate(event.target.value)}
          tabIndex={1}
          rows={4}
        />

        <Button className="w-full" disabled={!diff.length}>
          {isLoading ? (
            <AiOutlineLoading className="mx-2 animate-spin stroke-[3rem] font-bold" />
          ) : (
            'Generate'
          )}
        </Button>
      </form>

      {!!error && (
        <Warning.FailedResponse className="my-4" variant="destructive">
          {error.message}
        </Warning.FailedResponse>
      )}

      {!!message && (
        <>
          <hr className="my-4 w-64 mx-auto" />

          <Message message={message} />
        </>
      )}
    </div>
  );
}
