import { VideoData } from "./video-data";

export class PlaylistsDataResponse {
    data: PlaylistsInfo[] = [];
    message: string = '';
    status: string = '';
}

export class PlaylistsInfo {
    playlist_id: number = -1
    playlist_title: string = ''
    playlist_uploader: string = ''
    item_count: number = -1
    yt_playlist_id: string = ''
    thumbnail: string = ''
}

export class SelectedPlaylist {
    info: PlaylistsInfo = new PlaylistsInfo()
    active_video: number = -1
    video_data: VideoData[] = []
}
