import { ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from '@angular/core';
import { provideRouter, withHashLocation } from '@angular/router';
import { provideAnimations } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';


import { routes } from './app.routes';
import { provideHttpClient, withFetch } from '@angular/common/http';

export const appConfig: ApplicationConfig = {
  providers: [provideAnimationsAsync(), provideZoneChangeDetection({ eventCoalescing: true }), [provideRouter(routes, withHashLocation())],
  importProvidersFrom(BrowserModule), importProvidersFrom(BrowserAnimationsModule), provideAnimations(), provideHttpClient(withFetch())
  ]
};

