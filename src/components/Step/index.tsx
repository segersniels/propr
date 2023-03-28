import styles from './Step.module.css';
import ReactMarkdown from 'react-markdown';

const Step = ({ step, children }: { step: number; children: string }) => {
  return (
    <div className="flex items-center my-1 ml-2">
      <div className={styles.decimal}>{step}</div>
      <ReactMarkdown className={styles.step}>{children}</ReactMarkdown>
    </div>
  );
};

export default Step;
