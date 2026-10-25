import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router, provideRouter } from '@angular/router';

import { PlaylistDetailsComponent } from './playlist-details.component';
import { SharedDataService } from '../../services/shared-data.service';
import { PlaylistsService } from '../../services/playlists.service';
import { PlaylistsInfo, SelectedPlaylist } from '../../classes/playlists';
import { VideoData } from '../../classes/video-data';

function video(id: number, index: number): VideoData {
    return Object.assign(new VideoData(), {
        video_id: id, playlist_video_index: index, title: 'video ' + id, description: 'about ' + id,
        media_url: `/media/${id}.mp4`, thumbnail: `/media/${id}.webp`
    });
}

function playlist(videos: VideoData[], activeVideo = -1): SelectedPlaylist {
    return Object.assign(new SelectedPlaylist(), {
        info: Object.assign(new PlaylistsInfo(), { playlist_id: 7, playlist_title: 'A playlist' }),
        active_video: activeVideo,
        video_data: videos
    });
}

describe('PlaylistDetailsComponent', () => {
    let fixture: ComponentFixture<PlaylistDetailsComponent>;
    let component: PlaylistDetailsComponent;
    let sharedData: SharedDataService;
    let playlistsSvc: PlaylistsService;
    let navigate: jasmine.Spy;

    beforeEach(() => {
        localStorage.removeItem('playlist');
        TestBed.configureTestingModule({ imports: [PlaylistDetailsComponent], providers: [provideRouter([])] });
        navigate = spyOn(TestBed.inject(Router), 'navigate').and.resolveTo(true);
        sharedData = TestBed.inject(SharedDataService);
        playlistsSvc = TestBed.inject(PlaylistsService);
        fixture = TestBed.createComponent(PlaylistDetailsComponent);
        component = fixture.componentInstance;
        spyOn(component, 'createPlayer');
    });

    afterEach(() => {
        fixture.destroy();
        localStorage.removeItem('playlist');
    });

    describe('getPlaylistsVideos', () => {
        beforeEach(() => sharedData.setPlaylist(playlist([])));

        it('gives no videos for a null response', async () => {
            spyOn(console, 'error');
            spyOn(playlistsSvc, 'getPlaylistVideos').and.resolveTo(null as any);

            expect(await component.getPlaylistsVideos(7)).toEqual([]);
        });

        it('gives no videos when the response has null data', async () => {
            spyOn(playlistsSvc, 'getPlaylistVideos').and.resolveTo({ data: null } as any);

            expect(await component.getPlaylistsVideos(7)).toEqual([]);
        });

        it('gives no videos for an empty list', async () => {
            spyOn(playlistsSvc, 'getPlaylistVideos').and.resolveTo({ data: [] } as any);

            expect(await component.getPlaylistsVideos(7)).toEqual([]);
        });

        it('gives no videos, and logs, when the request fails', async () => {
            const error = spyOn(console, 'error');
            spyOn(playlistsSvc, 'getPlaylistVideos').and.rejectWith(new Error('down'));

            expect(await component.getPlaylistsVideos(7)).toEqual([]);
            expect(error).toHaveBeenCalled();
        });

        it('returns the videos the service gives', async () => {
            const videos = [video(1, 1)];
            spyOn(playlistsSvc, 'getPlaylistVideos').and.resolveTo({ data: videos } as any);

            expect(await component.getPlaylistsVideos(7)).toBe(videos);
        });

        it('uses the cached videos without asking the service', async () => {
            const cached = [video(2, 1)];
            sharedData.setPlaylist(playlist(cached));
            const get = spyOn(playlistsSvc, 'getPlaylistVideos');

            expect((await component.getPlaylistsVideos(7)).map(v => v.video_id)).toEqual([2]);
            expect(get).not.toHaveBeenCalled();
        });
    });

    describe('ngOnInit', () => {
        it('goes back to the playlists when no playlist was picked', async () => {
            await component.ngOnInit();

            expect(navigate).toHaveBeenCalledWith(['/playlists']);
            expect(component.loaded).toBeFalse();
            expect(component.createPlayer).not.toHaveBeenCalled();
        });

        it('finishes loading an empty playlist with nothing selected', async () => {
            sharedData.setPlaylist(playlist([]));
            spyOn(playlistsSvc, 'getPlaylistVideos').and.resolveTo({ data: null } as any);

            await component.ngOnInit();
            fixture.detectChanges();

            expect(component.loaded).toBeTrue();
            expect(component.playlistVideos).toEqual([]);
            expect(component.hasVideo).toBeFalse();
            expect(navigate).not.toHaveBeenCalled();
            expect(fixture.nativeElement.textContent).toContain('This playlist has no videos.');
        });

        it('finishes loading when the request for the videos fails', async () => {
            spyOn(console, 'error');
            sharedData.setPlaylist(playlist([]));
            spyOn(playlistsSvc, 'getPlaylistVideos').and.rejectWith(new Error('down'));

            await component.ngOnInit();

            expect(component.loaded).toBeTrue();
            expect(component.hasVideo).toBeFalse();
        });

        it('falls back to the first video when the saved active video is no longer in the playlist', async () => {
            sharedData.setPlaylist(playlist([video(1, 1), video(2, 2)], 99));

            await component.ngOnInit();

            expect(component.selectedVideo.video_id).toBe(1);
            expect(component.playlist.active_video).toBe(1);
        });

        it('sorts the videos and selects the first one', async () => {
            sharedData.setPlaylist(playlist([video(2, 2), video(1, 1)]));

            await component.ngOnInit();
            fixture.detectChanges();

            expect(component.loaded).toBeTrue();
            expect(component.playlistVideos.map(v => v.video_id)).toEqual([1, 2]);
            expect(component.selectedVideo.video_id).toBe(1);
            expect(component.hasVideo).toBeTrue();
            expect(fixture.nativeElement.textContent).not.toContain('This playlist has no videos.');
        });
    });

    describe('changeContent', () => {
        let videos: VideoData[];

        beforeEach(async () => {
            videos = [video(1, 1), video(2, 2)];
            videos[1].media_url = '/media/a #2 b.mp4';
            sharedData.setPlaylist(playlist(videos));
            await component.ngOnInit();
            component.player = { source: undefined };
        });

        it('switches the video through the player, without a reload', () => {
            component.changeContent(2);

            expect(component.selectedVideo.video_id).toBe(2);
            expect(component.player.source.sources[0].src).toBe('/media/a %232 b.mp4');
            expect(component.player.source.poster).toBe('/media/2.webp');
            expect(component.loaded).toBeTrue();
        });

        it('saves the new active video so a reload restores it', () => {
            component.changeContent(2);

            expect(sharedData.getPlaylist().active_video).toBe(2);
        });

        it('encodes the url once however often the user switches back and forth', () => {
            component.changeContent(2);
            component.changeContent(1);
            component.changeContent(2);
            component.changeContent(1);
            component.changeContent(2);

            expect(component.player.source.sources[0].src).toBe('/media/a %232 b.mp4');
            expect(component.selectedVideo.media_url).toBe('/media/a %232 b.mp4');
        });

        it('leaves the cached videos, which are saved back to storage, as they came', () => {
            component.changeContent(2);

            expect(component.playlist.video_data.find(v => v.video_id === 2)!.media_url).toBe('/media/a #2 b.mp4');
            expect(sharedData.getPlaylist().video_data.find(v => v.video_id === 2)!.media_url).toBe('/media/a #2 b.mp4');
        });

        it('ignores an unknown video and the one already selected', () => {
            component.changeContent(99);
            component.changeContent(1);

            expect(component.player.source).toBeUndefined();
            expect(component.selectedVideo.video_id).toBe(1);
        });
    });
});
