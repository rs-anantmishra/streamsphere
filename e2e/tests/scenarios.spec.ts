import { expect, test } from '@playwright/test';
import { scenarios, pageUrl } from '../support/scenarios';
import { startApp } from '../support/stack';
import { download, expectThumbnailsFrom, expectVideoPlaying, shot, stubPlayerIcons, thumbnails, useDarkTheme, watch } from '../support/page';

for (const s of scenarios) {
    test(s.id, async ({ page }) => {
        test.setTimeout(s.download ? 240_000 : 120_000);
        const started = await startApp(s.env);

        if (s.stops) {
            expect(started.state).toBe('exited');
            expect(started.exitCode).toBe(1);
            expect(started.logs).toContain(`${s.stops} is not set`);
            return;
        }
        expect(started.state, started.logs).toBe('running');
        for (const line of s.prints ?? []) { expect(started.logs).toContain(line) }

        await stubPlayerIcons(page);
        const seen = watch(page);
        const base = pageUrl(s);
        await page.goto(`${base}/#/home`);
        await useDarkTheme(page);

        if (s.apiUnreachable) {
            // the library and the playlists say the api cannot be reached, with the address they tried
            const message = `Cannot reach the Streamsphere API at ${s.api}`;
            await page.goto(`${base}/#/videos`);
            await expect(page.getByText(message)).toBeVisible();
            await expect(page.locator('p-progressspinner, p-progressSpinner')).toHaveCount(0);
            await expect(page.getByText('There are no videos added to the library!')).toHaveCount(0);
            await shot(page, 'library-api-unreachable');
            await page.goto(`${base}/#/playlists`);
            await expect(page.getByText(message)).toBeVisible();
            await expect(page.getByText('There are no playlists added to the library!')).toHaveCount(0);
            await shot(page, 'playlists-api-unreachable');
            return;
        }

        if (s.contentBroken) {
            // the api address is the page's own, so it works whatever PROTOCOL says; the content is on another place
            // and uses the wrong protocol, so the cards are drawn without their thumbnails
            await page.goto(`${base}/#/videos`);
            await expect.poll(() => [...seen.api]).toEqual([s.api]);
            expect(seen.failedApi).toEqual([]);
            await expect.poll(async () => (await thumbnails(page)).length).toBeGreaterThan(0);
            const shown = await thumbnails(page);
            expect(shown.map(t => t.origin)).toEqual(shown.map(() => s.content));
            expect(shown.some(t => t.loaded)).toBe(false);
            await shot(page, `library-${s.id.split(' ')[0]}`);
            return;
        }

        // the api, and the live status socket
        await page.goto(`${base}/#/videos`);
        await expect.poll(() => [...seen.api]).toEqual([s.api]);
        const socket = s.api!.replace(/^http/, 'ws') + '/ws/status';
        await expect.poll(() => seen.sockets, { message: 'the status socket' }).toContain(socket);

        // library and player
        await expectThumbnailsFrom(page, s.content!);
        await shot(page, `library-${s.id.split(' ')[0]}`);
        await page.locator('p-card', { hasText: 'Me at the zoo' }).locator('img').first().click();
        await expectVideoPlaying(page, s.content!);

        // playlists and the playlist player
        await page.goto(`${base}/#/playlists`);
        await expectThumbnailsFrom(page, s.content!);
        await page.locator('app-playlist-card img, p-card img').first().click();
        await expect(page).toHaveURL(/playlist-details/);
        await expectThumbnailsFrom(page, s.content!);
        await expectVideoPlaying(page, s.content!);

        // a real download: the panel thumbnail and the status events
        if (s.download) {
            await page.goto(`${base}/#/home`);
            expect(await download(page, 'youtube.com/watch?v=jNQXAC9IVRw')).toBe(s.content);
        }

        expect(seen.errors).toEqual([]);
        // the image must ship the optimized bundle, not Angular's development mode
        expect(seen.logs.filter(l => /development mode/i.test(l))).toEqual([]);
    });
}
