'use client';

import { experimental_useAssistant as useAssistant } from 'ai/react';
import { useEffect, useRef, useState } from 'react';
import Step from './step';
import { Textarea } from './ui/textarea';
import { Button } from './ui/button';
import { AiOutlineLoading } from 'react-icons/ai';
import Message from './message';
import Warning from './warning';
import useLocalStorage from 'hooks/use-local-storage';

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
  const [shouldShowError, setShouldShowError] = useState(false);
  const [template, setTemplate] = useLocalStorage('template', '');
  const ref = useRef<null | HTMLDivElement>(null);

  const { status, messages, input, submitMessage, handleInputChange, error } =
    useAssistant({
      api: '/api/assistant',
      body: {
        template,
      },
    });

  const scrollIntoView = () => {
    if (ref.current !== null) {
      ref.current.scrollIntoView({ behavior: 'smooth' });
    }
  };

  useEffect(() => {
    if (!error) {
      return setShouldShowError(false);
    }

    setShouldShowError(true);
  }, [error]);

  const isLoading = status === 'in_progress';
  const lastMessage = messages[messages.length - 1];
  const message =
    lastMessage?.role === 'assistant' ? lastMessage.content : null;

  return (
    <div className="flex flex-col w-full">
      <form
        onSubmit={(e) => {
          submitMessage(e);
          scrollIntoView();
        }}
      >
        <Step step={1}>Provide a template</Step>

        <Textarea
          className="min-h-[60px] w-full resize-none bg-transparent p-2 my-4 focus-within:outline-none sm:text-sm font-mono"
          value={template}
          placeholder={TEMPLATE_PLACEHOLDER}
          onChange={(event) => setTemplate(event.target.value)}
          tabIndex={1}
          rows={4}
        />

        <Step step={2}>Navigate to your PR on GitHub</Step>
        <Step step={3}>Add `.diff` to the end of the URL</Step>
        <Step step={4}>Copy paste 🚀</Step>

        <Textarea
          className="min-h-[60px] w-full resize-none bg-transparent p-2 my-4 focus-within:outline-none sm:text-sm font-mono"
          value={input}
          placeholder={DIFF_PLACEHOLDER}
          tabIndex={0}
          rows={10}
          spellCheck={false}
          onChange={handleInputChange}
        />

        <Button className="w-full" disabled={!input.length}>
          {isLoading ? (
            <AiOutlineLoading className="mx-2 animate-spin stroke-[3rem] font-bold" />
          ) : (
            'Generate'
          )}
        </Button>
      </form>

      {shouldShowError && (
        <Warning.FailedResponse className="my-4" variant="destructive">
          {'message' in (error as Error)
            ? (error as Error).message
            : (error as any).toString()}
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
