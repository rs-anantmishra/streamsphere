import { defineConfig } from '@playwright/test';
import * as fs from 'node:fs';
import { detectPlatformOverride } from './support/platform.ts';

// see support/platform.ts: only LMDE needs telling what it really is
const osRelease = fs.existsSync('/etc/os-release') ? fs.readFileSync('/etc/os-release', 'utf8') : undefined;
const platform = detectPlatformOverride(osRelease, process.arch);
if (platform) { process.env.PLAYWRIGHT_HOST_PLATFORM_OVERRIDE ??= platform }

// the visible window is a rule for these tests, never headless.
// The stack (docker, nginx, certificate) is built in support/stack.ts; needs the image: see the top of that file
export default defineConfig({
    testDir: './tests',
    globalSetup: './support/global-setup.ts',
    workers: 1,
    fullyParallel: false,
    timeout: 120_000,
    expect: { timeout: 15_000 },
    reporter: [['list']],
    use: {
        headless: false,
        colorScheme: 'dark',
        ignoreHTTPSErrors: true,
    },
    projects: [
        { name: 'seed', testMatch: /seed\.spec\.ts/ },
        { name: 'scenarios', testMatch: /(scenarios|old-image|queue-drawer|unsubstituted|upgrade|cache)\.spec\.ts/, dependencies: ['seed'] },
        { name: 'install', testMatch: /install\.spec\.ts/ },
    ],
});
