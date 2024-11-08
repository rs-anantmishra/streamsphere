export class ContentSearchResponse {
    data: ContentSearch[] = []
    message: string = '';
    status: string = '';
}

export class ContentSearch {
    video_id: number = 0;
    channel: string = '';
    title: string = '';
}
