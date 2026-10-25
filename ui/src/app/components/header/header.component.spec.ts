import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NavigationEnd, Router } from '@angular/router';
import { MessageService } from 'primeng/api';
import { Subject } from 'rxjs';

import { HeaderComponent } from './header.component';
import { VideosService } from '../../services/videos.service';

describe('HeaderComponent homepage flag', () => {
    let fixture: ComponentFixture<HeaderComponent>;
    let component: HeaderComponent;
    let routerEvents: Subject<NavigationEnd>;
    let routerStub: { url: string, events: Subject<NavigationEnd>, navigate: jasmine.Spy };

    function navigationEnd(url: string): NavigationEnd {
        return new NavigationEnd(1, url, url);
    }

    function createHeader(url: string) {
        routerEvents = new Subject<NavigationEnd>();
        routerStub = { url, events: routerEvents, navigate: jasmine.createSpy('navigate') };
        TestBed.configureTestingModule({ imports: [HeaderComponent] });
        // the component provides its own services, so they are replaced at component level
        TestBed.overrideComponent(HeaderComponent, {
            set: {
                providers: [
                    { provide: Router, useValue: routerStub },
                    { provide: VideosService, useValue: { getContentSearchInfo: () => Promise.resolve({ data: null }) } },
                    MessageService
                ]
            }
        });
        fixture = TestBed.createComponent(HeaderComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    }

    beforeEach(() => {
        // the theme code styles the page body, which index.html names app-dlbg
        document.body.id = 'app-dlbg';
    });

    afterEach(() => {
        fixture?.destroy();
        document.body.removeAttribute('id');
        document.body.className = '';
        document.documentElement.classList.remove('app-dark');
        localStorage.clear();
    });

    it('is the homepage when the app starts on /home', () => {
        createHeader('/home');
        expect(component.isHomepage).toBeTrue();
    });

    it('is not the homepage when the app is reloaded on another page', () => {
        createHeader('/videos');
        expect(component.isHomepage).toBeFalse();
    });

    it('follows the router whichever way the page was reached', () => {
        createHeader('/home');

        routerEvents.next(navigationEnd('/videos/play'));
        expect(component.isHomepage).toBeFalse();

        routerEvents.next(navigationEnd('/home'));
        expect(component.isHomepage).toBeTrue();
    });

    it('stops following the router once destroyed', () => {
        createHeader('/home');
        fixture.destroy();

        expect(routerEvents.observed).toBeFalse();
    });
});
