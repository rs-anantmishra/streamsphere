import { expect, test } from '@playwright/test';
import { OLD_IMAGE, run } from '../support/stack';

// the release notes must say a new-style compose file needs the new image: the old entrypoint has no HOST_IP to build
// the address from, and starts anyway with an address that has no host
test('a new-style compose on the 0.1.17 image builds an address with no host', async () => {
    test.skip(run('docker', ['image', 'inspect', OLD_IMAGE], { allowFail: true }).status !== 0, `${OLD_IMAGE} is not on this machine`);
    run('docker', ['rm', '-f', 'streamsphere-e2e-old'], { allowFail: true });
    run('docker', ['run', '-d', '--name', 'streamsphere-e2e-old', '-e', 'APPLICATION_PORT=1282', '-e', 'CONTENT_PORT=1288', OLD_IMAGE]);
    try {
        await expect.poll(() => run('docker', ['logs', 'streamsphere-e2e-old'], { allowFail: true }).out, { timeout: 30_000 }).toContain('http://:1282');
    } finally {
        run('docker', ['rm', '-f', 'streamsphere-e2e-old'], { allowFail: true });
    }
});
