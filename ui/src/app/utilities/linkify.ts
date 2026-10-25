export function linkify(text: string): string {
    var urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
    return text.replace(urlRegex, function (url: string) {
        return '<a href="' + url + '">' + url + '</a>';
    });
}
