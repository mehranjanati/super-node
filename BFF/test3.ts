import { convertToModelMessages } from 'ai';

try {
  const msgs = [{ role: 'user', content: 'hello' }];
  console.log(convertToModelMessages(msgs as any));
} catch (e) {
  console.error(e);
}
