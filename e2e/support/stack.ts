// the docker + nginx stack the end-to-end tests run against. Everything lives in one temp folder and one compose project
// (streamsphere-e2e) so nothing of a running Streamsphere is touched. Needs: docker, nginx, openssl, and the two images:
//   docker build -f Dockerfile.amd64 -t streamsphere/streamsphere:local-https .   (from the repo root)
import { spawnSync } from 'node:child_process';
import * as fs from 'node:fs';
import * as os from 'node:os';
import * as path from 'node:path';

export const PORTS = {
    app: Number(process.env.E2E_APP_PORT ?? 1282),
    content: Number(process.env.E2E_CONTENT_PORT ?? 1288),
    tlsApp: Number(process.env.E2E_TLS_APP_PORT ?? 8443),
    tlsContent: Number(process.env.E2E_TLS_CONTENT_PORT ?? 8444),
    bare: Number(process.env.E2E_BARE_PORT ?? 1292),
    boundContent: Number(process.env.E2E_BOUND_CONTENT_PORT ?? 1293),
};
export const IMAGE = process.env.STREAMSPHERE_IMAGE ?? 'streamsphere/streamsphere:local-https';
export const CONTENT_IMAGE = process.env.STREAMSPHERE_CONTENT_IMAGE ?? 'streamsphere/streamsphere-content:0.1.3';
export const OLD_IMAGE = process.env.STREAMSPHERE_OLD_IMAGE ?? 'streamsphere/streamsphere:0.1.17';
export const NGINX = process.env.NGINX_PATH ?? (fs.existsSync('/usr/sbin/nginx') ? '/usr/sbin/nginx' : 'nginx');

const PROJECT = 'streamsphere-e2e';
export const APP = `${PROJECT}-app`;
export const VOLUMES = { db: `${PROJECT}_db-data`, content: `${PROJECT}_content-data` };
export const CONTENT = `${PROJECT}-content`;

export const workDir = () => {
    if (!process.env.E2E_DIR) { throw new Error('E2E_DIR is not set: the global setup did not run') }
    return process.env.E2E_DIR;
};

export function run(cmd: string, args: string[], opts: { cwd?: string, allowFail?: boolean } = {}) {
    const r = spawnSync(cmd, args, { cwd: opts.cwd, encoding: 'utf8' });
    if (r.status !== 0 && !opts.allowFail) { throw new Error(`${cmd} ${args.join(' ')} failed (${r.status}): ${r.stderr || r.stdout}`) }
    return { status: r.status ?? -1, out: (r.stdout ?? '') + (r.stderr ?? '') };
}

const compose = (args: string[], allowFail = false) => run('docker', ['compose', '-p', PROJECT, '-f', path.join(workDir(), 'docker-compose.yaml'), ...args], { cwd: workDir(), allowFail });

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms));

function portInUse(port: number) {
    const r = run('ss', ['-ltn'], { allowFail: true });
    return new RegExp(`:${port}\\s`).test(r.out);
}

export async function setUp() {
    for (const image of [IMAGE, CONTENT_IMAGE]) {
        if (run('docker', ['image', 'inspect', image], { allowFail: true }).status !== 0) {
            throw new Error(`image ${image} is missing. Build it first: docker build -f Dockerfile.amd64 -t ${IMAGE} . (from the repo root)`);
        }
    }
    for (const port of Object.values(PORTS)) {
        if (portInUse(port)) { throw new Error(`port ${port} is in use. Stop what is using it, or set E2E_APP_PORT, E2E_CONTENT_PORT, E2E_TLS_APP_PORT and E2E_TLS_CONTENT_PORT`) }
    }

    const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'streamsphere-e2e-'));
    process.env.E2E_DIR = dir;
    fs.writeFileSync(path.join(dir, 'scenario.env'), '');

    // one self-signed certificate for both names; the browser contexts ignore https errors
    run('openssl', ['req', '-x509', '-newkey', 'rsa:2048', '-nodes', '-days', '2', '-subj', '/CN=streamsphere.localhost',
        '-addext', 'subjectAltName=DNS:streamsphere.localhost,DNS:content.localhost',
        '-keyout', path.join(dir, 'key.pem'), '-out', path.join(dir, 'cert.pem')]);

    fs.writeFileSync(path.join(dir, 'docker-compose.yaml'), `services:
  streamsphere:
    image: ${IMAGE}
    container_name: ${APP}
    ports:
      - ${PORTS.app}:3000
    env_file: scenario.env
    volumes:
      - db-data:/app/database/db
      - content-data:/app/content
  content:
    image: ${CONTENT_IMAGE}
    container_name: ${CONTENT}
    ports:
      - ${PORTS.content}:3500
    volumes:
      - content-data:/content
volumes:
  db-data:
  content-data:
`);

    const cert = path.join(dir, 'cert.pem'), key = path.join(dir, 'key.pem');
    const tls = (listen: number, name: string, extra: string, target: number) => `
    server {
        listen ${listen} ssl; server_name ${name};
        ssl_certificate ${cert}; ssl_certificate_key ${key};
        location / { proxy_pass http://127.0.0.1:${target}; proxy_http_version 1.1; proxy_set_header Host $host; ${extra} }
    }`;
    const upgrade = 'proxy_set_header Upgrade $http_upgrade; proxy_set_header Connection $connection_upgrade;';
    fs.mkdirSync(path.join(dir, 'tmp'));
    fs.writeFileSync(path.join(dir, 'nginx.conf'), `pid ${dir}/nginx.pid;
error_log ${dir}/nginx-error.log;
events {}
http {
    access_log ${dir}/nginx-access.log;
    client_body_temp_path ${dir}/tmp/body; proxy_temp_path ${dir}/tmp/proxy; fastcgi_temp_path ${dir}/tmp/fcgi; uwsgi_temp_path ${dir}/tmp/uwsgi; scgi_temp_path ${dir}/tmp/scgi;
    map $http_upgrade $connection_upgrade { default upgrade; '' close; }
${tls(PORTS.tlsApp, 'streamsphere.localhost', upgrade, PORTS.app)}
${tls(PORTS.tlsApp, 'content.localhost', '', PORTS.content)}
${tls(PORTS.tlsContent, 'streamsphere.localhost', '', PORTS.content)}
}
`);

    compose(['up', '-d', 'content']);
    run(NGINX, ['-c', path.join(dir, 'nginx.conf'), '-e', path.join(dir, 'nginx-error.log')]);
}

export async function tearDown() {
    const dir = process.env.E2E_DIR;
    if (!dir) { return }
    removeBound();
    compose(['down', '-v'], true);
    const pid = path.join(dir, 'nginx.pid');
    if (fs.existsSync(pid)) { run('kill', [fs.readFileSync(pid, 'utf8').trim()], { allowFail: true }) }
    await sleep(500);
    fs.rmSync(dir, { recursive: true, force: true });
}

export type Started = { state: string, exitCode: number, logs: string };

// recreates the app container with the given environment; resolves once it is serving or has exited
export async function startApp(env: Record<string, string>): Promise<Started> {
    fs.writeFileSync(path.join(workDir(), 'scenario.env'), Object.entries(env).map(([k, v]) => `${k}=${v}`).join('\n') + '\n');
    compose(['up', '-d', '--force-recreate', 'streamsphere']);
    for (let i = 0; i < 60; i++) {
        const state = run('docker', ['inspect', '-f', '{{.State.Status}}', APP], { allowFail: true }).out.trim();
        if (state === 'exited' || logs().includes('Connected to database')) { break }
        await sleep(1000);
    }
    return {
        state: run('docker', ['inspect', '-f', '{{.State.Status}}', APP], { allowFail: true }).out.trim(),
        exitCode: Number(run('docker', ['inspect', '-f', '{{.State.ExitCode}}', APP], { allowFail: true }).out.trim()),
        logs: logs(),
    };
}

export const logs = () => run('docker', ['logs', APP], { allowFail: true }).out;

export function stopApp() { run('docker', ['stop', APP], { allowFail: true }) }
export function startAppAgain() { run('docker', ['start', APP]) }
export function copyDb(direction: 'out' | 'in') {
    const local = path.join(workDir(), 'streamsphere.db');
    if (direction === 'out') { run('docker', ['cp', `${APP}:/app/database/db/streamsphere.db`, local]) }
    else { run('docker', ['cp', local, `${APP}:/app/database/db/streamsphere.db`]) }
    return local;
}

// runs an image next to the stack, on the stack's volumes, with its own entrypoint (or the image's own) and environment
export function runContainer(name: string, image: string, hostPort: number, env: Record<string, string>, entrypoint?: string[]) {
    run('docker', ['rm', '-f', name], { allowFail: true });
    const args = ['run', '-d', '--name', name, '-p', `${hostPort}:3000`,
        '-v', `${VOLUMES.db}:/app/database/db`, '-v', `${VOLUMES.content}:/app/content`,
        ...Object.entries(env).flatMap(([k, v]) => ['-e', `${k}=${v}`])];
    if (entrypoint) { args.push('--entrypoint', entrypoint[0]) }
    args.push(image, ...(entrypoint ? entrypoint.slice(1) : []));
    run('docker', args);
}
export const removeContainer = (name: string) => run('docker', ['rm', '-f', name], { allowFail: true });
export const containerLogs = (name: string) => run('docker', ['logs', name], { allowFail: true }).out;

// a fresh install on host folders (bind mounts) instead of named volumes, as some users run it (#47): an app and a
// content server sharing one content folder, next to the stack. A bind mount hides what the image has at the mount point
export const BOUND = { app: `${PROJECT}-bound`, content: `${PROJECT}-bound-content` };
const boundRoot = () => path.join(workDir(), 'bound');

// empty host folders for a new install; anything a previous one left is removed first
export function boundFolders() {
    removeBound();
    const folders = { db: path.join(boundRoot(), 'db'), content: path.join(boundRoot(), 'content') };
    fs.mkdirSync(folders.db, { recursive: true });
    fs.mkdirSync(folders.content, { recursive: true });
    return folders;
}

// starts both containers on the folders; resolves once the app is serving or has exited
export async function startBound(folders: { db: string, content: string }, opts: { readOnlyDb?: boolean } = {}): Promise<Started> {
    run('docker', ['run', '-d', '--name', BOUND.content, '-p', `${PORTS.boundContent}:3500`, '-v', `${folders.content}:/content`, CONTENT_IMAGE]);
    run('docker', ['run', '-d', '--name', BOUND.app, '-p', `${PORTS.bare}:3000`,
        '-v', `${folders.db}:/app/database/db${opts.readOnlyDb ? ':ro' : ''}`, '-v', `${folders.content}:/app/content`,
        '-e', `APPLICATION_PORT=${PORTS.bare}`, '-e', `CONTENT_PORT=${PORTS.boundContent}`, IMAGE]);
    const state = () => run('docker', ['inspect', '-f', '{{.State.Status}}', BOUND.app], { allowFail: true }).out.trim();
    // a new container updates yt-dlp before it serves, which can take a while
    for (let i = 0; i < 180; i++) {
        if (state() === 'exited' || containerLogs(BOUND.app).includes('Connected to database')) { break }
        await sleep(1000);
    }
    return {
        state: state(),
        exitCode: Number(run('docker', ['inspect', '-f', '{{.State.ExitCode}}', BOUND.app], { allowFail: true }).out.trim()),
        logs: containerLogs(BOUND.app),
    };
}

// removes both containers and their files; the containers write as root, so the files go through a container too
export function removeBound() {
    removeContainer(BOUND.app);
    removeContainer(BOUND.content);
    const root = process.env.E2E_DIR && path.join(process.env.E2E_DIR, 'bound');
    if (root && fs.existsSync(root)) {
        run('docker', ['run', '--rm', '-v', `${root}:/bound`, '--entrypoint', 'rm', IMAGE, '-rf', '/bound/db', '/bound/content'], { allowFail: true });
        fs.rmSync(root, { recursive: true, force: true });
    }
}
