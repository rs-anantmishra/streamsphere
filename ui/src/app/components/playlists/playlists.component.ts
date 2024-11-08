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
    standalone: true,
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

    @ViewChild('paginator', { static: true }) paginator!: Paginator
    lstPlaylists: any

    constructor(private svcPlaylists: PlaylistsService, private svcSharedData: SharedDataService) {
        let pageCount = -1;
        this.subscription = this.svcSharedData.getPageSizeCountPlaylist().subscribe(x => pageCount = x)
        if (pageCount < 0) {
            if (this.svcSharedData.getPlaylistsPageSizeCount() < 0) {
                this.svcSharedData.setPageSizeCountPlaylist(this.rows)
            } else {
                this.svcSharedData.setPageSizeCountPlaylist(this.svcSharedData.getPlaylistsPageSizeCount());
            }
        }
        this.subscription = this.svcSharedData.getPageSizeCountPlaylist().subscribe(rows => this.rows = rows);
    }

    async ngOnInit() {
        await this.getAllPlaylists();
        this.loadPage = 0
    }

    async getAllPlaylists() {
        let result = await this.svcPlaylists.getAllPlaylists();
        if (result.data !== null && result.data.length > 0) {
            this.svcSharedData.setlstPlaylists(result.data)
            this.lstPlaylists = this.getPagedResult(this.first, this.rows);
        } else if (result.data === null) {
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