export type DownloadState = 'idle' | 'downloading' | 'saving' | 'completed' | 'failed';

export interface DownloadStatusEvent {
    download: string;
    video_url: string;
    state: DownloadState;
    video_id: number;
    title: string;
    channel: string;
    thumbnail: string;
}
