/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx,md,mdx}", "../docs/**/*.{md,mdx}"],
  corePlugins: {
    preflight: false,
  },
  darkMode: ['class', '[data-theme="dark"]'],
  theme: {
    extend: {},
  },
  plugins: [],
};
