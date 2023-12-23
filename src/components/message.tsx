import { Markdown } from './markdown';
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
        className="prose break-words font-light"
        components={{
          p({ children }) {
            return <p className="mb-2 last:mb-0">{children}</p>;
          },
        }}
      >
        {message}
      </Markdown>
    </div>
  );
}
