import { Injectable } from '@angular/core';
import { Socket } from 'ngx-socket-io';
import { apiBaseUrl, toWsUrl } from '../utilities/url-transforms';

@Injectable({
    providedIn: 'root',
})
export class WebSocketService {
    private webSocket!: Socket;
    constructor() {
        let wsUrl = toWsUrl(apiBaseUrl()) + '/ws'
        this.webSocket = new Socket({
            url:  wsUrl + '/status',
            options: {},
        });
    }

    // this method is used to start connection/handhshake of socket with server
    connectSocket(message: any) {
        this.webSocket.emit('connect', message);
    }

    // this method is used to get response from server
    receiveStatus() {
        return this.webSocket.fromEvent('/ws/status')
    }

    //this method is used to send request to server
    sendStatus(message: string) {
        this.webSocket.on("connection", () => {
            this.webSocket.emit(message);
        })
    }

    // this method is used to end web socket connection
    disconnectSocket() {
        this.webSocket.disconnect();
    }
}