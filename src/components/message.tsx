import { Markdown } from './markdown';
import remarkGfm from 'remark-gfm';
import { CodeBlock } from './ui/codeblock';

interface Props {
  message: string;
}

export default function Message(props: Props) {
  const { message } = props;

  return (
    <div
      className="cursor-copy flex flex-col p-4 pl-6 shadow-md rounded-2xl border border-gray-100 hover:bg-gray-50"
      onClick={() => {
        return navigator.clipboard.writeText(message);
      }}
    >
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
