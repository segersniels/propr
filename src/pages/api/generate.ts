import { NextRequest, NextResponse } from 'next/server';
import { generate } from 'lib/api';

export const config = {
  runtime: 'edge',
};

export default async function handler(req: NextRequest) {
  const body = await req.json();
  const prompt = `
    Generate a concise PR description from the provided git diff according to a provided template.
    Be thorough and concise. The PR description should be a good summary of the changes made. Don't be afraid to go into detail.
    If a section has no content, you can leave the entire section out.

    It's not worth mentioning that you added tests when you mention you added a new feature.
    It implies that you added tests in the first place.

    Here is the template: """
    ${body.template}
    """

    Here is the diff: """
    ${body.diff}
    """
  `;

  const message = await generate(prompt);

  return NextResponse.json({
    message,
  });
}
