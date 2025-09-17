import { VideoData } from "./video-data";

export class CardGroupDataResponse {
    data: CardGroup[] = [];
    message: string = '';
    status: string = '';
}

export class CardGroup {
    id: number = -1;
    title: string = '';
    subtitle: string = '';
    uploader: string = '';
    item_count: number = -1;
    external_id: string = '';
    thumbnail: string = '';
}

export class SelectedCardGroup {
    groupInfo: CardGroup = new CardGroup()
    active_video: number = -1
    video_data: VideoData[] = []
}