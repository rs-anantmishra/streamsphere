import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { SharedDataService } from './shared-data.service';
import { StorageStatus, StorageStatusResponse } from '../classes/storage';
import { environment } from '../../environments/environment';

const apiUrl: string = environment.baseUrl + "/api"

@Injectable({
    providedIn: 'root'
})
export class FilesService {
    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    //getStorageStatus
    async getStorageStatus(): Promise<StorageStatusResponse> {
        let url = '/storage/status'
        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })
    }
}

/*
    //getStorageStatus
    getStorageStatus(): Observable<StorageStatusResponse> {
        let url = apiUrl + '/storage/status'
        return new Observable(subscriber => {
            this.http.get(url)
                .subscribe(response => {
                    subscriber.next(JSON.parse(JSON.stringify(response)));
                })
        });
    }

// service.ts
getData(): Observable<any> {
    return new Observable(subscriber => {
        this.http.get(url)
          .pipe(catchError(this.handleError)
          .subscribe(res => {
              // Do my service.ts logic.
              // ...
              subscriber.next(res)
              subscriber.complete()
          }, err => subscriber.error(err))
    })
}

// component.ts
ngOnInit() {
    this.service.getData().subscribe(res => {
        // Do my component logic.
        // ...
    }, err => this.errors = err)
}

*/