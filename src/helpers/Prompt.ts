import { OpenAIStreamPayload } from './Stream';
import { Tiktoken } from '@dqbd/tiktoken/lite/init';

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

/**
 * Removes lines from the diff that don't start with a special character
 */
function removeExcessiveLinesFromChunk(diff: string, minify = false) {
  if (!minify) {
    return diff;
  }

  return diff
    .split('\n')
    .filter((line) => /^\W/.test(line))
    .join('\n');
}

/**
 * Prepare a diff for use in the prompt by removing stuff like
 * the lockfile changes and removing some of the whitespace.
 */
function prepareDiff(diff: string, minify = false) {
  const chunks = Array.from(
    diff.matchAll(/diff --git[\s\S]*?(?=diff --git|$)/g),
    (match) => match[0]
  ).map((chunk) => chunk.replace(/ {2,}/g, ''));

  return chunks
    .filter((chunk) => {
      const firstLine = chunk.split('\n')[0];

      for (const file of FILES_TO_IGNORE) {
        if (firstLine.includes(file)) {
          return false;
        }
      }

      return true;
    })
    .map((chunk) => removeExcessiveLinesFromChunk(chunk, minify))
    .join('\n');
}

export function generatePrompt(diff: string, template: string, minify = false) {
  return `
    Generate a concise PR description from the provided git diff according to a provided template.
    The PR description should be a good summary of the changes made.
    Do not reference each file and function added but rather give a general explanation of the changes made.
    Do not treat imports and requires as changes or new features.
    You can leave out the entire heading of the template if no applicable changes are found.

    The PR description should be structured as follows: """
    ${template}
    """

    Here is the diff: """
    ${prepareDiff(diff, minify)}
    """
  `;
}

export function generateConsolidatePrompt(
  descriptions: string[],
  template: string
) {
  return `
    Consolidate the following PR descriptions. Keep it concise.
    Respect the template as follows: """
    ${template}
    """

    Here are the diffs: """
    ${descriptions.join('\n\n')}
    """
  `;
}

/**
 * Get token length of default prompt excluding the diff
 */
function getDefaultPromptTokenLength(template: string, encoding: Tiktoken) {
  return encoding.encode(generatePrompt('', template)).length;
}

/**
 * Get max token length of model excluding the default prompt
 */
function getMaxTokenLength(template: string, encoding: Tiktoken) {
  return 4096 - getDefaultPromptTokenLength(template, encoding);
}

/**
 * Split the large diff into separate chunks
 */
export function split(diff: string, template: string, encoding: Tiktoken) {
  const combinedChunks = [];
  const maxTokenLength = getMaxTokenLength(template, encoding);

  // Split into smaller chunks
  const chunks = Array.from(
    diff.matchAll(/diff --git[\s\S]*?(?=diff --git|$)/g),
    (match) => match[0]
  );

  /**
   * Add chunks together as long as they do not exceed token length
   * to limit the number of requests we do to the OpenAI API
   */
  let currentChunk = '';
  for (const chunk of chunks) {
    const currentChunkLength = encoding.encode(currentChunk).length;
    const chunkLength = encoding.encode(chunk).length;

    if (currentChunkLength + chunkLength <= maxTokenLength) {
      currentChunk += chunk;
    } else {
      combinedChunks.push(currentChunk);
      currentChunk = chunk;
    }
  }

  // Add any remaining chunk to the array
  if (currentChunk) {
    combinedChunks.push(currentChunk);
  }

  return combinedChunks;
}

/**
 * Create payload to be sent to OpenAI API
 */
export function createPayload(
  content: string,
  stream = false
): OpenAIStreamPayload {
  return {
    model: 'gpt-3.5-turbo',
    messages: [{ role: 'user', content }],
    temperature: 0.7,
    top_p: 1,
    frequency_penalty: 0,
    presence_penalty: 0,
    max_tokens: 500,
    stream,
    n: 1,
  };
}
