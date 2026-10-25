import { Component } from '@angular/core';
import { ComponentFixture, TestBed, fakeAsync, flushMicrotasks, tick } from '@angular/core/testing';
import { Router, provideRouter } from '@angular/router';
import { ConfirmationService, MessageService } from 'primeng/api';
import { throwError } from 'rxjs';

import { VideoDetailsComponent } from './video-details.component';
import { SharedDataService } from '../../services/shared-data.service';
import { VideosService } from '../../services/videos.service';
import { VideoData } from '../../classes/video-data';
import { UrlEncode } from '../../utilities/url-encode';

@Component({ template: '' })
class OtherPageComponent { }

function video(id: number, title: string): VideoData {
    return Object.assign(new VideoData(), {
        video_id: id, title, description: 'about ' + title, media_url: `/media/${id}.mp4`, thumbnail: `/media/${id}.webp`
    });
}

describe('VideoDetailsComponent', () => {
    let fixture: ComponentFixture<VideoDetailsComponent>;
    let component: VideoDetailsComponent;
    let router: Router;
    let saved: VideoData;
    let player: { source: any, destroy: jasmine.Spy };
    let playerBuiltInPage: boolean;

    beforeEach(async () => {
        saved = video(1, 'First');
        TestBed.configureTestingModule({
            imports: [VideoDetailsComponent],
            providers: [provideRouter([{ path: 'other', component: OtherPageComponent }])]
        });
        // the component provides its own services, so they are replaced at component level (the router stays the real one)
        TestBed.overrideComponent(VideoDetailsComponent, {
            set: {
                providers: [
                    { provide: SharedDataService, useValue: { getActivePlayerMetadata: () => saved } },
                    { provide: VideosService, useValue: {} },
                    ConfirmationService, MessageService, UrlEncode
                ]
            }
        });
        router = TestBed.inject(Router);
        fixture = TestBed.createComponent(VideoDetailsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();

        // swap the real player for a recorder, so what the page asks of it can be checked
        playerBuiltInPage = fixture.nativeElement.querySelector('.plyr video') !== null;
        component.player.destroy();
        player = { source: undefined, destroy: jasmine.createSpy('destroy') };
        component.player = player;
    });

    afterEach(() => fixture.destroy());

    it('shows the saved video when it opens', () => {
        expect(component.selectedVideo.title).toBe('First');
        expect(playerBuiltInPage).toBeTrue();
    });

    it('plays the newly chosen video when navigation ends with a different saved video', async () => {
        saved = video(2, 'Second');
        await router.navigateByUrl('/other');

        expect(component.selectedVideo.title).toBe('Second');
        expect(player.source.sources[0].src).toBe('/media/2.mp4');
        expect(player.source.poster).toBe('/media/2.webp');
    });

    it('leaves the player alone when the saved video is the one already shown', async () => {
        await router.navigateByUrl('/other');

        expect(player.source).toBeUndefined();
        expect(component.selectedVideo.title).toBe('First');
    });

    it('destroys the player and stops listening to the router when the page is left', async () => {
        fixture.destroy();
        expect(player.destroy).toHaveBeenCalledTimes(1);

        saved = video(3, 'Third');
        await router.navigateByUrl('/other');
        expect(player.source).toBeUndefined();
    });

    it('has no video to play when nothing was saved', async () => {
        saved = new VideoData();
        await router.navigateByUrl('/other');

        expect(component.hasVideo).toBeFalse();
        expect(player.source).toBeUndefined();
    });

    it('does not ask for a file when there is no video', async () => {
        saved = new VideoData();
        await router.navigateByUrl('/other');
        const svc = fixture.debugElement.injector.get(VideosService) as any;
        svc.download = jasmine.createSpy('download');

        await component.download();

        expect(svc.download).not.toHaveBeenCalled();
    });

    it('tells the user when the download fails', async () => {
        const svc = fixture.debugElement.injector.get(VideosService) as any;
        svc.download = () => Promise.resolve(throwError(() => new Error('404')));
        const toast = spyOn(fixture.debugElement.injector.get(MessageService), 'add');

        await component.download();

        expect(toast).toHaveBeenCalledWith(jasmine.objectContaining({ severity: 'error', detail: 'Download failed' }));
    });

    describe('deleting the video', () => {
        function confirmDeleteRightAway() {
            spyOn(fixture.debugElement.injector.get(ConfirmationService), 'confirm').and.callFake((options: any) => options.accept());
        }

        it('goes to the videos page 3 seconds after the delete', fakeAsync(() => {
            confirmDeleteRightAway();
            (fixture.debugElement.injector.get(VideosService) as any).deleteVideoById = () => Promise.resolve({ data: true });
            const navigate = spyOn(router, 'navigate').and.resolveTo(true);

            component.confirm({ target: null } as any, 1);
            flushMicrotasks();
            expect(navigate).not.toHaveBeenCalled();
            tick(3000);

            expect(navigate).toHaveBeenCalledOnceWith(['/videos']);
        }));

        it('stays on the page when the delete fails', fakeAsync(() => {
            confirmDeleteRightAway();
            (fixture.debugElement.injector.get(VideosService) as any).deleteVideoById = () => Promise.resolve({ data: false });
            const navigate = spyOn(router, 'navigate').and.resolveTo(true);

            component.confirm({ target: null } as any, 1);
            flushMicrotasks();
            tick(3000);

            expect(navigate).not.toHaveBeenCalled();
        }));
    });
});
