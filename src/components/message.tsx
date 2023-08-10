import { Markdown } from './markdown';
import remarkGfm from 'remark-gfm';
import { CodeBlock } from './ui/codeblock';
import { Button } from './ui/button';
import { AiOutlineCopy } from 'react-icons/ai';

interface Props {
  message: string;
}

export default function Message(props: Props) {
  const { message } = props;

  return (
    <div className="relative flex flex-col p-6 shadow-md rounded-md border border-gray-100">
      <Button
        variant="ghost"
        size="icon"
        className="absolute top-2 right-2 text-neutral-500"
        onClick={() => {
          return navigator.clipboard.writeText(message);
        }}
      >
        <AiOutlineCopy className="w-4 h-4" />
      </Button>

      <Markdown
        className="prose break-words dark:prose-invert font-light"
        remarkPlugins={[remarkGfm]}
        components={{
          p({ children }) {
            return <p className="mb-2 last:mb-0">{children}</p>;
          },
          code({ node, inline, className, children, ...props }) {
            if (children.length) {
              if (children[0] == '▍') {
                return (
                  <span className="mt-1 animate-pulse cursor-default">▍</span>
                );
              }

              children[0] = (children[0] as string).replace('`▍`', '▍');
            }

            const match = /language-(\w+)/.exec(className || '');

            if (inline) {
              return (
                <code
                  className="text-sm text-gray-700 bg-gray-100 p-1 rounded-lg"
                  {...props}
                >
                  {children}
                </code>
              );
            }

            return (
              <CodeBlock
                key={Math.random()}
                language={(match && match[1]) || ''}
                value={String(children).replace(/\n$/, '')}
                {...props}
              />
            );
          },
        }}
      >
        {message}
      </Markdown>
    </div>
  );
}
