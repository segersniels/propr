import { NextRequest } from 'next/server';
import { OpenAIStream } from 'helpers/Stream';
import * as PromptHelper from 'helpers/Prompt';
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
  const prompt = PromptHelper.generatePrompt(body.diff, body.template);

  // Within model token length so pass entire diff
  if (encoding.encode(prompt).length < 4096) {
    const stream = await OpenAIStream(PromptHelper.createPayload(prompt, true));

    encoding.free();

    return new Response(stream);
  }

  // Split diff into chunks and generate prompts for each chunk
  const descriptions = await Promise.all(
    PromptHelper.split(body.diff, encoding).map(async (chunk) => {
      let chunkPrompt = PromptHelper.generatePrompt(chunk, body.template);

      if (encoding.encode(chunkPrompt).length > 4096) {
        chunkPrompt = PromptHelper.generatePrompt(chunk, body.template, true);

        // Check if minified prompt is still too long
        if (encoding.encode(chunkPrompt).length > 4096) {
          return new Response(
            `The diff is too large (${
              encoding.encode(chunkPrompt).length
            }), try reducing the number of staged changes.`,
            {
              status: 400,
            }
          );
        }
      }

      const response = await fetch(
        'https://api.openai.com/v1/chat/completions',
        {
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${process.env.OPENAI_API_KEY ?? ''}`,
          },
          method: 'POST',
          body: JSON.stringify(PromptHelper.createPayload(chunkPrompt, false)),
        }
      );

      if (!response.ok) {
        return new Response(await response.text(), {
          status: 500,
        });
      }

      const data = await response.json();

      return data.choices[0].message.content;
    })
  );

  // Ask ChatGPT to consolidate into one description
  const stream = await OpenAIStream(
    PromptHelper.createPayload(
      PromptHelper.generateConsolidatePrompt(descriptions, body.template),
      true
    )
  );

  encoding.free();

  return new Response(stream);
}
