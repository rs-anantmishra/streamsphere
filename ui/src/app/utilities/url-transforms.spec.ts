import { appBaseUrl, toWsUrl } from './url-transforms';

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

describe('appBaseUrl', () => {
    it('uses the origin when nothing is configured', () => {
        expect(appBaseUrl('', 'https://ss.example.com')).toBe('https://ss.example.com');
    });

    it('uses the origin when the placeholder was never substituted', () => {
        expect(appBaseUrl('${baseUrl}', 'http://localhost:3000')).toBe('http://localhost:3000');
    });

    it('uses the configured url when set', () => {
        expect(appBaseUrl('https://api.example.com', 'https://ss.example.com')).toBe('https://api.example.com');
    });

    it('drops a trailing slash', () => {
        expect(appBaseUrl('https://api.example.com/', 'http://x')).toBe('https://api.example.com');
    });
});
