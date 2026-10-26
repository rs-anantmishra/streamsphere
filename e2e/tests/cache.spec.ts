import { expect, test } from '@playwright/test';
import { PORTS, startApp } from '../support/stack';
import { expectThumbnailsFrom, useDarkTheme } from '../support/page';

// no request interception (page.route) in this test: Playwright turns the browser's HTTP cache off while a route is active,
// and the cache is exactly what is under test
// main-*.js is rewritten with the environment's addresses at every container start, under the same file name. The same
// browser must pick up a changed variable with an ordinary navigation, no hard reload
test('a changed variable is picked up by a browser that already has the page', async ({ page }) => {
    const base = `http://localhost:${PORTS.app}`;
    const ports = { APPLICATION_PORT: String(PORTS.app), CONTENT_PORT: String(PORTS.content) };

    expect((await startApp(ports)).state).toBe('running');
    await page.goto(`${base}/#/videos`);
    await useDarkTheme(page);
    await expectThumbnailsFrom(page, `http://localhost:${PORTS.content}`);

    // the same content server, reached by another name
    expect((await startApp({ ...ports, CONTENT_IP: '127.0.0.1' })).state).toBe('running');
    await page.goto('about:blank');
    await page.goto(`${base}/#/videos`);
    await expectThumbnailsFrom(page, `http://127.0.0.1:${PORTS.content}`);
});
