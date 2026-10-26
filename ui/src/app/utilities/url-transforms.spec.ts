import { addressUrl, apiAddress, contentAddress, contentLink, linkContent, toWsUrl } from './url-transforms';

describe('toWsUrl', () => {
    it('turns http into ws and keeps the port', () => {
        expect(toWsUrl('http://192.168.1.11:1282')).toBe('ws://192.168.1.11:1282');
    });

    it('turns https into wss', () => {
        expect(toWsUrl('https://streamsphere.example.com')).toBe('wss://streamsphere.example.com');
    });

    it('keeps an https port', () => {
        expect(toWsUrl('https://streamsphere.example.com:8443')).toBe('wss://streamsphere.example.com:8443');
    });

    it('leaves "api" and "http" inside the host alone', () => {
        expect(toWsUrl('https://api.example.com')).toBe('wss://api.example.com');
        expect(toWsUrl('http://httpbin.local:3000')).toBe('ws://httpbin.local:3000');
    });

    it('keeps a base path', () => {
        expect(toWsUrl('https://example.com/streamsphere')).toBe('wss://example.com/streamsphere');
    });

    it('drops a trailing slash', () => {
        expect(toWsUrl('http://localhost:3000/')).toBe('ws://localhost:3000');
    });
});

describe('addressUrl', () => {
    const lan = { protocol: 'http:', hostname: '192.168.1.11', port: '1282' };
    const secure = { protocol: 'https:', hostname: 'ss.example.com', port: '' };

    it('uses the page host with the port', () => {
        expect(addressUrl({ protocol: '', host: '', port: '1288' }, lan)).toBe('http://192.168.1.11:1288');
        expect(addressUrl({ protocol: '', host: '', port: '1288' }, secure)).toBe('https://ss.example.com:1288');
    });

    it('uses the host and the port that were set', () => {
        expect(addressUrl({ protocol: '', host: '10.0.0.5', port: '3500' }, lan)).toBe('http://10.0.0.5:3500');
    });

    it('uses the protocol that was set, else the page protocol', () => {
        expect(addressUrl({ protocol: 'https', host: 'ss.example.com', port: '443' }, lan)).toBe('https://ss.example.com:443');
        expect(addressUrl({ protocol: '', host: 'ss.example.com', port: '443' }, secure)).toBe('https://ss.example.com:443');
    });

    it('is the page address when nothing was set', () => {
        expect(addressUrl({ protocol: '', host: '', port: '' }, lan)).toBe('http://192.168.1.11:1282');
        expect(addressUrl({ protocol: '', host: '', port: '' }, secure)).toBe('https://ss.example.com');
    });

    it('the page\'s own host and port always use the protocol the page was opened with', () => {
        // https asked for on an http page, and the other way round
        expect(addressUrl({ protocol: 'https', host: '', port: '1282' }, lan)).toBe('http://192.168.1.11:1282');
        expect(addressUrl({ protocol: 'http', host: 'ss.example.com', port: '443' }, secure)).toBe('https://ss.example.com:443');
        expect(addressUrl({ protocol: 'http', host: '', port: '' }, secure)).toBe('https://ss.example.com');
    });

    it('the protocol is still honored for another host', () => {
        expect(addressUrl({ protocol: 'https', host: 'api.example.com', port: '443' }, lan)).toBe('https://api.example.com:443');
        expect(addressUrl({ protocol: 'http', host: 'content.example.com', port: '8443' }, secure)).toBe('http://content.example.com:8443');
    });

    it('the protocol is still honored for another port on the page host', () => {
        expect(addressUrl({ protocol: 'https', host: '', port: '1288' }, lan)).toBe('https://192.168.1.11:1288');
    });

    it('treats unsubstituted placeholders as not set', () => {
        expect(addressUrl({ protocol: '${protocol}', host: '${applicationIp}', port: '${applicationPort}' }, lan)).toBe('http://192.168.1.11:1282');
    });
});

describe('apiAddress and contentAddress', () => {
    const lan = { protocol: 'http:', hostname: '192.168.1.11', port: '1282' };
    const ports = { protocol: '', applicationIp: '', applicationPort: '1282', contentIp: '', contentPort: '1288' };

    it('both ports only: the page host on each port', () => {
        expect(apiAddress(ports, lan)).toBe('http://192.168.1.11:1282');
        expect(contentAddress(ports, lan)).toBe('http://192.168.1.11:1288');
    });

    it('APPLICATION_IP: the API there, and the content there too', () => {
        const s = { ...ports, applicationIp: '192.168.1.15' };
        expect(apiAddress(s, lan)).toBe('http://192.168.1.15:1282');
        expect(contentAddress(s, lan)).toBe('http://192.168.1.15:1288');
    });

    it('CONTENT_IP only: the API stays on the page host', () => {
        const s = { ...ports, contentIp: 'content.example.com' };
        expect(apiAddress(s, lan)).toBe('http://192.168.1.11:1282');
        expect(contentAddress(s, lan)).toBe('http://content.example.com:1288');
    });

    it('both IPs: each address uses its own', () => {
        const s = { ...ports, applicationIp: 'ss.example.com', contentIp: 'content.example.com' };
        expect(apiAddress(s, lan)).toBe('http://ss.example.com:1282');
        expect(contentAddress(s, lan)).toBe('http://content.example.com:1288');
    });

    it('PROTOCOL applies to every address except the page\'s own', () => {
        const s = { ...ports, protocol: 'https' };
        expect(apiAddress(s, lan)).toBe('http://192.168.1.11:1282');
        expect(contentAddress(s, lan)).toBe('https://192.168.1.11:1288');
    });
});

describe('contentLink', () => {
    it('puts the base in front of a content path', () => {
        expect(contentLink('/youtube.com/a/x.webm', 'https://content.example.com:8443')).toBe('https://content.example.com:8443/youtube.com/a/x.webm');
    });

    it('adds the missing slash', () => {
        expect(contentLink('youtube.com/a/x.webm', 'http://10.0.0.5:3500')).toBe('http://10.0.0.5:3500/youtube.com/a/x.webm');
    });

    it('leaves empty values and local assets alone', () => {
        expect(contentLink('', 'http://h')).toBe('');
        expect(contentLink('./assets/noimage.png', 'http://h')).toBe('./assets/noimage.png');
    });

    it('leaves an already absolute url alone', () => {
        expect(contentLink('https://other.example.com/x.webm', 'http://h')).toBe('https://other.example.com/x.webm');
    });
});

describe('linkContent', () => {
    const link = (path: string) => 'https://c' + path;

    it('rewrites media_url and thumbnail of a record', () => {
        expect(linkContent({ title: 't', media_url: '/a.webm', thumbnail: '/a.webp' }, link))
            .toEqual({ title: 't', media_url: 'https://c/a.webm', thumbnail: 'https://c/a.webp' });
    });

    it('rewrites every record of an array', () => {
        expect(linkContent([{ thumbnail: '/1' }, { thumbnail: '/2' }], link)).toEqual([{ thumbnail: 'https://c/1' }, { thumbnail: 'https://c/2' }]);
    });

    it('rewrites the records in a data array', () => {
        expect(linkContent({ status: 'ok', data: [{ media_url: '/a' }] }, link)).toEqual({ status: 'ok', data: [{ media_url: 'https://c/a' }] });
    });

    it('leaves other fields and non-string values alone', () => {
        expect(linkContent({ title: '/keep', thumbnail: null }, link)).toEqual({ title: '/keep', thumbnail: null });
    });
});
