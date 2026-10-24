import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { Subject } from 'rxjs';
import { MessageService } from 'primeng/api';

import { DownloadsComponent } from './downloads.component';
import { DownloadService } from '../../services/download.service';
import { FilesService } from '../../services/files.service';
import { SharedDataService } from '../../services/shared-data.service';
import { Messages } from '../../constants/messages';
import { UrlEncode } from '../../utilities/url-encode';
import { DownloadStatusEvent } from '../../interfaces/download-status-event';

function statusEvent(overrides: Partial<DownloadStatusEvent>): DownloadStatusEvent {
    return { download: '', video_url: '', state: 'downloading', video_id: 1, title: '', channel: '', thumbnail: '', ...overrides };
}

describe('DownloadsComponent status events', () => {
    let fixture: ComponentFixture<DownloadsComponent>;
    let component: DownloadsComponent;
    let events: Subject<DownloadStatusEvent>;
    let sharedData: SharedDataService;
    let messageService: MessageService;

    beforeEach(() => {
        events = new Subject<DownloadStatusEvent>();
        TestBed.configureTestingModule({ imports: [DownloadsComponent] });
        // the component provides its own services, so they are replaced at component level
        TestBed.overrideComponent(DownloadsComponent, {
            set: {
                providers: [
                    { provide: DownloadService, useValue: { statusEvents: () => events, getQueuedItems: () => Promise.resolve([]) } },
                    { provide: FilesService, useValue: {} },
                    MessageService, Messages, UrlEncode
                ]
            }
        });
        fixture = TestBed.createComponent(DownloadsComponent);
        component = fixture.componentInstance;
        sharedData = TestBed.inject(SharedDataService);
        messageService = fixture.debugElement.injector.get(MessageService);
        spyOn(sharedData, 'setRefreshAutoCompleteValue');
        spyOn(messageService, 'add');
    });

    afterEach(() => fixture.destroy());

    it('renders title, channel and progress from a downloading event', () => {
        component.handleStatusEvent(statusEvent({
            download: '[download]  45.3% of  229.20MiB at 19.19MiB/s ETA 00:10',
            title: 'A Title', channel: 'A Channel', thumbnail: 'thumb.jpg'
        }));

        expect(component.activeDLTitle).toBe('A Title');
        expect(component.activeDLChannel).toBe('A Channel');
        expect(component.activeDLImage).toBe('thumb.jpg');
        expect(component.serverLogs).toBe('[download]  45.3% of  229.20MiB at 19.19MiB/s ETA 00:10');
        expect(component.downloadPercent).toBe(45.3);
    });

    it('refreshes search once per completed video', () => {
        component.handleStatusEvent(statusEvent({ state: 'completed', video_id: 7, download: 'Download completed successfully.' }));
        component.handleStatusEvent(statusEvent({ state: 'completed', video_id: 7, download: 'Download completed successfully.' }));
        component.handleStatusEvent(statusEvent({ state: 'completed', video_id: 8, download: 'Download completed successfully.' }));

        expect(sharedData.setRefreshAutoCompleteValue).toHaveBeenCalledTimes(2);
        expect(component.downloadPercent).toBe(100);
    });

    it('refreshes search once when saving goes idle without a completed event', () => {
        component.handleStatusEvent(statusEvent({ state: 'saving', video_id: 5, download: '[download] Saving file to library...' }));
        component.handleStatusEvent(statusEvent({ state: 'idle', video_id: 0 }));
        component.handleStatusEvent(statusEvent({ state: 'idle', video_id: 0 }));

        expect(sharedData.setRefreshAutoCompleteValue).toHaveBeenCalledTimes(1);
        expect(component.downloadPercent).toBe(100);
    });

    it('clears the progress line when a download goes idle mid-way', () => {
        component.handleStatusEvent(statusEvent({ video_id: 6, download: '[download]  45.3% of  229.20MiB', title: 'A Title' }));
        component.handleStatusEvent(statusEvent({ state: 'idle', video_id: 0 }));

        expect(component.serverLogs).toBe('');
        expect(component.downloadPercent).toBeNull();
        expect(component.activeDLTitle).toBe('A Title');
        expect(sharedData.setRefreshAutoCompleteValue).not.toHaveBeenCalled();
    });

    it('refreshes twice when the same video completes, downloads again and completes', () => {
        component.handleStatusEvent(statusEvent({ state: 'completed', video_id: 7, download: 'Download completed successfully.' }));
        component.handleStatusEvent(statusEvent({ video_id: 7, download: '[download]  10.0% of  1.00MiB' }));
        component.handleStatusEvent(statusEvent({ state: 'completed', video_id: 7, download: 'Download completed successfully.' }));

        expect(sharedData.setRefreshAutoCompleteValue).toHaveBeenCalledTimes(2);
    });

    it('shows one toast per failed video', () => {
        component.handleStatusEvent(statusEvent({ state: 'failed', video_id: 9, download: 'ERROR: boom' }));
        component.handleStatusEvent(statusEvent({ state: 'failed', video_id: 9, download: 'ERROR: boom' }));

        expect(messageService.add).toHaveBeenCalledTimes(1);
        expect(sharedData.setRefreshAutoCompleteValue).not.toHaveBeenCalled();
        expect(component.downloadPercent).toBeNull();
    });

    it('keeps the last item on idle', () => {
        component.handleStatusEvent(statusEvent({ state: 'completed', title: 'A Title', channel: 'A Channel', download: 'Download completed successfully.' }));
        component.handleStatusEvent(statusEvent({ state: 'idle', video_id: 0 }));

        expect(component.activeDLTitle).toBe('A Title');
        expect(component.activeDLChannel).toBe('A Channel');
        expect(component.serverLogs).toBe('Download completed successfully.');
    });

    it('keeps the default logs on idle when nothing is shown', () => {
        component.handleStatusEvent(statusEvent({ state: 'idle', video_id: 0 }));

        expect(component.serverLogs).toBe('no logs available');
        expect(component.activeDLTitle).toBe('');
    });

    it('keeps the previous title when an event field is empty', () => {
        component.handleStatusEvent(statusEvent({ title: 'A Title', channel: 'A Channel' }));
        component.handleStatusEvent(statusEvent({ state: 'saving', download: 'saving' }));

        expect(component.activeDLTitle).toBe('A Title');
        expect(component.activeDLChannel).toBe('A Channel');
    });

    it('handles events from the service after init and stops after destroy', () => {
        fixture.detectChanges();
        events.next(statusEvent({ title: 'Live' }));
        expect(component.activeDLTitle).toBe('Live');

        fixture.destroy();
        expect(events.observed).toBeFalse();
    });

    it('rejects an empty URL', async () => {
        const svc = fixture.debugElement.injector.get(DownloadService) as any;
        svc.getMetadata = jasmine.createSpy('getMetadata');
        component.options.Identifier = '';

        await component.GetMedia();

        expect(svc.getMetadata).not.toHaveBeenCalled();
        expect(messageService.add).toHaveBeenCalledWith(jasmine.objectContaining({ severity: 'error', detail: 'Enter a video or playlist URL' }));
        expect(component.loading).toBeFalse();
        expect(component.urlInputDisabled).toBeFalse();
    });

    it('polls the queue while the drawer is open', fakeAsync(() => {
        const svc = fixture.debugElement.injector.get(DownloadService) as any;
        const spy = spyOn(svc, 'getQueuedItems').and.returnValue(Promise.resolve([]));

        component.getQueuedItems(true);
        tick(6000);
        expect(spy.calls.count()).toBeGreaterThanOrEqual(3);

        component.sidebarVisible = false;
        tick(3000);
        const calls = spy.calls.count();
        tick(6000);
        expect(spy.calls.count()).toBe(calls);
        fixture.destroy();
    }));
});
