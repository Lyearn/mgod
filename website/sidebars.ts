import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'category',
      label: 'Introduction',
      items: [
        'about',
      ],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Beginner\'s Guide',
      items: [
        'installation',
        'basic_usage',
      ],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Features',
      items: [
        'schema_options',
        'field_options',
        'field_transformers',
        'meta_fields',
      ],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Advanced Guide',
      items: [
        'union_types',
      ],
      collapsed: false,
    },
  ],
};

export default sidebars;
