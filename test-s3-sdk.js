// Test S3 SDK download instead of HTTP
const { S3Client, GetObjectCommand } = require('@aws-sdk/client-s3');

async function testS3Sdk() {
  const accessKey = process.env.S3_ACCESS_KEY;
  const secretKey = process.env.S3_SECRET_KEY;
  const endpoint = process.env.S3_ENDPOINT || 'https://s3.ru1.storage.beget.cloud';
  const region = process.env.S3_REGION || 'ru1';
  const bucketName = process.env.S3_BUCKET_NAME;

  console.log('Bucket:', bucketName);
  console.log('Endpoint:', endpoint);
  console.log('AccessKey:', accessKey ? accessKey.substring(0, 4) + '...' : 'MISSING');

  if (!accessKey || !secretKey || !bucketName) {
    console.error('Missing S3 credentials!');
    return;
  }

  const client = new S3Client({
    region,
    endpoint,
    credentials: { accessKeyId: accessKey, secretAccessKey: secretKey },
    forcePathStyle: true,
  });

  const key = 'products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  console.log('Fetching key:', key);

  try {
    const command = new GetObjectCommand({ Bucket: bucketName, Key: key });
    const response = await client.send(command);

    console.log('ContentType:', response.ContentType);
    console.log('ContentLength:', response.ContentLength);

    // Read the stream
    const chunks = [];
    for await (const chunk of response.Body) {
      chunks.push(chunk);
    }
    const buf = Buffer.concat(chunks);
    console.log('Downloaded via S3 SDK:', buf.length, 'bytes');
    const base64 = buf.toString('base64');
    console.log('Base64 length:', base64.length);
    console.log('SUCCESS!');
  } catch (e) {
    console.error('S3 SDK error:', e.name, e.message);
  }
}

// Test 2: Also try http.get without timeout
async function testHttpNoTimeout() {
  const http = require('https');
  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
  console.log('\n--- Test: HTTPS without timeout ---');
  return new Promise((resolve) => {
    const req = http.get(url, { headers: { 'Accept-Encoding': 'identity' } }, (res) => {
      console.log('Status:', res.statusCode);
      const chunks = [];
      res.on('data', (chunk) => {
        chunks.push(chunk);
        process.stdout.write(`\rReceived: ${chunks.reduce((s, c) => s + c.length, 0)} bytes`);
      });
      res.on('end', () => {
        const buf = Buffer.concat(chunks);
        console.log('\nTotal downloaded:', buf.length, 'bytes');
        resolve(buf);
      });
      res.on('error', (e) => {
        console.error('\nStream error:', e.message);
        resolve(null);
      });
    });
    req.on('error', (e) => {
      console.error('Request error:', e.message);
      resolve(null);
    });
    // Absolute fallback timeout at 60s
    setTimeout(() => { req.destroy(); console.error('\nForce timeout at 60s'); resolve(null); }, 60000);
  });
}

(async () => {
  console.log('=== S3 SDK Test ===');
  await testS3Sdk();
  
  console.log('\n=== HTTP Direct Test ===');
  await testHttpNoTimeout();
})();
