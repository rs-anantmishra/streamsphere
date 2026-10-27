import { expect, test } from '@playwright/test';
import * as fs from 'node:fs';
import * as path from 'node:path';
import { boundFolders, PORTS, removeBound, startBound } from '../support/stack';
import { download, stubPlayerIcons, useDarkTheme, watch } from '../support/page';

// #47: on host folders (bind mounts) the image's database was hidden, the app ran on an empty one, and downloads
// landed outside the library while the page said they completed
test.afterAll(() => removeBound());

const filesIn = (dir: string): string[] => fs.readdirSync(dir, { withFileTypes: true })
    .flatMap(entry => entry.isDirectory() ? filesIn(path.join(dir, entry.name)) : [path.join(dir, entry.name)]);

test('a fresh install on empty host folders downloads into the library', async ({ page }) => {
    test.setTimeout(240_000);
    const folders = boundFolders();
    const started = await startBound(folders);
    expect(started.state, started.logs).toBe('running');

    await stubPlayerIcons(page);
    const seen = watch(page);
    await page.goto(`http://localhost:${PORTS.bare}/#/home`);
    await useDarkTheme(page);
    await download(page, 'youtube.com/watch?v=jNQXAC9IVRw');

    await page.goto(`http://localhost:${PORTS.bare}/#/videos`);
    await expect(page.getByText('Me at the zoo').first()).toBeVisible();
    expect(filesIn(folders.content).filter(f => /\[jNQXAC9IVRw\]\.(mp4|webm|mkv)$/.test(f))).toHaveLength(1);
    expect(seen.errors).toEqual([]);
});

test('a database folder it cannot write to stops the start and names the folder', async () => {
    test.setTimeout(240_000);
    const folders = boundFolders();
    const started = await startBound(folders, { readOnlyDb: true });
    expect(started.state, started.logs).toBe('exited');
    expect(started.exitCode).toBe(1);
    expect(started.logs).toContain('/app/database/db');
});
