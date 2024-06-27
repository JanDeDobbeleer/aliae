const path = require('path');

module.exports = {
  title: 'aliae',
  tagline: 'Cross platform shell management.',
  url: 'https://aliae.dev',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'jandedobbeleer',
  projectName: 'aliae',
  onBrokenLinks: 'ignore',
  plugins: [
    path.resolve(__dirname, 'plugins', 'appinsights'),
    'docusaurus-node-polyfills'
  ],
  stylesheets: [
    "https://rsms.me/inter/inter.css",
    "https://fonts.googleapis.com/css2?family=Fira+Code&display=swap"
  ],
  themeConfig: {
    prism: {
      additionalLanguages: ['powershell', 'lua', 'jsstacktrace', 'yaml'],
    },
    navbar: {
      title: 'aliae 🌱',
      items: [
        {
          to: 'docs',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {
          href: 'https://github.com/sponsors/JanDeDobbeleer',
          label: 'Sponsor',
          position: 'left',
        },
        {
          href: 'https://github.com/jandedobbeleer/aliae',
          className: 'header-github-link',
          'aria-label': 'GitHub repository',
          position: 'right',
        },
        {
          href: 'https://www.gitkraken.com/invite/nQmDPR9D',
          className: 'header-gk-link',
          'aria-label': 'GitKraken',
          position: 'right',
        },
        {
          href: 'https://discord.gg/n7E3DkXssv',
          className: 'header-discord-link',
          'aria-label': 'Discord',
          position: 'right',
        },
        {
          href: 'https://staging.bsky.app/profile/aliae.dev',
          className: 'header-bluesky-link',
          'aria-label': 'Bluesky',
          position: 'right',
        }
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'How to',
          items: [
            {
              label: 'Getting started',
              to: 'docs/',
            },
            {
              label: 'Contributing',
              to: 'docs/contributing/started',
            },
          ],
        },
        {
          title: 'Social',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/jandedobbeleer/aliae',
            },
            {
              label: 'Discord',
              href: 'https://discord.gg/n7E3DkXssv',
            },
            {
              label: 'Bluesky',
              href: 'https://staging.bsky.app/profile/aliae.dev',
            }
          ],
        },
        {
          title: 'Links',
          items: [
            {
              label: 'Sponsor',
              href: 'https://github.com/sponsors/JanDeDobbeleer',
            },
            {
              label: 'GitKraken',
              href: 'https://www.gitkraken.com/invite/nQmDPR9D',
            },
            {
              label: 'Docusaurus',
              href: 'https://github.com/facebook/docusaurus',
            },
            {
              label: 'Privacy',
              href: '/privacy',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} <a href='https://github.com/sponsors/JanDeDobbeleer' target='_blank'>Jan De Dobbeleer</a> and <a href='/docs/contributors'>contributors</a>.`,
    },
    appInsights: {
      instrumentationKey: 'b204d2e6-9a10-473c-9655-2ad79af82e5c',
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/jandedobbeleer/aliae/edit/main/website/',
        },
        theme: {
          customCss: [
            require.resolve('./src/css/prism-rose-pine-moon.css'),
            require.resolve('./src/css/custom.css')
          ],
        },
      },
    ],
  ],
};
