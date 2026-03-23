const https = require('https');
const http = require('http');

const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';

function download(imageUrl) {
  return new Promise((resolve, reject) => {
    const proto = imageUrl.startsWith('https') ? https : http;
    const req = proto.get(imageUrl, { 
      timeout: 30000,
      headers: { 'User-Agent': 'MedVito-ModerationWorker/1.0' }
    }, (res) => {
      console.log('HTTP status:', res.statusCode);
      console.log('Headers:', JSON.stringify(res.headers, null, 2));
      
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        console.log('Redirecting to:', res.headers.location);
        download(res.headers.location).then(resolve).catch(reject);
        return;
      }
      
      const chunks = [];
      res.on('data', (chunk) => chunks.push(chunk));
      res.on('end', () => {
        const buf = Buffer.concat(chunks);
        console.log('Downloaded size:', buf.length, 'bytes');
        resolve(buf);
      });
      res.on('error', (e) => {
        console.error('Stream error:', e.message);
        reject(e);
      });
    });
    
    req.on('error', (e) => {
      console.error('Request error:', e.message);
      reject(e);
    });
    
    req.on('timeout', () => {
      console.error('Request timeout!');
      req.destroy();
      reject(new Error('timeout'));
    });
  });
}

download(url)
  .then((buf) => {
    console.log('Success! Size:', buf.length);
    const base64 = buf.toString('base64');
    console.log('Base64 length:', base64.length);
  })
  .catch((e) => console.error('Failed:', e.message));
