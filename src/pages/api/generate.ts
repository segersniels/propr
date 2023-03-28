import type { NextApiRequest, NextApiResponse } from 'next';
import { generate } from 'lib/api';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<any>
) {
  const prompt = `
    Generate a concise PR description from the provided git diff according to a provided template.
    Be thorough and concise. The PR description should be a good summary of the changes made. Don't be afraid to go into detail.
    If a section has no content, you can leave the entire section out.

    It's not worth mentioning that you added tests when you mention you added a new feature.
    It implies that you added tests in the first place.

    Here is the template: """
    ${req.body.template}
    """

    Here is the diff: """
    ${req.body.diff}
    """
  `;

  const message = await generate(prompt);

  return res.status(200).json({
    message,
  });
}
