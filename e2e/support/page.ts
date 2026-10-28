import { expect, Page } from '@playwright/test';

// the player (Plyr) loads its icon sprite from cdn.plyr.io with a bare XMLHttpRequest on every page that has a player;
// when that request fails Angular logs "ERROR Error: 0". The tests are about addresses, not icons, so the CDN is answered
// locally and a flaky connection cannot fail them
export async function stubPlayerIcons(page: Page) {
    await page.route('https://cdn.plyr.io/**', route => route.fulfill({ contentType: 'image/svg+xml', body: '<svg xmlns="http://www.w3.org/2000/svg"></svg>' }));
}

// switches the app to its dark theme (a visible window should not blind anyone)
export async function useDarkTheme(page: Page) {
    const toggle = page.getByRole('button', { name: '#dark', exact: true });
    if (await toggle.count()) { await toggle.click() }
}

// what a page did on the network and in the console, collected as it happens
export function watch(page: Page) {
    const seen = { api: new Set<string>(), failedApi: [] as string[], sockets: [] as string[], errors: [] as string[], logs: [] as string[] };
    page.on('request', r => { const u = new URL(r.url()); if (u.pathname.startsWith('/api/')) { seen.api.add(u.origin) } });
    page.on('requestfailed', r => { if (new URL(r.url()).pathname.startsWith('/api/')) { seen.failedApi.push(r.url()) } });
    page.on('websocket', ws => seen.sockets.push(ws.url()));
    // net::ERR_NETWORK_CHANGED is Chromium reacting to this machine's network devices changing (docker creates and removes
    // them while the stack is recreated). It retries the request itself, and every test still checks that its images
    // loaded and its video played, so it is not counted as an app error
    page.on('console', m => { if (m.type() === 'error' && !/ERR_NETWORK_CHANGED/.test(m.text())) { seen.errors.push(m.text()) } });
    page.on('console', m => { if (m.type() === 'log') { seen.logs.push(m.text()) } });
    page.on('pageerror', e => seen.errors.push(String(e)));
    return seen;
}

// the thumbnails on the page (not the app's own assets): where they came from and whether they loaded
export async function thumbnails(page: Page) {
    return page.locator('img').evaluateAll(imgs => (imgs as HTMLImageElement[])
        .filter(i => i.currentSrc && !/\/assets\//.test(i.currentSrc))
        .map(i => ({ origin: new URL(i.currentSrc).origin, loaded: i.complete && i.naturalWidth > 0 })));
}

export async function expectThumbnailsFrom(page: Page, origin: string) {
    // one read per attempt: the page can re-render between two reads, so the check is made on a single snapshot
    await expect.poll(async () => {
        const shown = await thumbnails(page);
        if (shown.length === 0) { return 'no thumbnails on the page' }
        const unloaded = shown.filter(t => !t.loaded);
        if (unloaded.length) { return `${unloaded.length} of ${shown.length} thumbnails did not load`; }
        const elsewhere = shown.filter(t => t.origin !== origin);
        if (elsewhere.length) { return `thumbnails from ${[...new Set(elsewhere.map(t => t.origin))].join(', ')}, expected ${origin}` }
        return 'ok';
    }).toBe('ok');
}

// plays the video on the page for a moment
export async function expectVideoPlaying(page: Page, origin: string) {
    const video = page.locator('video').first();
    await expect(video).toBeAttached();
    const state = await video.evaluate(async (v: HTMLVideoElement) => {
        if (v.paused) { await v.play().catch(() => {}) }
        await new Promise(r => setTimeout(r, 2500));
        const out = { src: v.currentSrc, ready: v.readyState, time: v.currentTime, error: v.error?.code ?? null };
        v.pause();
        return out;
    });
    expect(new URL(state.src).origin).toBe(origin);
    expect(state.error).toBeNull();
    expect(state.ready).toBeGreaterThanOrEqual(3);
    expect(state.time).toBeGreaterThan(1);
}

export const urlField = (page: Page) => page.getByRole('textbox', { name: /videoId|playlistId|youtu/ });

// starts a real download from the home page and waits for it to finish; returns the panel thumbnail's origin
export async function download(page: Page, url: string) {
    await urlField(page).fill(url);
    await page.getByRole('button', { name: 'Get', exact: true }).click();
    await expect(page.getByText('Download completed successfully.')).toBeVisible({ timeout: 120_000 });
    const panel = page.locator('img.info-image-size');
    await expect(panel).toHaveJSProperty('complete', true);
    return new URL(await panel.evaluate((i: HTMLImageElement) => i.currentSrc)).origin;
}

// saves a screenshot to the folder in E2E_SHOTS (nothing happens when it is not set), to look at a page afterwards
export async function shot(page: Page, name: string) {
    if (process.env.E2E_SHOTS) { await page.screenshot({ path: `${process.env.E2E_SHOTS}/${name}.png` }) }
}
