export const environment = {
    DefaultLanguage: "en",
    production: true,
    development: false,
    environmentName:"PROD",
    baseUrl: '${baseUrl}',
    apiUrl: '${apiUrl}',
    contentUrl: '${contentUrl}',    
};

// environment.ts
export interface Environment {
    apiUrl: string;
    baseUrl: string;
    contentUrl: string;
  }