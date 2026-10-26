// base URL the UI calls the API on: the configured one, or the origin the page was loaded from
// an unsubstituted '${baseUrl}' (build served without entrypoint.sh) counts as not configured
export function appBaseUrl(configured: string, origin: string = location.origin): string {
    const value = !configured || configured.startsWith('${') ? origin : configured;
    return value.replace(/\/$/, '');
}

// websocket URL for an http(s) base URL; only the scheme is swapped
export function toWsUrl(baseUrl: string): string {
    return baseUrl.replace(/^http/, 'ws').replace(/\/$/, '');
}
