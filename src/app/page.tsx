import Form from 'components/form';

export default function Page() {
  return (
    <div className="flex max-w-prose flex-col items-center justify-center p-4 md:py-24">
      <h1 className="text-5xl md:text-6xl font-bold text-center tracking-tighter mb-8">
        Generate your next proper PR description
      </h1>

      <Form />
    </div>
  );
}
