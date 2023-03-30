import { NextRequest } from 'next/server';
import { OpenAIStream } from 'helpers/Stream';
import { prepareDiff } from 'helpers/Prompt';

if (!process.env.OPENAI_API_KEY) {
  throw new Error('Missing OPENAI_API_KEY environment variable');
}

export const config = {
  runtime: 'edge',
};

export default async function handler(req: NextRequest) {
  const body = await req.json();
  const prompt = `
    Generate a concise PR description from the provided git diff according to a provided template.
    The PR description should be a good summary of the changes made.
    Do not reference each file and function added but rather give a general explanation of the changes made.
    Don't mention each change individually, but rather group them together.
    You are free to make a calculated guess as to which changes and files are related to each other so you can group them together.
    When endpoints/routes are in the diff try to reference these when describing features.
    It's not worth mentioning that you added tests when you mention you added a new feature as it implies that you added tests in the first place.
    If notes or reason why the change happened are requested, make sure you try to explain the reasoning without using too much technical jargon.
    If a section has no content, you can leave the entire section out.
    If the diff provided is not actually a diff I want you to respond with an appropriate message accordingly.

    The PR description should be structured as follows: """
    ${body.template}
    """

    Here is the diff: """
    ${prepareDiff(body.diff)}
    """
  `;

  const stream = await OpenAIStream({
    model: 'gpt-3.5-turbo',
    messages: [{ role: 'user', content: prompt }],
    temperature: 0.7,
    top_p: 1,
    frequency_penalty: 0,
    presence_penalty: 0,
    max_tokens: 300,
    stream: true,
    n: 1,
  });

  return new Response(stream);
}
