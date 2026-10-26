import { expect, test } from '@playwright/test';
import { containerLogs, IMAGE, PORTS, removeContainer, runContainer } from '../support/stack';
import { expectThumbnailsFrom, stubPlayerIcons, thumbnails, useDarkTheme, watch } from '../support/page';

// the built UI served without entrypoint.sh: its ${...} placeholders are never replaced. The app must not break: every
// address falls back to the address the page was opened at
test('a build whose placeholders were never substituted uses the page address for everything', async ({ page }) => {
    const name = 'streamsphere-e2e-bare';
    const origin = `http://localhost:${PORTS.bare}`;
    runContainer(name, IMAGE, PORTS.bare, {}, ['sh', '-c', 'cd /app/application && chmod +x streamsphere && ./streamsphere']);
    try {
        await expect.poll(() => containerLogs(name), { timeout: 30_000 }).toContain('Connected to database');

        const seen = watch(page);
        await stubPlayerIcons(page);
        await page.goto(`${origin}/#/videos`);
        await useDarkTheme(page);

        await expect.poll(() => [...seen.api]).toEqual([origin]);
        expect(seen.failedApi).toEqual([]);
        // the library is there (the stack's volumes are mounted); the content is looked for at the page address, where
        // nothing serves it, because without the entrypoint nothing says where the content server is
        await expect.poll(async () => (await thumbnails(page)).length).toBeGreaterThan(0);
        expect((await thumbnails(page)).every(t => t.origin === origin)).toBe(true);
        expect(seen.errors.filter(e => !/Failed to load resource: the server responded with a status of 404/.test(e))).toEqual([]);
    } finally {
        removeContainer(name);
    }
});
