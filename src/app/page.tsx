import Form from 'components/form';

export default function Page() {
  return (
    <div className="flex flex-row mx-auto items-center justify-center p-4 min-h-screen">
      <div className="flex flex-col flex-1 max-w-xl w-full items-center">
        <h1 className="text-5xl md:text-7xl font-bold text-center tracking-tighter mb-8">
          Generate your next proper PR description
        </h1>

        <Form />
      </div>
    </div>
  );
}
