import styles from './Step.module.css';
import ReactMarkdown from 'react-markdown';

export default function Step({
  step,
  children,
}: {
  step: number;
  children: string;
}) {
  return (
    <div className="flex items-center my-1 ml-2">
      <div className="rounded-full text-white bg-black flex items-center justify-center w-6 h-6 mr-2">
        {step}
      </div>
      <ReactMarkdown className="list-decimal list-inside font-light text-lg antialiased">
        {children}
      </ReactMarkdown>
    </div>
  );
}
