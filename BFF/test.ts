import { streamText } from 'ai';
import { google } from '@ai-sdk/google';

async function main() {
  try {
    const result = streamText({
      model: google('gemini-2.5-flash'),
      messages: [{ role: 'user', content: 'hello' }],
    });
    console.log("Success");
  } catch (e) {
    console.error(e);
  }
}
main();
