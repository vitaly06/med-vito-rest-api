const axios = require('axios');

async function test() {
  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  
  try {
    // Test download
    console.log('Downloading image...');
    const resp = await axios.get(url, { 
      responseType: 'arraybuffer', 
      timeout: 30000,
      headers: { 'User-Agent': 'MedVito-Test/1.0' }
    });
    const buf = Buffer.from(resp.data);
    console.log('Downloaded OK, size:', buf.length, 'bytes');
    console.log('Content-Type:', resp.headers['content-type']);
    
    const base64 = buf.toString('base64');
    console.log('Base64 length:', base64.length);
    
    // Test API call
    const folderId = process.env.YANDEX_FOLDER_ID;
    const apiKey = process.env.YANDEX_API_KEY;
    
    if (!folderId || !apiKey) {
      console.log('No YANDEX_FOLDER_ID or YANDEX_API_KEY, skipping API test');
      return;
    }
    
    console.log('Sending to YandexGPT vision...');
    const mimeType = resp.headers['content-type'] || 'image/jpeg';
    
    const apiResp = await axios.post(
      'https://llm.api.cloud.yandex.net/foundationModels/v1/completion',
      {
        modelUri: `gpt://${folderId}/yandexgpt-pro/latest`,
        completionOptions: { stream: false, temperature: 0.05, maxTokens: 150 },
        messages: [
          {
            role: 'user',
            content: [
              { type: 'image_url', image_url: { url: `data:${mimeType};base64,${base64}` } },
              { type: 'text', text: 'What do you see in this image? Answer briefly in Russian.' },
            ],
          },
        ],
      },
      {
        headers: { Authorization: `Api-Key ${apiKey}`, 'Content-Type': 'application/json' },
        timeout: 60000,
      },
    );
    
    console.log('API Response status:', apiResp.status);
    console.log('API Response:', JSON.stringify(apiResp.data, null, 2));
  } catch (e) {
    console.error('Error:', e.code || e.response?.status || 'unknown');
    console.error('Message:', e.message);
    if (e.response?.data) {
      const data = e.response.data;
      if (Buffer.isBuffer(data)) {
        console.error('Response body:', data.toString('utf-8'));
      } else {
        console.error('Response body:', JSON.stringify(data));
      }
    }
  }
}

test();
