import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'Intext Documentation',
  url: 'https://docs.codekeg.dev',
  baseUrl: '/',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
        }
      }
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    colorMode: {
      respectPrefersColorScheme: true
    },
    navbar: {
      logo: {
        src: 'img/itx-logo.png',
        alt: 'Intext\'s Logo'
      },
      items: [
        {
          type: 'docSidebar',
          position: 'left',
          label: 'Getting Started',
          sidebarId: 'gettingStartedSidebar'
        },
        {
          type: 'docSidebar',
          position: 'left',
          label: 'Functions',
          sidebarId: 'functionSidebar'
        },
        {
          type: 'docSidebar',
          position: 'left',
          label: 'Concepts',
          sidebarId: 'conceptSidebar'
        },
        {
          type: 'docSidebar',
          position: 'left',
          label: 'Syntax',
          sidebarId: 'syntaxSidebar'
        },
        {
          href: 'https://github.com/Elemnto56/Intextual_File_Language',
          label: 'GitHub',
          position: 'right',
        },
        {
          type: 'dropdown',
          position: 'right',
          label: 'External Intext Sites',
          items: [
            {
              href: 'https://docs.codekeg.dev/playground/',
              label: 'Intext Playground'
            }
          ]
        }
      ],
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
