const { S3Client, GetObjectCommand } = require('@aws-sdk/client-s3');
const fs = require('fs');

async function main() {
  const log = [];
  const l = (msg) => { log.push(msg); };

  try {
    // Test 1: S3 SDK download
    l('=== S3 SDK Test ===');
    const accessKey = process.env.S3_ACCESS_KEY;
    const secretKey = process.env.S3_SECRET_KEY;
    const endpoint = process.env.S3_ENDPOINT || 'https://s3.ru1.storage.beget.cloud';
    const region = process.env.S3_REGION || 'ru1';
    const bucketName = process.env.S3_BUCKET_NAME;
    l('Bucket: ' + bucketName);
    l('AccessKey: ' + (accessKey ? accessKey.substring(0, 4) + '...' : 'MISSING'));

    const client = new S3Client({
      region, endpoint,
      credentials: { accessKeyId: accessKey, secretAccessKey: secretKey },
      forcePathStyle: true,
    });
    
    const key = 'products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
    l('Fetching key: ' + key);
    
    const command = new GetObjectCommand({ Bucket: bucketName, Key: key });
    const response = await client.send(command);
    l('ContentType: ' + response.ContentType);
    l('ContentLength: ' + response.ContentLength);
    
    const stream = response.Body;
    const chunks = [];
    let received = 0;
    
    // Use stream.transformToByteArray() instead of manual iteration
    const bytes = await stream.transformToByteArray();
    const buf = Buffer.from(bytes);
    l('Downloaded: ' + buf.length + ' bytes');
    l('Base64 length: ' + buf.toString('base64').length);
    l('S3_SDK_SUCCESS');
  } catch (e) {
    l('S3 SDK Error: ' + e.name + ' ' + e.message);
  }
  
  try {
    // Test 2: fetch API
    l('\n=== Fetch Test ===');
    const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 45000);
    const r = await fetch(url, { signal: controller.signal });
    clearTimeout(timeout);
    l('Fetch status: ' + r.status);
    const ab = await r.arrayBuffer();
    l('Fetch downloaded: ' + ab.byteLength + ' bytes');
    l('FETCH_SUCCESS');
  } catch (e) {
    l('Fetch Error: ' + e.message);
  }
  
  // Write results to file
  fs.writeFileSync('/tmp/test-results.txt', log.join('\n'));
}

main().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/test-results.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
