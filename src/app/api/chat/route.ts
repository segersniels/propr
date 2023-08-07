import { OpenAIStream, StreamingTextResponse } from 'ai';
import { Configuration, OpenAIApi } from 'openai-edge';

export const runtime = 'edge';

function generateSystemMessage(template: string) {
  return `You will be asked to write a descriptive GitHub PR description based on a provided git diff.
    Analyze the code changes and provide a concise explanation of the changes, their context and why they were made.
    Don't reference file names directly, instead focus on explaining the changes in a broader context.
    Do not treat imports and requires as changes or new features.

    Use the following template to write your description:
    """
    ${template}
    """

    If a section from the template does not apply (no significant changes in that category), omit that section from your final output.
  `;
}

export async function POST(req: Request) {
  const { diff, template } = await req.json();
  const config = new Configuration({
    apiKey: process.env.OPENAI_API_KEY,
  });

  const openai = new OpenAIApi(config);
  const response = await openai.createChatCompletion({
    model: 'gpt-3.5-turbo-16k',
    stream: true,
    messages: [
      {
        role: 'system',
        content: generateSystemMessage(template),
      },
      {
        role: 'user',
        content: diff,
      },
    ],
  });

  const stream = OpenAIStream(response);

  return new StreamingTextResponse(stream);
}
