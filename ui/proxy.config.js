const PROXY_CONFIG = {
    "/youtube.com": {
        "target": "http://localhost:3500/youtube.com",
        "secure": false,
        "changeOrigin": true,
        "logLevel": "debug"
    },
    "/api": {
        "target": "http://localhost:3000",
        "secure": false,
        "changeOrigin": true,
        "logLevel": "debug"
    },
};

module.exports = PROXY_CONFIG