import { convertToCoreMessages } from 'ai';

try {
  const msgs = [{ role: 'user', content: 'hello' }];
  console.log(convertToCoreMessages(msgs as any));
} catch (e) {
  console.error(e);
}
