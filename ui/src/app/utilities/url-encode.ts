

export class UrlEncode {
    //list of special chars for reference
    chars = {
        "specialChars": [
            { "char": "%", "encoded": "%25" }, 
            { "char": "#", "encoded": "%23" }, 
            //{ "char": "!", "encoded": "%21" }, 
            //{ "char": "\"", "encoded": "%22" }, 
            //{ "char": "$", "encoded": "%24" }, 
            //{ "char": "&", "encoded": "%26" }, 
            //{ "char": "'", "encoded": "%27" }, 
            //{ "char": "(", "encoded": "%28" }, 
            //{ "char": ")", "encoded": "%29" }, 
            //{ "char": "*", "encoded": "%2A" }, 
            //{ "char": "+", "encoded": "%2B" }, 
            //{ "char": ",", "encoded": "%2C" }, 
            //{ "char": "-", "encoded": "%2D" },
            // { "char": ".", "encoded": "%2E" }
        ]
    };

    encodedUrl(url: string): string {
        for (let k = 0; k < this.chars.specialChars.length; k++) {

            if (url.indexOf(this.chars.specialChars[k].char) > 0) {
                url = url.replace(this.chars.specialChars[k].char, this.chars.specialChars[k].encoded)
            }
        }
        return url
    }
}
