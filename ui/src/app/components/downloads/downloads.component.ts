
import { Component, OnInit, OnDestroy, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { InputGroupModule } from 'primeng/inputgroup';
import { InputGroupAddonModule } from 'primeng/inputgroupaddon';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { CheckboxModule } from 'primeng/checkbox';
import { PanelModule } from 'primeng/panel';
import { CardModule } from 'primeng/card';
import { DrawerModule } from 'primeng/drawer';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { ProgressBarModule } from 'primeng/progressbar';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { FieldsetModule } from 'primeng/fieldset';
import { DividerModule } from 'primeng/divider';
import { PopoverModule } from 'primeng/popover';
import { Subscription } from 'rxjs';
import { UrlEncode } from '../../utilities/url-encode';


//Services & Classes
import { QueueDownloads, DownloadMedia, VideoData, VideoDataRequest } from '../../classes/video-data';
import { DownloadService } from '../../services/download.service';
import { SimplecardComponent } from "../simplecard/simplecard.component";
import { SharedDataService, Operation } from '../../services/shared-data.service';
import { Messages, Severity } from '../../constants/messages'
import { RemovePrefixPipe } from '../../utilities/pipes/remove-prefix.pipe'
import { FilesService } from '../../services/files.service';
import { FilesizeConversionPipe } from '../../utilities/pipes/filesize-conversion.pipe'
import { DownloadStatusEvent } from '../../interfaces/download-status-event';
import { parseDownloadPercent } from '../../utilities/parse-download-percent';


interface ExtractionOptions {
    Identifier: string;
    GetAudioOnly: boolean;
    GetSubs: boolean;
}


@Component({
    selector: 'app-downloads',
    imports: [ToastModule, ProgressBarModule, FieldsetModule, ProgressSpinnerModule, DrawerModule, CardModule, FormsModule,
        InputGroupModule, InputGroupAddonModule, InputTextModule, ButtonModule, CommonModule, CheckboxModule, PanelModule,
        SimplecardComponent, RemovePrefixPipe, ScrollPanelModule, DividerModule, PopoverModule, FilesizeConversionPipe],
    providers: [DownloadService, MessageService, Messages, FilesService, UrlEncode],
    templateUrl: './downloads.component.html',
    styleUrl: './downloads.component.scss'
})
export class DownloadsComponent implements OnInit, OnDestroy {

    serverLogs: string;
    queuedItems: VideoData[] = []
    nilMetadata = new VideoData()

    loading: boolean = false;
    urlInputDisabled: boolean = false;

    statusSubscription?: Subscription
    downloadPercent: number | null = null
    private completedIds = new Set<number>()
    private failedIds = new Set<number>()

    activeDLImage = ''
    activeDLChannel = ''
    activeDLTitle = ''

    //css-classes
    homeBoxActive = 'home-box'
    contentBoxActive = 'content-box'

    urlPlaceholder = 'Video or Playlist URL'
    sidebarVisible: boolean = false;
    options: ExtractionOptions = { Identifier: '', GetAudioOnly: false, GetSubs: false }
    isPlaylist: boolean = false

    constructor(private messageService: MessageService,
        private svcDownload: DownloadService,
        private svcFiles: FilesService,
        private sharedData: SharedDataService,
        readonly urlEncode: UrlEncode,
        private msg: Messages) {

        this.serverLogs = msg.serverLogs
    }

    async ngOnInit() {
        this.startPlaceholderRotation()

        this.sharedData.isPlaylist = this.sharedData.getIsPlaylist()

        //the server keeps the socket open and sends the current state every second
        this.statusSubscription = this.svcDownload.statusEvents().subscribe(event => this.handleStatusEvent(event))

        //get queued-items on reload
        await this.getQueuedItems(false)
    }

    handleStatusEvent(event: DownloadStatusEvent) {
        //idle - leave the last item displayed
        if (event.state === 'idle') { return }

        this.activeDLTitle = event.title || this.activeDLTitle
        this.activeDLChannel = event.channel || this.activeDLChannel
        this.activeDLImage = event.thumbnail || this.activeDLImage
        this.serverLogs = event.download

        if (event.state === 'downloading') {
            this.downloadPercent = parseDownloadPercent(event.download)
            //a retried or repeated download of the same video can complete or fail again
            this.completedIds.delete(event.video_id)
            this.failedIds.delete(event.video_id)
        } else if (event.state === 'completed') {
            this.downloadPercent = 100
        } else {
            this.downloadPercent = null
        }

        //the server repeats completed/failed until the next item starts - act once per video
        if (event.state === 'completed' && !this.completedIds.has(event.video_id)) {
            this.completedIds.add(event.video_id)
            //the server sends completion only after the file is saved, so search can refresh right away
            this.sharedData.setRefreshAutoCompleteValue(true);
        } else if (event.state === 'failed' && !this.failedIds.has(event.video_id)) {
            this.failedIds.add(event.video_id)
            this.showMessage('Download failed: ' + (event.title || event.download), 'error', 'error')
        }
    }

    flipCheckbox(event: any, option: string): void {
        if (option === 'GetAudioOnly') {
            this.options.GetAudioOnly = !event.checked;
        } else if (option === 'GetSubs') {
            this.options.GetSubs = !event.checked;
        }
    }

    async GetMedia() {
        if (this.options.Identifier.trim() === '') {
            this.showMessage('Enter a video or playlist URL', 'error', 'error')
            return
        }

        this.loading = true;
        this.activeDLTitle = this.msg.getTitle
        this.activeDLChannel = this.msg.getChannel
        this.serverLogs = this.msg.getInfo

        this.urlInputDisabled = true;
        let metadataRequest = this.GetMetadataRequest(this.options)

        try {
            let metadataResponse = await this.svcDownload.getMetadata(metadataRequest)
            if (metadataResponse.status === 'failure') {
                this.showMessage(metadataResponse.message, 'error', 'error')
                //Completion Process
                this.GetMediaCompleteResult(false)
            } else {
                //show the queued item right away, status events take over once it starts downloading
                let metadata = metadataResponse.data;
                if (metadata && metadata.length > 0) {
                    this.activeDLTitle = metadata[0].title
                    this.activeDLChannel = metadata[0].channel
                    this.activeDLImage = metadata[0].thumbnail
                }

                //refresh autocomplete cache
                this.sharedData.setRefreshAutoCompleteValue(true);

                //Completion Process
                this.GetMediaCompleteResult()
            }
        } catch (err) {
            //request failed (api down, network error, invalid response) - reset the ui instead of leaving it stuck
            console.error('error: failed to queue download.', err)
            this.showMessage('Failed to queue download: ' + err, 'error', 'error')
            this.GetMediaCompleteResult(false)
        }
    }

    GetMediaCompleteResult(showMessage: boolean = true) {

        //failure case
        if (this.activeDLTitle == this.msg.getTitle) { this.activeDLTitle = "" }
        if (this.activeDLChannel == this.msg.getChannel) { this.activeDLChannel = "" }
        if (this.serverLogs == this.msg.getInfo) { this.serverLogs = "" }

        //Complete Result
        this.resetDownloadOptions();
        this.loading = false;
        this.urlInputDisabled = false;
        if (showMessage) {
            this.showMessage('Video/Playlist queued', 'info', 'Info')
        }
    }

    GetMetadataRequest(options: ExtractionOptions) {

        let request = new VideoDataRequest()

        request.Indicator = options.Identifier
        request.SubtitlesReq = options.GetSubs
        request.IsAudioOnly = options.GetAudioOnly

        return request
    }

    resetDownloadOptions() {
        this.options.GetAudioOnly = false
        this.options.GetSubs = false
        this.options.Identifier = ''
    }

    async getQueuedItems(openSidebar: boolean) {
        if (openSidebar) {
            this.sidebarVisible = true
            this.startQueuePolling()
        }
        await this.svcDownload.getQueuedItems("queued").then(item => { this.queuedItems = item; })
    }

    //keep the open drawer current
    queueTimer?: ReturnType<typeof setInterval>
    startQueuePolling() {
        if (this.queueTimer) { return }
        this.queueTimer = setInterval(() => {
            if (!this.sidebarVisible) { this.stopQueuePolling(); return }
            this.getQueuedItems(false)
        }, 3000)
    }

    stopQueuePolling() {
        clearInterval(this.queueTimer)
        this.queueTimer = undefined
    }

    fsStorageStatus: number = 0
    dbStorageStatus: number = 0
    async getStorageStatus() {
        try {
            let filesResult = await this.svcFiles.getStorageStatus();
            this.fsStorageStatus = filesResult.data.storage_used_fs;
            this.dbStorageStatus = filesResult.data.storage_used_db;
        } catch (e) {
            console.log('error fetching storage status:', e);
        }
    }

    ngOnDestroy() {
        // prevent memory leak when component destroyed
        this.statusSubscription?.unsubscribe();
        clearInterval(this.placeholderTimer);
        this.stopQueuePolling();
    }

    //placeholders = ['Video or Playlist URL', 'youtube.com/watch?v=cXdwGt_xPCI', 'cXdwGt_xPCI', 'youtu.be/cXdwGt_xPCI?si=r7GvF9rndAtAkb0n', 'youtube.com/watch?v=cABfEaP2IHo&list=PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU', 'youtube.com/playlist?list=PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU', 'PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU'];
    placeholders = ['Video or Playlist URL', 'youtube.com/watch?v=videoId', 'videoId', 'youtu.be/videoId', 'youtube.com/watch?v=videoId&list=playlistId', 'youtube.com/playlist?list=playlistId', 'playlistId'];
    index = 0;
    placeholderTimer?: ReturnType<typeof setInterval>;
    startPlaceholderRotation() {
        this.placeholderTimer = setInterval(() => {
            if (this.index === this.placeholders.length - 1) { this.index = 0; } //reset
            this.index = this.index + 1;
            this.urlPlaceholder = this.placeholders[this.index]
        }, 3000);
    }

    //Toast Messages
    showMessage(message: string, severity: string, summary: string) {
        this.messageService.add({ severity: severity, summary: summary, detail: message });
    }
}


