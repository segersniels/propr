import { GitHubLogoIcon } from '@radix-ui/react-icons';

export default function Footer() {
  return (
    <div className="flex flex-col fixed bottom-0 w-full">
      <a
        href="https://github.com/segersniels/propr-cli"
        className="self-end"
        target="_blank"
      >
        <GitHubLogoIcon className="h-6 w-6 m-4" />
      </a>
    </div>
  );
}
