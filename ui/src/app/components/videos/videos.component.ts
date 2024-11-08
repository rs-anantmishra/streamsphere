import { Component, OnDestroy, OnInit, ViewChild, effect } from '@angular/core';
import { Paginator, PaginatorModule } from 'primeng/paginator';
import { ButtonModule } from 'primeng/button';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { SimplecardComponent } from "../simplecard/simplecard.component";
import { Router, RouterModule } from '@angular/router'
import { PlaylistCardComponent } from '../playlist-card/playlist-card.component'

//services
import { SharedDataService } from '../../services/shared-data.service'
import { VideosService } from '../../services/videos.service'
import { CommonModule } from '@angular/common';
import { Subscription } from 'rxjs';

@Component({
    selector: 'app-videos',
    standalone: true,
    imports: [SimplecardComponent, CommonModule, PaginatorModule, ButtonModule, ScrollPanelModule, RouterModule, ProgressSpinnerModule, PlaylistCardComponent],
    providers: [Router, SharedDataService],
    templateUrl: './videos.component.html',
    styleUrl: './videos.component.scss'
})

export class VideosComponent implements OnInit, OnDestroy {

    subscription!: Subscription;
    isHomepage = false

    visibility = 'visible'
    first: number = 0;
    rows: number = 10;
    totalRecords: number = -1
    loadPage: number = -1

    //videos presenter
    lstVideos: any
    @ViewChild('paginator', { static: true }) paginator!: Paginator

    constructor(private svcVideos: VideosService, private svcSharedData: SharedDataService) {
        let pageCount = -1;
        this.subscription = this.svcSharedData.getPageSizeCount().subscribe(x => pageCount = x)
        if (pageCount < 0) {
            if (this.svcSharedData.getVideosPageSizeCount() < 0) {
                this.svcSharedData.setPageSizeCount(this.rows)
            } else {
                this.svcSharedData.setPageSizeCount(this.svcSharedData.getVideosPageSizeCount());
            }
        }
        this.subscription = this.svcSharedData.getPageSizeCount().subscribe(rows => this.rows = rows);
    }

    async ngOnInit() {
        await this.getAllVideos();
        this.loadPage = 0
    }

    async getAllVideos() {
        let result = await this.svcVideos.getAllVideos();
        if (result !== null && result.length > 0) {
            this.svcSharedData.setlstVideos(result)
            this.lstVideos = this.getPagedResult(this.first, this.rows);
        } else if (result === null) {
            this.totalRecords = 0
        }
    }

    getPagedResult(first: number, rows: number): any {
        let result = this.svcSharedData.getlstVideos()
        this.totalRecords = result.length
        return result.slice(first, (first + rows))
    }

    onPageChange(event: any) {
        //remember the page-size change
        this.svcSharedData.setPageSizeCount(event.rows)
        //set array to match page
        this.lstVideos = this.getPagedResult(event.first, event.rows)
        this.first = event.first
        this.rows = event.rows;
    }

    ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }
}
