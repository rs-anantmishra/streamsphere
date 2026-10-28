import { environment } from '../../environments/environment';

export interface PageLocation { protocol: string, hostname: string, port: string }
export interface Address { protocol: string, host: string, port: string }
export interface Settings { protocol: string, applicationIp: string, applicationPort: string, contentIp: string, contentPort: string }

// one setting; empty, or an unsubstituted '${...}' placeholder (build served without entrypoint.sh), is not set
const given = (value: string) => value && !value.startsWith('${') ? value : ''

// full address from the parts that were set; every part that was not set is the page's own.
// when a host is set the address is described completely, so a missing port means the protocol's default
// an address that is the page's own host and port always uses the protocol the page was opened with,
// one port speaks one protocol
export function addressUrl(address: Address, page: PageLocation = location): string {
    // not named after the ${...} placeholders of environment.ts: entrypoint.sh substitutes those in the whole bundle
    const pageScheme = page.protocol.replace(/:$/, '')
    const scheme = given(address.protocol) || pageScheme
    const host = given(address.host)
    const port = given(address.port) || (host ? '' : page.port)
    const hostname = host || page.hostname

    // no port given and no host set is the page's own port; a host without a port is the protocol's default port
    const defaultPort = (s: string) => s === 'https' ? '443' : '80'
    const pagePort = page.port || defaultPort(pageScheme)
    const ownPort = given(address.port) || (host ? defaultPort(scheme) : pagePort)
    const samePlace = hostname === page.hostname && ownPort === pagePort

    return (samePlace ? pageScheme : scheme) + '://' + hostname + (port ? ':' + port : '')
}

// where the UI calls the api: PROTOCOL://APPLICATION_IP:APPLICATION_PORT
export function apiAddress(settings: Settings, page: PageLocation = location): string {
    return addressUrl({ protocol: settings.protocol, host: settings.applicationIp, port: settings.applicationPort }, page)
}

// where content is served from: PROTOCOL://CONTENT_IP:CONTENT_PORT, the content ip being APPLICATION_IP when it is not set
export function contentAddress(settings: Settings, page: PageLocation = location): string {
    return addressUrl({ protocol: settings.protocol, host: given(settings.contentIp) || settings.applicationIp, port: settings.contentPort }, page)
}

// websocket URL for an http(s) base URL; only the scheme is swapped
export function toWsUrl(baseUrl: string): string {
    return baseUrl.replace(/^http/, 'ws').replace(/\/$/, '');
}

// a content path (media, thumbnail) on the content server
export function contentLink(path: string, base: string): string {
    if (!path || path.startsWith('./') || /^(https?:)?\/\//.test(path)) { return path }
    return base.replace(/\/$/, '') + (path.startsWith('/') ? path : '/' + path)
}

// rewrites the media_url and thumbnail of a response (a record, an array, or a record with a data array)
export function linkContent<T>(value: T, link: (path: string) => string): T {
    if (Array.isArray(value)) { return value.map(item => linkContent(item, link)) as T }
    if (value && typeof value === 'object') {
        const out: any = { ...value }
        for (const key of Object.keys(out)) {
            if ((key === 'media_url' || key === 'thumbnail') && typeof out[key] === 'string') { out[key] = link(out[key]) }
            else if (key === 'data') { out[key] = linkContent(out[key], link) }
        }
        return out
    }
    return value
}

// the api address of this deployment
export function apiBaseUrl(): string {
    return apiAddress(environment)
}

// puts the content host in front of the media and thumbnail paths the backend returns
export function withContentLinks<T>(value: T): T {
    return linkContent(value, path => contentLink(path, contentAddress(environment)))
}

// what the library and playlists pages say when the api cannot be reached at all
export function apiUnreachableMessage(address: string): string {
    return 'Cannot reach the Streamsphere API at ' + address + '. Check PROTOCOL, APPLICATION_IP and APPLICATION_PORT in your docker compose file.'
}
