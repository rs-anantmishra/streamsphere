// Playwright reads Linux Mint (ID=linuxmint) as Ubuntu 20.04, which has no Chromium build. LMDE is Debian underneath
// (DEBIAN_CODENAME in /etc/os-release), so on LMDE 7 (Debian 13, trixie) Playwright is told the real platform.
// Every other system is left to Playwright's own detection
export function detectPlatformOverride(osRelease: string | undefined, arch: string): string | undefined {
    const fields = Object.fromEntries((osRelease ?? '').split('\n').filter(line => line.includes('=')).map(line => {
        const [key, ...rest] = line.split('=');
        return [key, rest.join('=').replace(/^"|"$/g, '')];
    }));
    const suffix = arch === 'x64' ? 'x64' : arch === 'arm64' ? 'arm64' : undefined;
    if (suffix && fields.ID === 'linuxmint' && fields.DEBIAN_CODENAME === 'trixie') { return `debian13-${suffix}` }
    return undefined;
}
