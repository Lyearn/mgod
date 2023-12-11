import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'mgod',
  tagline: 'MongoDB ODM designed to work with Go modelsl',
  favicon: '/img/favicon.ico',

  // Set the production url of your site here
  url: 'https://aryan02420.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/mgod/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'aryan02420',
  projectName: 'mgod',
  trailingSlash: false, // https://github.com/slorber/trailing-slash-guide 

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          path: '../docs/',
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/aryan02420/mgod/tree/main/docs/',
          exclude: ['**/README.md'],
          breadcrumbs: false,
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/social-card.jpg', // twitter:image
    colorMode: {
      disableSwitch: true,
    },
    navbar: {
      title: 'mgod',
      logo: {
        alt: 'mgod',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          label: 'Docs',
          position: 'right',
        },
        { 
          href: 'https://pkg.go.dev/github.com/Lyearn/mgod',
          label: 'API',
          position: 'right',
        },
        // {
        //   href: 'https://blog.lyearn.com',
        //   label: 'Blog',
        //   position: 'right',
        // },
        {
          href: 'https://github.com/aryan02420/mgod',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    // footer: {
    //   style: 'dark',
    //   links: [
    //     {
    //       title: 'Docs',
    //       items: [
    //         {
    //           label: 'Tutorial',
    //           to: '/docs/intro',
    //         },
    //       ],
    //     },
    //     {
    //       title: 'Community',
    //       items: [
    //         {
    //           label: 'Stack Overflow',
    //           href: 'https://stackoverflow.com/questions/tagged/docusaurus',
    //         },
    //         {
    //           label: 'Discord',
    //           href: 'https://discordapp.com/invite/docusaurus',
    //         },
    //         {
    //           label: 'Twitter',
    //           href: 'https://twitter.com/docusaurus',
    //         },
    //       ],
    //     },
    //     {
    //       title: 'More',
    //       items: [
    //         {
    //           label: 'Blog',
    //           href: 'https://blog.lyearn.com',
    //         },
    //         {
    //           label: 'GitHub',
    //           href: 'https://github.com/Lyearn/mgod',
    //         },
    //       ],
    //     },
    //   ],
    //   copyright: `Copyright © ${new Date().getFullYear()} Lyearn, Inc.`,
    // },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
