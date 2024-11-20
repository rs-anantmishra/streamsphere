export class Messages {
    //arrays
    Severities: string[] = ['Secondary', 'Success', 'Info', 'Warning', 'Help', 'Danger', 'Contrast']

    //literals
    wsMessage: string = 'no active downloads'
    serverLogs: string = 'no logs available'
    downloadComplete: string = 'Download completed successfully.'
    triggerDownloadApiSuccessResponse: string = 'Item added to download queue successfully.'
    downloadInfoIdentifier: string = '[download]'
    getInfo: string = 'fetching content metadata & thumbnail before downloading, this might take a while!'
    getChannel: string = 'fetching channel'
    getTitle: string = 'fetching title'
}

//export const wsApiUrl: string = 'ws://localhost:3000/ws/downloadstatus'
export enum Severity {
    secondary = 1,
    success = 2,
    info = 3,
    warning = 4,
    help = 5,
    danger = 6,
    contrast = 7
}