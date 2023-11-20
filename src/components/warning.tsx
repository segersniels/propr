import { Alert, AlertDescription, AlertTitle } from 'components/ui/alert';
import { ExclamationTriangleIcon } from '@radix-ui/react-icons';

interface Props {
  className?: string;
  children?: React.ReactNode;
  variant?: 'default' | 'destructive';
}

function Warning(props: Props) {
  const { className, children, variant = 'default' } = props;

  return (
    <Alert className={className} variant={variant}>
      <ExclamationTriangleIcon className="h-4 w-4" />
      <AlertTitle>Heads up!</AlertTitle>
      <AlertDescription>{children}</AlertDescription>
    </Alert>
  );
}

Warning.FailedResponse = function FailedResponse(props: Props) {
  return <Warning {...props}>{props.children}</Warning>;
};

export default Warning;
