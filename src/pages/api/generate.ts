import { NextRequest } from 'next/server';
import * as PromptHelper from 'helpers/Prompt';
import * as OpenAIHelper from 'helpers/OpenAI';
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

/**
 * Split diff into chunks and generate prompts for each chunk
 */
async function consolidateUsingChunks(
  encoding: Tiktoken,
  diff: string,
  template: string
) {
  const maxTokenLength = PromptHelper.getMaxTokenLength({
    encoding,
  });

  return await Promise.all(
    PromptHelper.split({
      diff,
      template,
      encoding,
    }).map(async (chunk) => {
      let chunkPrompt = PromptHelper.generatePrompt(chunk, template);

      // Check if minified prompt is still too long
      if (encoding.encode(chunkPrompt).length > maxTokenLength) {
        return new Response(
          `The diff is too large (${
            encoding.encode(chunkPrompt).length
          }), try reducing the number of staged changes.`,
          {
            status: 400,
          }
        );
      }

      const response = await fetch(
        'https://api.openai.com/v1/chat/completions',
        {
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${process.env.OPENAI_API_KEY ?? ''}`,
          },
          method: 'POST',
          body: JSON.stringify(OpenAIHelper.createPayload(chunkPrompt, false)),
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
}

export default async function handler(req: NextRequest) {
  await init((imports) => WebAssembly.instantiate(wasm, imports));
  const encoding = new Tiktoken(
    model.bpe_ranks,
    model.special_tokens,
    model.pat_str
  );

  const body = await req.json();
  let maxTokenLength = PromptHelper.getMaxTokenLength({
    encoding,
  });

  // Within model token length so pass entire diff
  let prompt = PromptHelper.generatePrompt(body.diff, body.template);
  if (encoding.encode(prompt).length < maxTokenLength) {
    const stream = await OpenAIHelper.createOpenAIStream(
      OpenAIHelper.createPayload(prompt, true)
    );

    encoding.free();

    return new Response(stream);
  }

  // Check whether the minified body is within model token length
  prompt = PromptHelper.generatePrompt(body.diff, body.template, true);
  if (encoding.encode(prompt).length < maxTokenLength) {
    const stream = await OpenAIHelper.createOpenAIStream(
      OpenAIHelper.createPayload(prompt, true)
    );

    encoding.free();

    return new Response(stream);
  }

  // Ask ChatGPT to consolidate chunks into one description
  const descriptions = await consolidateUsingChunks(
    encoding,
    body.diff,
    body.template
  );

  const stream = await OpenAIHelper.createOpenAIStream(
    OpenAIHelper.createPayload(
      PromptHelper.generateConsolidatePrompt(descriptions, body.template),
      true
    )
  );

  encoding.free();

  return new Response(stream);
}
