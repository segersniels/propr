import { experimental_useAssistant as useAssistant } from 'ai/react';
import Warning from './warning';
import useLocalStorage from 'hooks/use-local-storage';
  const [shouldShowError, setShouldShowError] = useState(false);
  const [template, setTemplate] = useLocalStorage('template', '');
  const { status, messages, input, submitMessage, handleInputChange, error } =
    useAssistant({
      api: '/api/assistant',
      body: {
        template,
      },
    });
  useEffect(() => {
    if (!error) {
      return setShouldShowError(false);
    }

    setShouldShowError(true);
  }, [error]);

  const isLoading = status === 'in_progress';
          submitMessage(e);
        <Step step={1}>Provide a template</Step>
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
      {shouldShowError && (
        <Warning.FailedResponse className="my-4" variant="destructive">
          {'message' in (error as Error)
            ? (error as Error).message
            : (error as any).toString()}
        </Warning.FailedResponse>
      )}
