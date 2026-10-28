import { expect, test } from '@playwright/test';
import { DatabaseSync } from 'node:sqlite';
import { copyDb, PORTS, startApp, stopApp } from '../support/stack';
import { download, stubPlayerIcons, useDarkTheme, watch } from '../support/page';

// runs first: puts a real downloaded video (needs the network) and one playlist into the app's library,
// which every scenario in scenarios.spec.ts then reads
test('the library is seeded with a real download and a playlist', async ({ page }) => {
    test.setTimeout(180_000);
    const started = await startApp({ APPLICATION_PORT: String(PORTS.app), CONTENT_PORT: String(PORTS.content) });
    expect(started.state, started.logs).toBe('running');

    await stubPlayerIcons(page);
    const seen = watch(page);
    await page.goto(`http://localhost:${PORTS.app}/#/home`);
    await useDarkTheme(page);
    const panelOrigin = await download(page, 'youtube.com/watch?v=jNQXAC9IVRw');
    expect(panelOrigin).toBe(`http://localhost:${PORTS.content}`);
    expect(seen.errors).toEqual([]);

    // a playlist needs two rows pointing at the downloaded video and its files
    stopApp();
    const file = copyDb('out');
    const db = new DatabaseSync(file);
    const files = db.prepare('select Id, FileType from tblFiles where VideoId = 1').all() as { Id: number, FileType: string }[];
    const thumbnail = files.find(f => f.FileType === 'Thumbnail')!.Id;
    const video = files.find(f => f.FileType === 'Video')!.Id;
    const now = Math.floor(Date.now() / 1000);
    db.prepare("insert into tblPlaylists(Title, ItemCount, PlaylistChannel, PlaylistUploader, ThumbnailFileId, YoutubePlaylistId, CreatedDate) values('Test Playlist', 1, 'jawed', 'jawed', ?, 1, ?)").run(thumbnail, now);
    const playlist = (db.prepare('select max(Id) as id from tblPlaylists').get() as { id: number }).id;
    db.prepare('insert into tblPlaylistVideoFiles(VideoId, PlaylistId, PlaylistVideoIndex, FileId, ThumbnailFileId, CreatedDate) values(1, ?, 1, ?, ?, ?)').run(playlist, video, thumbnail, now);
    db.close();
    copyDb('in');
});
