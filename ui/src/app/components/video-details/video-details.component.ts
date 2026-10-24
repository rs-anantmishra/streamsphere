import { Component, ElementRef, OnInit, OnDestroy, ViewEncapsulation, ViewChild } from '@angular/core';
import { NavigationEnd, Router, RouterModule } from '@angular/router';
import { Subscription, filter } from 'rxjs';
import { ButtonModule } from 'primeng/button';
import { SharedDataService } from '../../services/shared-data.service';
import { VideoData } from '../../classes/video-data';
import { PanelModule } from 'primeng/panel';
import { TagModule } from 'primeng/tag';
import { ChipModule } from 'primeng/chip';
import { FieldsetModule } from 'primeng/fieldset';
import { ConfirmPopup, ConfirmPopupModule } from 'primeng/confirmpopup';
import { ConfirmationService, MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { UrlEncode } from '../../utilities/url-encode';
import { linkify } from '../../utilities/linkify';
import { cp1252_to_utf8 } from '../../utilities/cp1252-to-utf8';
import { downloadBlob } from '../../utilities/download-blob';
import { confirmDelete } from '../../utilities/confirm-delete';
import Plyr from 'plyr';

import { MinifiedViewCount } from '../../utilities/pipes/views-conversion.pipe'
import { MinifiedLikeCount } from '../../utilities/pipes/likes-conversion.pipe';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CommaSepStringFromArray } from "../../utilities/pipes/array-comma-sep.pipe";
import { FormattedResolutionPipe } from "../../utilities/pipes/format-resolution.pipe";
import { VideosService } from '../../services/videos.service';
import { MinifiedDatePipe } from "../../utilities/pipes/formatted-date.pipe";
import { FilesizeConversionPipe } from "../../utilities/pipes/filesize-conversion.pipe";

@Component({
    selector: 'app-video-details',
    imports: [RouterModule, ButtonModule, PanelModule, ScrollPanelModule, TagModule, ChipModule, ConfirmPopupModule, ToastModule,
        MinifiedViewCount, MinifiedLikeCount, CommaSepStringFromArray, FormattedResolutionPipe, MinifiedDatePipe, FieldsetModule, FilesizeConversionPipe],
    providers: [Router, SharedDataService, ConfirmationService, MessageService, UrlEncode],
    templateUrl: './video-details.component.html',
    styleUrl: './video-details.component.scss',
    encapsulation: ViewEncapsulation.None
})
export class VideoDetailsComponent implements OnInit, OnDestroy {

    public player: any;
    private routeSubscription?: Subscription;
    selectedVideo: VideoData = new VideoData()
    @ViewChild(ConfirmPopup) confirmPopup!: ConfirmPopup;
    @ViewChild('playerVideo', { static: true }) playerVideo!: ElementRef<HTMLVideoElement>;
    constructor(private confirmationService: ConfirmationService, private messageService: MessageService, private router: Router, private urlEncode: UrlEncode,
        private svcSharedData: SharedDataService, private svcVideos: VideosService) {
    }

    async ngOnInit(): Promise<void> {
        this.loadSelectedVideo();

        this.player = new Plyr(this.playerVideo.nativeElement, { captions: { active: true }, keyboard: { global: true }, autoplay: true });
        this.playSelectedVideo();

        //picking another video while on this page is a navigation to the same url - the page is reused, so look again
        this.routeSubscription = this.router.events.pipe(filter(e => e instanceof NavigationEnd))
            .subscribe(() => this.switchVideoIfChanged());
    }

    ngOnDestroy(): void {
        this.routeSubscription?.unsubscribe();
        this.player?.destroy();
    }

    private loadSelectedVideo() {
        this.selectedVideo = this.svcSharedData.getActivePlayerMetadata();

        this.selectedVideo.media_url = this.selectedVideo.media_url
        this.selectedVideo.thumbnail = this.selectedVideo.thumbnail
        this.selectedVideo.webpage_url = this.selectedVideo.webpage_url

        this.selectedVideo.description = cp1252_to_utf8(this.selectedVideo.description)
        this.selectedVideo.description = linkify(this.selectedVideo.description)
    }

    private switchVideoIfChanged() {
        if (this.svcSharedData.getActivePlayerMetadata().video_id === this.selectedVideo.video_id) { return }

        this.loadSelectedVideo();
        this.playSelectedVideo();
    }

    private playSelectedVideo() {
        this.player.source = {
            type: 'video',
            sources: [{ src: this.selectedVideo.media_url, type: 'video/webm' }],
            poster: this.selectedVideo.thumbnail
        };
    }

    async download(): Promise<void> {
        (await this.svcVideos
            .download(this.selectedVideo.media_url))
            .subscribe(blob => downloadBlob(blob, this.selectedVideo.title + '.' + this.selectedVideo.extension))
    }

    accept() {
        this.confirmPopup.onAccept();
    }

    reject() {
        this.confirmPopup.onReject();
    }

    confirm(event: Event, contentId: number) {
        confirmDelete(event, contentId, id => this.svcVideos.deleteVideoById(id), this.confirmationService, this.messageService, this.router);
    }
}
