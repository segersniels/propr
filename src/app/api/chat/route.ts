import { OpenAIStream, StreamingTextResponse } from 'ai';
import { OpenAI } from 'openai';

export const runtime = 'edge';

function generateSystemMessage(template: string) {
  return `You will be asked to write a concise GitHub PR description based on a provided git diff.
    Analyze the code changes and provide a concise explanation of the changes, their context and why they were made.
    Don't reference file names or directories directly, instead give a general explanation of the changes made.
    Do not treat imports and requires as changes or new features.

    Use the following template to write your description:
    """
    ${template}
    """

    If you can't determine changes for a specific section in the template just omit that section out entirely.
    If the provided message is not a diff respond with an appropriate message.
  `;
}

export async function POST(req: Request) {
  const { diff, template } = await req.json();
  const openai = new OpenAI({
    apiKey: process.env.OPENAI_API_KEY,
  });

  const response = await openai.chat.completions.create({
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
