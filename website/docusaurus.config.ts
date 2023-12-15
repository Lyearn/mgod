import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'mgod',
  tagline: 'Transform your MongoDB interactions in Go effortlessly with mgod',
  favicon: '/img/favicon.ico',

  url: 'https://lyearn.github.io',
  baseUrl: '/mgod/',

  // GitHub pages deployment config.
  organizationName: 'Lyearn',
  projectName: 'mgod',
  trailingSlash: false, // https://github.com/slorber/trailing-slash-guide 

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // you can use this field to set useful metadata like html lang
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
          // editUrl: 'https://github.com/Lyearn/mgod/tree/main/docs/',
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

  plugins: [
    // https://www.swyx.io/tailwind-docusaurus-2022
    async function tailwindPlugin(context, options) {
      const tailwindCssNesting = (await import("tailwindcss/nesting"));
      const tailwindCss = (await import("tailwindcss"));
      const autoprefixer = (await import("autoprefixer"));
    
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          postcssOptions.plugins.push(tailwindCssNesting);
          postcssOptions.plugins.push(tailwindCss);
          postcssOptions.plugins.push(autoprefixer);
          return postcssOptions;
        },
      };
    },
  ],

  themeConfig: {
    image: 'img/social-card.jpg', // twitter:image
    colorMode: {
      disableSwitch: true,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'mgod',
      logo: {
        alt: 'mgod',
        src: 'img/logo-with-text.svg',
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
        {
          'aria-label': 'GitHub',
          className: 'navbar-github-link',
          href: 'https://github.com/Lyearn/mgod',
          position: 'right',
        },
      ],
    },
    tableOfContents: {
      minHeadingLevel: 2,
      maxHeadingLevel: 2,
    },
    footer: {
      links: [
        {
          label: 'Github Discussions',
          href: 'https://github.com/Lyearn/mgod/discussions',
          className: 'footer-github-discussions-link',
        },
      ],
      // copyright: `Copyright © ${new Date().getFullYear()} Lyearn, Inc.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
