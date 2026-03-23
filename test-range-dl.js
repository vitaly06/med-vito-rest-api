const { S3Client, GetObjectCommand, HeadObjectCommand } = require('@aws-sdk/client-s3');
const fs = require('fs');

const CHUNK_SIZE = 8192; // 8KB chunks

async function main() {
  const log = [];
  const l = (m) => log.push(m);

  const accessKey = process.env.S3_ACCESS_KEY;
  const secretKey = process.env.S3_SECRET_KEY;
  const endpoint = process.env.S3_ENDPOINT || 'https://s3.ru1.storage.beget.cloud';
  const region = process.env.S3_REGION || 'ru1';
  const bucketName = process.env.S3_BUCKET_NAME;
  const key = 'products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';

  const client = new S3Client({
    region, endpoint,
    credentials: { accessKeyId: accessKey, secretAccessKey: secretKey },
    forcePathStyle: true,
    requestHandler: {
      requestTimeout: 10000,
      httpsAgent: { maxSockets: 1, keepAlive: false },
    },
  });

  try {
    // Step 1: Get file size
    l('=== Getting file size ===');
    const head = await client.send(new HeadObjectCommand({ Bucket: bucketName, Key: key }));
    const totalSize = head.ContentLength;
    const contentType = head.ContentType;
    l('Size: ' + totalSize + ', Type: ' + contentType);

    // Step 2: Download in chunks using Range requests
    l('\n=== Downloading in chunks ===');
    const chunks = [];
    let offset = 0;
    let chunkNum = 0;

    while (offset < totalSize) {
      const end = Math.min(offset + CHUNK_SIZE - 1, totalSize - 1);
      chunkNum++;

      try {
        const resp = await client.send(new GetObjectCommand({
          Bucket: bucketName,
          Key: key,
          Range: `bytes=${offset}-${end}`,
        }));

        const bytes = await resp.Body.transformToByteArray();
        chunks.push(Buffer.from(bytes));
        offset += bytes.length;

        if (chunkNum % 20 === 0) {
          l(`Chunk ${chunkNum}: ${offset}/${totalSize} bytes (${Math.round(offset/totalSize*100)}%)`);
        }
      } catch (e) {
        l(`Chunk ${chunkNum} (${offset}-${end}) FAILED: ${e.message}`);
        // Retry once
        try {
          const resp2 = await client.send(new GetObjectCommand({
            Bucket: bucketName,
            Key: key,
            Range: `bytes=${offset}-${end}`,
          }));
          const bytes = await resp2.Body.transformToByteArray();
          chunks.push(Buffer.from(bytes));
          offset += bytes.length;
        } catch (e2) {
          l(`Retry also failed: ${e2.message}`);
          break;
        }
      }
    }

    const buf = Buffer.concat(chunks);
    l('\nTotal downloaded: ' + buf.length + ' bytes');
    l('Expected: ' + totalSize + ' bytes');
    l('Match: ' + (buf.length === totalSize));

    if (buf.length === totalSize) {
      const base64 = buf.toString('base64');
      l('Base64 length: ' + base64.length);
      l('SUCCESS!');
    }
  } catch (e) {
    l('ERROR: ' + e.name + ' ' + e.message);
  }

  fs.writeFileSync('/tmp/range-test.txt', log.join('\n'));
}

main().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/range-test.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
