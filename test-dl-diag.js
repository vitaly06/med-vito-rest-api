const fs = require('fs');

async function main() {
  const log = [];
  const l = (m) => log.push(m);

  // Test 1: Download from httpbin (different server, ~200KB)
  try {
    l('=== Test 1: httpbin 200KB ===');
    const r = await fetch('https://httpbin.org/bytes/204800');
    const ab = await r.arrayBuffer();
    l('httpbin OK: ' + ab.byteLength + ' bytes');
  } catch (e) {
    l('httpbin FAIL: ' + e.message);
  }

  // Test 2: Download our S3 image using Node.js fetch with TLS options
  try {
    l('\n=== Test 2: S3 image via fetch ===');
    const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
    const controller = new AbortController();
    const t = setTimeout(() => controller.abort(), 30000);
    const r = await fetch(url, { signal: controller.signal });
    clearTimeout(t);
    l('S3 fetch status: ' + r.status);
    const ab = await r.arrayBuffer();
    l('S3 fetch OK: ' + ab.byteLength + ' bytes');
  } catch (e) {
    l('S3 fetch FAIL: ' + e.message);
  }

  // Test 3: Try S3 HTTP (not HTTPS)
  try {
    l('\n=== Test 3: S3 image via HTTP (not HTTPS) ===');
    const url = 'http://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
    const controller = new AbortController();
    const t = setTimeout(() => controller.abort(), 30000);
    const r = await fetch(url, { signal: controller.signal, redirect: 'follow' });
    clearTimeout(t);
    l('HTTP fetch status: ' + r.status);
    const ab = await r.arrayBuffer();
    l('HTTP fetch OK: ' + ab.byteLength + ' bytes');
  } catch (e) {
    l('HTTP fetch FAIL: ' + e.message);
  }

  // Test 4: Try S3 via Node.js http.get with explicit stream handling
  try {
    l('\n=== Test 4: S3 via http module with explicit buffer ===');
    const http = require('http');
    const result = await new Promise((resolve, reject) => {
      const url = 'http://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';
      http.get(url, (res) => {
        l('HTTP status: ' + res.statusCode);
        if (res.statusCode >= 300 && res.statusCode < 400) {
          l('Redirect to: ' + res.headers.location);
          resolve(null);
          return;
        }
        const chunks = [];
        res.on('data', (ch) => chunks.push(ch));
        res.on('end', () => resolve(Buffer.concat(chunks)));
        res.on('error', reject);
      }).on('error', reject);
      setTimeout(() => reject(new Error('manual timeout 30s')), 30000);
    });
    if (result) l('HTTP module OK: ' + result.length + ' bytes');
    else l('HTTP module: got redirect');
  } catch (e) {
    l('HTTP module FAIL: ' + e.message);
  }

  fs.writeFileSync('/tmp/dl-diag.txt', log.join('\n'));
}

main().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/dl-diag.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
