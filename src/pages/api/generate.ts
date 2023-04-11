import { NextRequest } from 'next/server';
import { OpenAIStream } from 'helpers/Stream';
import { generatePrompt } from 'helpers/Prompt';
// @ts-expect-error
import wasm from 'resources/tiktoken_bg.wasm?module';
import model from '@dqbd/tiktoken/encoders/cl100k_base.json';
import { init, Tiktoken } from '@dqbd/tiktoken/lite/init';

if (!process.env.OPENAI_API_KEY) {
  throw new Error('Missing OPENAI_API_KEY environment variable');
}

export const config = {
  runtime: 'edge',
};

export default async function handler(req: NextRequest) {
  await init((imports) => WebAssembly.instantiate(wasm, imports));
  const encoding = new Tiktoken(
    model.bpe_ranks,
    model.special_tokens,
    model.pat_str
  );

  const body = await req.json();
  let prompt = generatePrompt(body.diff, body.template);

  // Check if exceeding model max token length and minify accordingly
  if (encoding.encode(prompt).length > 4096) {
    prompt = generatePrompt(body.diff, body.template, true);

    // Check if minified prompt is still too long
    if (encoding.encode(prompt).length > 4096) {
      return new Response(
        `The diff is too large (${
          encoding.encode(prompt).length
        }), try reducing the number of staged changes.`,
        {
          status: 400,
        }
      );
    }
  }

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

  // Free the encoding to prevent memory leaks
  encoding.free();

  return new Response(stream);
}
