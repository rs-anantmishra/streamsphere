import { expect, test } from '@playwright/test';
import { PORTS, startApp } from '../support/stack';
import { stubPlayerIcons, useDarkTheme, watch } from '../support/page';

// the queue drawer is filled from /api/download/queued-items. An item only stays "queued" for a moment, so the response is
// supplied here, in the exact shape the backend returns it (the thumbnail is a plain path, as everywhere else); the
// backend side of that shape is covered by the Go tests
test('the queue drawer draws a queued item with its thumbnail from the content address', async ({ page }) => {
    const started = await startApp({ APPLICATION_PORT: String(PORTS.app), CONTENT_PORT: String(PORTS.content) });
    expect(started.state, started.logs).toBe('running');

    const queued = [{
        video_id: 1, title: 'Me at the zoo', description: '', duration: 19, webpage_url: 'https://www.youtube.com/watch?v=jNQXAC9IVRw',
        channel: 'jawed', domain: '', video_format: '', extension: '', watch_count: 0, views_count: 0, likes_count: 0, filesize: 0,
        upload_date: '', is_deleted: false, categories: null, tags: null,
        thumbnail: '/youtube.com/jawed/Videos/Thumbnails/Me at the zoo [jNQXAC9IVRw].webp', media_url: '',
    }];
    await page.route(url => url.pathname === '/api/download/queued-items', route => {
        const state = new URL(route.request().url()).searchParams.get('state');
        return route.fulfill({ json: state === 'queued' ? queued : [] });
    });

    const seen = watch(page);
    await stubPlayerIcons(page);
    await page.goto(`http://localhost:${PORTS.app}/#/home`);
    await useDarkTheme(page);
    await page.getByRole('button', { name: 'Queued', exact: true }).click();

    const drawer = page.locator('.queued-box');
    await expect(drawer).toBeVisible();
    await expect(drawer.getByText('Me at the zoo')).toBeVisible();
    const images = drawer.locator('img');
    await expect.poll(async () => images.evaluateAll(imgs => (imgs as HTMLImageElement[]).map(i => i.currentSrc))).toContainEqual(
        expect.stringContaining(`http://localhost:${PORTS.content}/youtube.com/jawed/Videos/Thumbnails/`));
    await expect.poll(async () => images.evaluateAll(imgs => (imgs as HTMLImageElement[])
        .filter(i => i.currentSrc.includes('/Thumbnails/')).every(i => i.complete && i.naturalWidth > 0))).toBe(true);
    expect(seen.errors).toEqual([]);
});
