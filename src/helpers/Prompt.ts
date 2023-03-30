const FILES_TO_IGNORE = [
  'package-lock.json',
  'yarn.lock',
  'npm-debug.log',
  'yarn-debug.log',
  'yarn-error.log',
  '.pnpm-debug.log',
];

/**
 * Removes lines from the diff that don't start with a special character
 */
function removeExcessiveLinesFromChunk(diff: string) {
  return diff
    .split('\n')
    .filter((line) => /^\W/.test(line))
    .join('\n');
}

/**
 * Prepare a diff for use in the prompt by removing stuff like
 * the lockfile changes and removing some of the whitespace.
 */
export function prepareDiff(diff: string) {
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
    .map(removeExcessiveLinesFromChunk)
    .join('\n');
}
