const fs = require('fs');

async function main() {
  const log = [];
  const l = (m) => log.push(m);

  const url = 'https://s3.ru1.storage.beget.cloud/c15b4d655f70-medvito-data/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png';

  try {
    l('=== Node.js fetch download test ===');
    l('URL: ' + url);
    const controller = new AbortController();
    const t = setTimeout(() => { controller.abort(); l('TIMEOUT!'); }, 20000);
    const r = await fetch(url, { signal: controller.signal });
    clearTimeout(t);
    l('Status: ' + r.status);
    l('Content-Length: ' + r.headers.get('content-length'));
    const ab = await r.arrayBuffer();
    l('Downloaded: ' + ab.byteLength + ' bytes');
    l('SUCCESS');
  } catch (e) {
    l('FAIL: ' + e.message);
  }

  fs.writeFileSync('/tmp/host-net-test.txt', log.join('\n'));
}

main().then(() => process.exit(0)).catch(e => {
  require('fs').writeFileSync('/tmp/host-net-test.txt', 'FATAL: ' + e.message);
  process.exit(1);
});
