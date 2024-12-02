import { CommonModule, DOCUMENT } from '@angular/common';
import { HostListener, Component, OnInit, inject, OnDestroy } from '@angular/core';
import { MenuItem, MessageService } from 'primeng/api';
import { SplitButtonModule } from 'primeng/splitbutton';
import { ToastModule } from 'primeng/toast';
import { FilterService, SelectItemGroup } from 'primeng/api';
import { AutoCompleteCompleteEvent, AutoCompleteModule } from 'primeng/autocomplete';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { InputSwitchModule } from 'primeng/inputswitch';
import { SharedDataService } from '../../services/shared-data.service';
import { VideosService } from '../../services/videos.service';
import { ContentSearchResponse, ContentSearch } from '../../classes/search';
import { Subscription } from 'rxjs';
import { VideoData, VideoDataResponse } from '../../classes/video-data';

@Component({
    selector: 'app-header',
    standalone: true,
    imports: [InputSwitchModule, CommonModule, SplitButtonModule, ToastModule, FormsModule, AutoCompleteModule],
    providers: [MessageService, Router, VideosService],
    templateUrl: './header.component.html',
    styleUrl: './header.component.scss'
})
export class HeaderComponent implements OnInit, OnDestroy {

    //isHomepage route?
    isHomepage = false

    //subs updates this
    searchCacheSubscription!: Subscription;

    #document = inject(DOCUMENT);
    themeIcon = ''
    activeIcon = '#dark'
    toggleTheme() {
        const linkElement = this.#document.getElementById('app-theme',) as HTMLLinkElement;
        const bodyElement = this.#document.getElementById('app-dlbg',) as HTMLBodyElement;

        if (linkElement.href.includes('light')) {
            this.setDarkMode();
        } else {
            this.setLightMode();
        }
    }

    setLightMode() {
        const linkElement = this.#document.getElementById('app-theme',) as HTMLLinkElement;
        const bodyElement = this.#document.getElementById('app-dlbg',) as HTMLBodyElement;

        linkElement.href = 'themes/aura-light-blue/theme.css';
        bodyElement.className = "downloads-bg-light"
        this.activeIcon = '#dark'
        this.sharedDataSvc.setIsDarkMode(false)
    }

    setDarkMode() {
        const linkElement = this.#document.getElementById('app-theme',) as HTMLLinkElement;
        const bodyElement = this.#document.getElementById('app-dlbg',) as HTMLBodyElement;

        linkElement.href = 'themes/aura-dark-blue/theme.css';
        bodyElement.className = "downloads-bg-dark"
        this.activeIcon = '#light'
        this.sharedDataSvc.setIsDarkMode(true)
    }

    //search-bar
    visible: string = 'visible'

    navItems!: MenuItem[];
    filteredGroups!: any[];
    selectedTitle: any;
    groupedTitles!: SelectItemGroup[];

    constructor(private router: Router,
        private videosSvc: VideosService,
        private filterService: FilterService,
        public sharedDataSvc: SharedDataService) {

        //check and set if homepage
        this.checkIsHomepage('/home')

        //nav-items
        this.initNavItems();

        //check and set theme
        let isDarkMode = this.sharedDataSvc.getIsDarkMode();
        if (isDarkMode === null) {
            this.sharedDataSvc.setIsDarkMode(true)
        } else {
            if (isDarkMode === true) {
                this.setDarkMode()
            } else if (isDarkMode === false) {
                this.setLightMode()
            }
        }
    }
    ngOnDestroy(): void {
        //unsubscribe
        this.searchCacheSubscription.unsubscribe()
    }

    navigate(route: string) {
        this.checkIsHomepage(route)
        this.router.navigate([route]);
    }

    home: MenuItem | undefined;
    async ngOnInit(): Promise<void> {
        //update cache
        this.searchCacheSubscription = this.sharedDataSvc._refreshAutoComplete$.subscribe(() => { this.buildAutoCompleteCache(); })
        await this.buildAutoCompleteCache();
    }

    async buildAutoCompleteCache() {
        let result = await this.videosSvc.getContentSearchInfo()
        if (result.data !== null) {
            this.groupedTitles = await this.buildAutoCompleteDataset(result)
        }
    }

    async buildAutoCompleteDataset(raw: ContentSearchResponse): Promise<SelectItemGroup[]> {
        let content: ContentSearch[] = raw.data;
        let result: SelectItemGroup[] = [];

        //group titles by channel
        let grouped = content.reduce(
            (result: any, currentValue: any) => {
                (result[currentValue['channel']] = result[currentValue['channel']] || []).push({ "label": currentValue['title'], 'value': currentValue['video_id'] });
                return result;
            }, {});

        //format json in requires manner
        for (let key in grouped) {
            if (grouped.hasOwnProperty(key)) {
                let val = { label: key, value: "", items: grouped[key] };
                result.push(val)
            }
        }
        return result;
    }

    //keyboard shortcuts
    @HostListener("document:keydown", ["$event"]) handleKeyboardEvent(event: KeyboardEvent) {
        if (event.key === 'P') {
            this.navigate('/playlists')
        }
        if (event.key === 'V') {
            this.navigate('/videos')
        }
        if (event.key === 'C' && event.altKey) {
            this.navigate('/channels')
        }
        if (event.key === 'L' && event.altKey) {
            this.navigate('/logs')
        }
        if ((event.key === 'H' || event.key === 'D')) {
            this.navigate('/home')
        }
    }

    filterGroupedContent(event: AutoCompleteCompleteEvent) {
        let query = event.query;
        let filteredGroups = [];

        for (let optgroup of this.groupedTitles) {
            let filteredSubOptions = this.filterService.filter(optgroup.items, ['label'], query, "contains");
            if (filteredSubOptions && filteredSubOptions.length) {
                filteredGroups.push({
                    label: optgroup.label,
                    value: optgroup.value,
                    items: filteredSubOptions
                });
            }
        }
        this.filteredGroups = filteredGroups;
    }

    initNavItems() {
        this.navItems = [
            { label: 'Home', routerLink: ['/home'], command: () => { this.navigate('/home'); } },
            { separator: true },
            { label: 'Videos', routerLink: ['/videos'], command: () => { this.navigate('/videos'); } },
            { label: 'Playlists', routerLink: ['/playlists'], command: () => { this.navigate('/playlists'); } },
            // { label: 'Tags', routerLink: ['/tags'], command: () => { this.navigate('/tags'); } },
            // { label: 'Categories', routerLink: ['/categories'], command: () => { this.navigate('/categories'); } },
            // { separator: true },
            // { label: 'Pattern Matching', routerLink: ['/recursive'] },
            // { label: 'Saved Patterns', routerLink: ['/notes'] },
            // { label: 'Source RegEx', routerLink: ['/source'] },

            // { separator: true },
            // { label: 'Activity Logs', routerLink: ['/activity-logs'], command: () => { this.navigate('/logs'); } },
        ];
    }

    checkIsHomepage(route: string) {
        //if any of the mentioned reourtes are not going to PROC then its going to home.
        if (route === '/videos' || route === '/tags' || route === '/categories' || route === '/videos/play' || route === '/playlists' || route === '/playlist-details') {
            this.isHomepage = false
        } else {
            this.isHomepage = true
        }
    }

    //search result - clicked
    async navigateToVideo(selected: any) {
        //click on an empty search field
        if (selected !== undefined) {
            let result = await this.videosSvc.getContentById(selected.value)
            if (result.data != null) {
                this.selectedVideo(result.data[0])
            }
        }
    }

    //search result - enter pressed
    async navigateToVideoOnEnter(selected: any, $event: KeyboardEvent) {
        if ($event.key == 'Enter') {
            //click on an empty search field
            if (selected !== undefined) {
                let result = await this.videosSvc.getContentById(selected.value)
                if (result.data != null) {
                    this.selectedVideo(result.data[0])
                }
            }
        }
    }

    selectedVideo(playVideo: VideoData) {
        this.selectedTitle = ''
        this.sharedDataSvc.setPlayVideo(playVideo);
        this.router.navigate(['videos', 'play'])
    }
}