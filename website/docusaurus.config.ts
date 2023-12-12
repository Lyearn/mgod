import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'mgod',
  tagline: 'Transform your MongoDB interactions in Go effortlessly with mgod',
  favicon: '/img/favicon.ico',

  // Set the production url of your site here
  url: 'https://lyearn.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/mgod/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'Lyearn',
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
      const tailwindCss = (await import("tailwindcss"));
      const autoprefixer = (await import("autoprefixer"));
    
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
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
        // {
        //   href: 'https://blog.lyearn.com',
        //   label: 'Blog',
        //   position: 'right',
        // },
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
