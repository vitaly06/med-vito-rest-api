const axios = require('axios');
const fs = require('fs');

async function test() {
  const log = [];
  const l = (msg) => { log.push(String(msg)); };

  const folderId = process.env.YANDEX_FOLDER_ID;
  const apiKey = process.env.YANDEX_API_KEY;
  const imageUrl = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';

  l('FolderId: ' + folderId);
  l('ApiKey: ' + (apiKey ? apiKey.substring(0, 6) + '...' : 'MISSING'));

  // Test 1: OpenAI-compatible API with Gemma 3 27B and image URL (not base64)
  try {
    l('\n=== Test 1: Gemma 3 27B with image URL ===');
    const resp = await axios.post(
      'https://ai.api.cloud.yandex.net/v1/chat/completions',
      {
        model: `gpt://${folderId}/gemma-3-27b-it/latest`,
        messages: [
          {
            role: 'user',
            content: [
              { type: 'image_url', image_url: { url: imageUrl } },
              { type: 'text', text: 'Опиши что на этом изображении в одном предложении.' },
            ],
          },
        ],
        max_tokens: 150,
        temperature: 0.1,
      },
      {
        headers: {
          Authorization: `Api-Key ${apiKey}`,
          'Content-Type': 'application/json',
          'x-folder-id': folderId,
        },
        timeout: 60000,
      },
    );
    l('Status: ' + resp.status);
    l('Response: ' + JSON.stringify(resp.data, null, 2));
  } catch (e) {
    l('Error: ' + (e.response?.status || e.code || 'unknown'));
    l('Message: ' + e.message);
    if (e.response?.data) {
      const d = typeof e.response.data === 'string' ? e.response.data : JSON.stringify(e.response.data);
      l('Body: ' + d);
    }
  }

  // Test 2: OpenAI-compatible API with YandexGPT Pro and image URL
  try {
    l('\n=== Test 2: YandexGPT Pro with image URL ===');
    const resp = await axios.post(
      'https://ai.api.cloud.yandex.net/v1/chat/completions',
      {
        model: `gpt://${folderId}/yandexgpt/latest`,
        messages: [
          {
            role: 'user',
            content: [
              { type: 'image_url', image_url: { url: imageUrl } },
              { type: 'text', text: 'Опиши что на этом изображении в одном предложении.' },
            ],
          },
        ],
        max_tokens: 150,
        temperature: 0.1,
      },
      {
        headers: {
          Authorization: `Api-Key ${apiKey}`,
          'Content-Type': 'application/json',
          'x-folder-id': folderId,
        },
        timeout: 60000,
      },
    );
    l('Status: ' + resp.status);
    l('Response: ' + JSON.stringify(resp.data, null, 2));
  } catch (e) {
    l('Error: ' + (e.response?.status || e.code || 'unknown'));
    l('Message: ' + e.message);
    if (e.response?.data) {
      const d = typeof e.response.data === 'string' ? e.response.data : JSON.stringify(e.response.data);
      l('Body: ' + d);
    }
  }

  fs.writeFileSync('/tmp/vision-test.txt', log.join('\n'));
}

test().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/vision-test.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
