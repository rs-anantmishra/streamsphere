import { ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from '@angular/core';
import { provideRouter, withHashLocation } from '@angular/router';
import { provideAnimations } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';
import { definePreset } from '@primeng/themes';
import Aura from '@primeng/themes/aura';


import { routes } from './app.routes';
import { provideHttpClient, withFetch } from '@angular/common/http';

// Aura defaults to emerald in v18; keep the blue primary from the old aura-*-blue themes
const AuraBlue = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{blue.50}', 100: '{blue.100}', 200: '{blue.200}', 300: '{blue.300}', 400: '{blue.400}',
      500: '{blue.500}', 600: '{blue.600}', 700: '{blue.700}', 800: '{blue.800}', 900: '{blue.900}', 950: '{blue.950}'
    }
  }
});

export const appConfig: ApplicationConfig = {
  providers: [provideAnimationsAsync(), provideZoneChangeDetection({ eventCoalescing: true }), [provideRouter(routes, withHashLocation())],
  importProvidersFrom(BrowserModule), importProvidersFrom(BrowserAnimationsModule), provideAnimations(), provideHttpClient(withFetch()),
  providePrimeNG({ theme: { preset: AuraBlue, options: { darkModeSelector: '.app-dark' } } })
  ]
};
