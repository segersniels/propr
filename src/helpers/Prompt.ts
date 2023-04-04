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
function prepareDiff(diff: string, minify = false) {
  if (!minify) {
    return diff;
  }

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

export function generatePrompt(diff: string, template: string, minify = false) {
  return `
    Generate a concise PR description from the provided git diff according to a provided template.
    The PR description should be a good summary of the changes made.
    Do not reference each file and function added but rather give a general explanation of the changes made.
    You are free to make a calculated guess as to which changes and files are related to each other so you can group them together.
    When endpoints or routes are added or altered reference these when describing features.
    It's not worth mentioning that you added tests when you mention you added a new feature as it implies that you added tests in the first place.
    If notes or reason why the change happened are requested, make sure you try to explain the reasoning without using too much technical jargon.
    If the diff provided is not actually a diff I want you to respond with an appropriate message accordingly.

    The PR description should be structured as follows: """
    ${template}
    """

    Here is the diff: """
    ${prepareDiff(diff, minify)}
    """
  `;
}
