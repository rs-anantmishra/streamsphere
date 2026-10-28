import { PORTS } from './stack';

// one row per way a user can configure the app. Plain http pages are opened at localhost on the published app port,
// https pages go through nginx at streamsphere.localhost (the content server is also reachable as content.localhost)
export type Scenario = {
    id: string;
    env: Record<string, string>;
    page: 'http' | 'https';
    // what the container prints at startup
    prints?: string[];
    // the address the page must use for the api, and for the content
    api?: string;
    content?: string;
    // the container must refuse to start, naming this variable
    stops?: string;
    // a wrong setting on the content address: the api works (an address that is the page's own always uses the page's
    // protocol), but the content is on another port or host that does not speak the protocol that was asked for
    contentBroken?: boolean;
    // nothing listens at the api address: the pages must say so instead of spinning
    apiUnreachable?: boolean;
    // also run a real download (needs the network)
    download?: boolean;
};

const http = (host: string, port: number) => `http://${host}:${port}`;
const https = (host: string, port: number) => `https://${host}:${port}`;
const ports = { APPLICATION_PORT: String(PORTS.app), CONTENT_PORT: String(PORTS.content) };
const tlsPorts = { APPLICATION_PORT: String(PORTS.tlsApp), CONTENT_PORT: String(PORTS.tlsApp) };

export const scenarios: Scenario[] = [
    { id: 'p1 ports only', page: 'http', env: { ...ports }, api: http('localhost', PORTS.app), content: http('localhost', PORTS.content), download: true },
    { id: 'p2 legacy HOST_IP', page: 'http', env: { HOST_IP: '127.0.0.1', ...ports }, prints: ['HOST_IP is deprecated, using it as APPLICATION_IP'], api: http('127.0.0.1', PORTS.app), content: http('127.0.0.1', PORTS.content) },
    { id: 'p3 APPLICATION_IP', page: 'http', env: { APPLICATION_IP: '127.0.0.1', ...ports }, api: http('127.0.0.1', PORTS.app), content: http('127.0.0.1', PORTS.content) },
    { id: 'p4 CONTENT_IP only', page: 'http', env: { CONTENT_IP: '127.0.0.1', ...ports }, api: http('localhost', PORTS.app), content: http('127.0.0.1', PORTS.content) },
    { id: 'p5 both IPs', page: 'http', env: { APPLICATION_IP: '127.0.0.1', CONTENT_IP: 'localhost', ...ports }, api: http('127.0.0.1', PORTS.app), content: http('localhost', PORTS.content) },
    { id: 'p6 HOST_IP and APPLICATION_IP', page: 'http', env: { HOST_IP: '10.9.9.9', APPLICATION_IP: '127.0.0.1', ...ports }, prints: ['HOST_IP is deprecated and ignored, APPLICATION_IP is set'], api: http('127.0.0.1', PORTS.app), content: http('127.0.0.1', PORTS.content) },
    { id: 'p7 HOST_IP and CONTENT_IP', page: 'http', env: { HOST_IP: '127.0.0.1', CONTENT_IP: 'localhost', ...ports }, prints: ['HOST_IP is deprecated, using it as APPLICATION_IP'], api: http('127.0.0.1', PORTS.app), content: http('localhost', PORTS.content) },
    { id: 'p8 PROTOCOL=http', page: 'http', env: { PROTOCOL: 'http', ...ports }, prints: ['protocol: http'], api: http('localhost', PORTS.app), content: http('localhost', PORTS.content) },
    { id: 'p9 PROTOCOL=https on an http page (wrong)', page: 'http', env: { PROTOCOL: 'https', ...ports }, api: http('localhost', PORTS.app), content: https('localhost', PORTS.content), contentBroken: true },
    { id: 'p10 wrong application port (nothing listens there)', page: 'http', env: { APPLICATION_PORT: '1299', CONTENT_PORT: String(PORTS.content) }, api: http('localhost', 1299), apiUnreachable: true },

    { id: 'h1 content on its own name', page: 'https', env: { ...tlsPorts, CONTENT_IP: 'content.localhost' }, api: https('streamsphere.localhost', PORTS.tlsApp), content: https('content.localhost', PORTS.tlsApp), download: true },
    { id: 'h2 content on the page host, own port', page: 'https', env: { APPLICATION_PORT: String(PORTS.tlsApp), CONTENT_PORT: String(PORTS.tlsContent) }, api: https('streamsphere.localhost', PORTS.tlsApp), content: https('streamsphere.localhost', PORTS.tlsContent) },
    { id: 'h3 PROTOCOL=https', page: 'https', env: { PROTOCOL: 'https', ...tlsPorts, CONTENT_IP: 'content.localhost' }, prints: ['protocol: https'], api: https('streamsphere.localhost', PORTS.tlsApp), content: https('content.localhost', PORTS.tlsApp) },
    { id: 'h4 legacy HOST_IP', page: 'https', env: { HOST_IP: 'streamsphere.localhost', APPLICATION_PORT: String(PORTS.tlsApp), CONTENT_PORT: String(PORTS.tlsContent) }, prints: ['HOST_IP is deprecated, using it as APPLICATION_IP'], api: https('streamsphere.localhost', PORTS.tlsApp), content: https('streamsphere.localhost', PORTS.tlsContent) },
    { id: 'h5 both IPs', page: 'https', env: { APPLICATION_IP: 'streamsphere.localhost', CONTENT_IP: 'content.localhost', ...tlsPorts }, api: https('streamsphere.localhost', PORTS.tlsApp), content: https('content.localhost', PORTS.tlsApp) },
    { id: 'h6 PROTOCOL=http on an https page (wrong)', page: 'https', env: { PROTOCOL: 'http', ...tlsPorts, CONTENT_IP: 'content.localhost' }, api: https('streamsphere.localhost', PORTS.tlsApp), content: http('content.localhost', PORTS.tlsApp), contentBroken: true },

    { id: 'f1 no ports', page: 'http', env: { PROTOCOL: 'http' }, stops: 'APPLICATION_PORT' },
    { id: 'f2 no application port', page: 'http', env: { CONTENT_PORT: String(PORTS.content) }, stops: 'APPLICATION_PORT' },
    { id: 'f3 no content port', page: 'http', env: { APPLICATION_PORT: String(PORTS.app) }, stops: 'CONTENT_PORT' },
];

export const pageUrl = (s: Scenario) => s.page === 'http' ? http('localhost', PORTS.app) : https('streamsphere.localhost', PORTS.tlsApp);
