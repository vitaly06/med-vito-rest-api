// Test 1: Node fetch API
async function testFetch() {
  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  console.log('--- Test: Node.js native fetch ---');
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 30000);
    const resp = await fetch(url, { signal: controller.signal });
    clearTimeout(timeout);
    console.log('Status:', resp.status);
    const arrayBuffer = await resp.arrayBuffer();
    const buf = Buffer.from(arrayBuffer);
    console.log('Size:', buf.length, 'bytes');
    console.log('Base64 length:', buf.toString('base64').length);
    return buf;
  } catch (e) {
    console.error('Fetch error:', e.message);
    return null;
  }
}

// Test 2: axios with decompression disabled
async function testAxiosNoDecomp() {
  const axios = require('axios');
  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  console.log('--- Test: axios with decompress:false ---');
  try {
    const resp = await axios.get(url, {
      responseType: 'arraybuffer',
      timeout: 30000,
      decompress: false,
      headers: { 'Accept-Encoding': 'identity' }
    });
    const buf = Buffer.from(resp.data);
    console.log('Size:', buf.length, 'bytes');
    return buf;
  } catch (e) {
    console.error('Axios error:', e.code, e.message);
    return null;
  }
}

// Test 3: axios with httpAgent settings
async function testAxiosAgent() {
  const axios = require('axios');
  const https = require('https');
  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  console.log('--- Test: axios with custom agent ---');
  try {
    const agent = new https.Agent({ keepAlive: false, timeout: 30000 });
    const resp = await axios.get(url, {
      responseType: 'arraybuffer',
      timeout: 30000,
      httpsAgent: agent,
      maxContentLength: 50 * 1024 * 1024,
      maxBodyLength: 50 * 1024 * 1024,
      headers: { 'Accept-Encoding': 'identity', 'User-Agent': 'MedVito/1.0' }
    });
    const buf = Buffer.from(resp.data);
    console.log('Size:', buf.length, 'bytes');
    return buf;
  } catch (e) {
    console.error('Axios agent error:', e.code, e.message);
    return null;
  }
}

(async () => {
  const buf = await testFetch();
  if (!buf) await testAxiosNoDecomp();
  if (!buf) await testAxiosAgent();
  if (buf) console.log('\nImage download successful!');
})();
