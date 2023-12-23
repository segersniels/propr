import { experimental_AssistantResponse } from 'ai';
import OpenAI from 'openai';
import { MessageContentText } from 'openai/resources/beta/threads/messages/messages';

const FILES_TO_IGNORE = [
  'package-lock.json',
  'yarn.lock',
  'npm-debug.log',
  'yarn-debug.log',
  'yarn-error.log',
  '.pnpm-debug.log',
  'Cargo.lock',
  'Gemfile.lock',
  'mix.lock',
  'Pipfile.lock',
  'composer.lock',
  'glide.lock',
];

const openai = new OpenAI({
  apiKey: process.env.OPENAI_API_KEY || '',
});

export const runtime = 'edge';

function splitDiffIntoChunks(diff: string) {
  return Array.from(
    diff.matchAll(/diff --git[\s\S]*?(?=diff --git|$)/g),
    (match) => match[0]
  ).map((chunk) => chunk.replace(/ {2,}/g, ''));
}

function removeLockFiles(chunks: string[]) {
  return chunks.filter((chunk) => {
    const firstLine = chunk.split('\n')[0];

    for (const file of FILES_TO_IGNORE) {
      if (firstLine.includes(file)) {
        return false;
      }
    }

    return true;
  });
}

/**
 * Prepare a diff for use in the prompt by removing stuff like
 * the lockfile changes and removing some of the whitespace.
 */
function prepareDiff(diff: string) {
  const chunks = splitDiffIntoChunks(diff);

  return removeLockFiles(chunks).join('\n');
}

function generateUserMessage(diff: string, template: string) {
  return `
    The diff:
    """
    ${prepareDiff(diff)}
    """

    Use the following template to write your description:
    """
    ${template}
    """
  `;
}

export async function POST(req: Request) {
  const input: {
    threadId: string | null;
    message: string;
    template: string;
  } = await req.json();

  const threadId = input.threadId ?? (await openai.beta.threads.create({})).id;
  const createdMessage = await openai.beta.threads.messages.create(threadId, {
    role: 'user',
    content: generateUserMessage(input.message, input.template),
  });

  return experimental_AssistantResponse(
    { threadId, messageId: createdMessage.id },
    async ({ threadId, sendMessage }) => {
      const run = await openai.beta.threads.runs.create(threadId, {
        assistant_id:
          process.env.ASSISTANT_ID ??
          (() => {
            throw new Error('ASSISTANT_ID is not set');
          })(),
      });

      async function waitForRun(run: OpenAI.Beta.Threads.Runs.Run) {
        while (run.status === 'queued' || run.status === 'in_progress') {
          await new Promise((resolve) => setTimeout(resolve, 500));

          run = await openai.beta.threads.runs.retrieve(threadId!, run.id);
        }

        if (
          run.status === 'cancelled' ||
          run.status === 'cancelling' ||
          run.status === 'failed' ||
          run.status === 'expired'
        ) {
          throw new Error(run.status);
        }
      }

      await waitForRun(run);

      // Get new thread messages (after our message)
      const responseMessages = (
        await openai.beta.threads.messages.list(threadId, {
          after: createdMessage.id,
          order: 'asc',
        })
      ).data;

      // Send the messages
      for (const message of responseMessages) {
        sendMessage({
          id: message.id,
          role: 'assistant',
          content: message.content.filter(
            (content) => content.type === 'text'
          ) as Array<MessageContentText>,
        });
      }
    }
  );
}
