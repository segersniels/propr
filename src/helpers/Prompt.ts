import { MAX_RESPONSE_LENGTH } from './OpenAI';
import { Tiktoken } from '@dqbd/tiktoken/lite/init';

const GPT4_TOKEN_LENGTH = 8192;
const GPT3_TOKEN_LENGTH = 4096;
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
function prepareDiff(diff: string, minify = false) {
  const chunks = splitDiffIntoChunks(diff);

  return removeLockFiles(chunks)
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

type GetMaxTokenLengthOptions = {
  model: 'gpt-4' | 'gpt-3.5-turbo';
  encoding: Tiktoken;
} & (
  | {
      template: string;
      excludePrompt: true;
    }
  | {
      excludePrompt?: false;
      template?: undefined;
    }
);

/**
 * Get max token length of model
 */
export function getMaxTokenLength(options: GetMaxTokenLengthOptions) {
  const tokenLength =
    options.model === 'gpt-4' ? GPT4_TOKEN_LENGTH : GPT3_TOKEN_LENGTH;
  const promptLength = options.excludePrompt
    ? getDefaultPromptTokenLength(options.template, options.encoding)
    : 0;

  return tokenLength - MAX_RESPONSE_LENGTH - promptLength;
}

interface SplitOptions {
  diff: string;
  template: string;
  encoding: Tiktoken;
  model: 'gpt-4' | 'gpt-3.5-turbo';
}

/**
 * Split the large diff into separate chunks
 */
export function split(options: SplitOptions) {
  const combinedChunks = [];
  const maxTokenLength = getMaxTokenLength({
    ...options,
    excludePrompt: true,
  });

  // Split into smaller chunks
  let chunks = splitDiffIntoChunks(options.diff);

  // Filter out lockfiles
  chunks = removeLockFiles(chunks);

  /**
   * Add chunks together as long as they do not exceed token length
   * to limit the number of requests we do to the OpenAI API
   */
  let currentChunk = '';
  for (const chunk of chunks) {
    const currentChunkLength = options.encoding.encode(currentChunk).length;
    const chunkLength = options.encoding.encode(chunk).length;

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
