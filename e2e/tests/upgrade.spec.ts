import { expect, test } from '@playwright/test';
import { containerLogs, OLD_IMAGE, PORTS, removeContainer, run, runContainer, startApp, stopApp } from '../support/stack';
import { download, expectThumbnailsFrom, expectVideoPlaying, stubPlayerIcons, useDarkTheme, watch } from '../support/page';

// a real upgrade: the library is written by the 0.1.17 image with a 0.1.17 compose file's environment, then the new image
// is started on the same volumes with that same, unchanged environment
test('a library written by 0.1.17 opens on the new image with the old compose file unchanged', async ({ page, browser }) => {
    test.skip(run('docker', ['image', 'inspect', OLD_IMAGE], { allowFail: true }).status !== 0, `${OLD_IMAGE} is not on this machine`);
    test.setTimeout(300_000);
    const old = 'streamsphere-e2e-old-app';
    const compose = { HOST_IP: 'localhost', APPLICATION_PORT: String(PORTS.app), CONTENT_PORT: String(PORTS.content) };

    // 1. the old image downloads a second copy of a video into the shared volumes
    stopApp();
    runContainer(old, OLD_IMAGE, PORTS.app, compose);
    try {
        // the old entrypoint updates yt-dlp before it serves anything
        await expect.poll(() => containerLogs(old), { timeout: 60_000 }).toContain('Connected to database');
        await stubPlayerIcons(page);
        await page.goto(`http://localhost:${PORTS.app}/#/home`);
        await useDarkTheme(page);
        await download(page, 'youtube.com/watch?v=jNQXAC9IVRw');
    } finally {
        removeContainer(old);
    }

    // 2. the new image, same volumes, same environment
    const started = await startApp(compose);
    expect(started.state, started.logs).toBe('running');
    expect(started.logs).toContain('HOST_IP is deprecated, using it as APPLICATION_IP');

    // a browser that has not seen the old UI
    const fresh = await browser.newContext({ ignoreHTTPSErrors: true });
    const user = await fresh.newPage();
    const seen = watch(user);
    const content = `http://localhost:${PORTS.content}`;
    await stubPlayerIcons(user);
    await user.goto(`http://localhost:${PORTS.app}/#/videos`);
    await expect.poll(() => [...seen.api]).toEqual([`http://localhost:${PORTS.app}`]);
    await expectThumbnailsFrom(user, content);
    await user.locator('p-card', { hasText: 'Me at the zoo' }).locator('img').first().click();
    await expectVideoPlaying(user, content);
    expect(seen.errors).toEqual([]);
    await fresh.close();

    // the tab that still has the old UI open: its code is already running, a hash-only navigation does not reload it, and
    // the new backend returns plain content paths the old UI cannot place. Recorded, not asserted: nothing the server sends
    // can fix an open tab, so the release notes say to reload open tabs after an upgrade
    await page.goto(`http://localhost:${PORTS.app}/#/videos`);
    await page.waitForTimeout(3000);
    const stale = await page.locator('img').evaluateAll(imgs => (imgs as HTMLImageElement[])
        .filter(i => i.currentSrc && !/\/assets\//.test(i.currentSrc)).map(i => ({ src: i.currentSrc, loaded: i.complete && i.naturalWidth > 0 })));
    test.info().annotations.push({ type: 'tab left open across the upgrade', description: stale.every(t => t.loaded) ? 'thumbnails loaded' : `thumbnails did NOT load: ${stale.map(t => t.src).join(', ')}` });
});
