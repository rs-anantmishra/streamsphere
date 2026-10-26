import { test } from 'node:test';
import * as assert from 'node:assert/strict';
import { detectPlatformOverride } from './platform.ts';

const lmde = 'NAME="LMDE"\nVERSION_ID="7"\nID=linuxmint\nID_LIKE=debian\nVERSION_CODENAME=gigi\nDEBIAN_CODENAME=trixie\n';
const ubuntu = 'NAME="Ubuntu"\nVERSION_ID="24.04"\nID=ubuntu\nID_LIKE=debian\n';
const mint = 'NAME="Linux Mint"\nVERSION_ID="22"\nID=linuxmint\nID_LIKE="ubuntu debian"\nUBUNTU_CODENAME=noble\n';

test('LMDE 7 is Debian 13 underneath', () => {
    assert.equal(detectPlatformOverride(lmde, 'x64'), 'debian13-x64');
    assert.equal(detectPlatformOverride(lmde, 'arm64'), 'debian13-arm64');
});

test('other systems are left to Playwright', () => {
    assert.equal(detectPlatformOverride(ubuntu, 'x64'), undefined);
    assert.equal(detectPlatformOverride(mint, 'x64'), undefined);
    assert.equal(detectPlatformOverride('', 'x64'), undefined);
    assert.equal(detectPlatformOverride(undefined, 'x64'), undefined);
});

test('an architecture Playwright has no Debian build for is left alone', () => {
    assert.equal(detectPlatformOverride(lmde, 'ia32'), undefined);
});
