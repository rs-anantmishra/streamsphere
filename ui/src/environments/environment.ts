export const environment = {
    DefaultLanguage: "en",
    production: true,
    development: false,
    environmentName:"PROD",
    // filled in by entrypoint.sh from the docker compose environment; empty means "same as the page"
    protocol: '${protocol}',
    applicationIp: '${applicationIp}',
    applicationPort: '${applicationPort}',
    contentIp: '${contentIp}',
    contentPort: '${contentPort}',
};

// environment.ts
export interface Environment {
    protocol: string;
    applicationIp: string;
    applicationPort: string;
    contentIp: string;
    contentPort: string;
  }