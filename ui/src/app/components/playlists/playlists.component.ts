import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { Paginator, PaginatorModule } from 'primeng/paginator';
import { Router, RouterModule } from '@angular/router'
import { ButtonModule } from 'primeng/button';
import { SharedDataService } from '../../services/shared-data.service';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { PlaylistCardComponent } from "../playlist-card/playlist-card.component";
import { Subscription } from 'rxjs';
import { PlaylistsService } from '../../services/playlists.service';


interface PageEvent {
    first: number;
    rows: number;
    page: number;
    pageCount: number;
}

@Component({
    selector: 'app-playlists',
    imports: [CommonModule, PaginatorModule, RouterModule, ButtonModule, ProgressSpinnerModule, ScrollPanelModule, PlaylistCardComponent],
    providers: [Router, SharedDataService],
    templateUrl: './playlists.component.html',
    styleUrl: './playlists.component.scss'
})
export class PlaylistsComponent implements OnInit {

    subscription!: Subscription;
    isHomepage = false

    visibility = 'visible'
    first: number = 0;
    rows: number = 10;
    totalRecords: number = -1
    loadPage: number = -1

    @ViewChild('paginator', { static: false }) paginator!: Paginator
    lstPlaylists: any

    constructor(private svcPlaylists: PlaylistsService, private svcSharedData: SharedDataService) {
    }

    async ngOnInit() {
        this.initPageSize();
        await this.getAllPlaylists();
        this.loadPage = 0
    }

    private initPageSize(): void {
        let pageCount = -1;
        this.svcSharedData.getPageSizeCountPlaylist().subscribe(x => pageCount = x).unsubscribe()
        if (pageCount < 0) {
            if (this.svcSharedData.getPlaylistsPageSizeCount() < 0) {
                this.svcSharedData.setPageSizeCountPlaylist(this.rows)
            } else {
                this.svcSharedData.setPageSizeCountPlaylist(this.svcSharedData.getPlaylistsPageSizeCount());
            }
        }
        this.subscription = this.svcSharedData.getPageSizeCountPlaylist().subscribe(rows => this.rows = rows);
    }

    async getAllPlaylists() {
        try {
            let result = await this.svcPlaylists.getAllPlaylists();
            if (result.data && result.data.length > 0) {
                this.svcSharedData.setlstPlaylists(result.data)
                this.lstPlaylists = this.getPagedResult(this.first, this.rows);
            } else {
                this.totalRecords = 0
            }
        } catch (err) {
            console.error('error: failed to load playlists.', err)
            this.totalRecords = 0
        }
    }

    getPagedResult(first: number, rows: number): any {
        let result = this.svcSharedData.getlstPlaylists()
        this.totalRecords = result.length
        return result.slice(first, (first + rows))
    }

    onPageChange(event: any) {
        //remember the page-size change
        this.svcSharedData.setPageSizeCountPlaylist(event.rows)
        //set array to match page
        this.lstPlaylists = this.getPagedResult(event.first, event.rows)
        this.first = event.first
        this.rows = event.rows;
    }

    ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

}