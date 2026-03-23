const { execSync } = require('child_process');
const axios = require('axios');
const fs = require('fs');

async function test() {
  const log = [];
  const l = (msg) => { log.push(String(msg)); };

  const imageUrl = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  const folderId = process.env.YANDEX_FOLDER_ID;
  const apiKey = process.env.YANDEX_API_KEY;

  // Test 1: Download via wget
  try {
    l('=== Downloading via wget child_process ===');
    const buf = execSync(`wget -q -O - "${imageUrl}"`, {
      maxBuffer: 50 * 1024 * 1024,
      timeout: 60000,
    });
    l('wget download OK, size: ' + buf.length + ' bytes');
    const base64 = buf.toString('base64');
    l('Base64 length: ' + base64.length);

    // Now test API call with the base64 image
    l('\n=== Sending to YandexGPT Pro via OpenAI-compatible API ===');
    const mimeType = 'image/png';
    const resp = await axios.post(
      'https://ai.api.cloud.yandex.net/v1/chat/completions',
      {
        model: `gpt://${folderId}/yandexgpt/latest`,
        messages: [
          {
            role: 'user',
            content: [
              { type: 'image_url', image_url: { url: `data:${mimeType};base64,${base64}` } },
              { type: 'text', text: 'Опиши что на этом изображении в одном предложении на русском.' },
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
        timeout: 90000,
      },
    );
    l('API Status: ' + resp.status);
    const text = resp.data?.choices?.[0]?.message?.content;
    l('AI response: ' + text);
    l('Usage: ' + JSON.stringify(resp.data?.usage));
    l('\nSUCCESS!');
  } catch (e) {
    l('Error: ' + (e.response?.status || e.code || e.status || 'unknown'));
    l('Message: ' + e.message);
    if (e.response?.data) {
      const d = typeof e.response.data === 'string' ? e.response.data : JSON.stringify(e.response.data);
      l('Body: ' + d.substring(0, 500));
    }
  }

  fs.writeFileSync('/tmp/wget-test.txt', log.join('\n'));
}

test().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/wget-test.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
