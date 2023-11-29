/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#F7F6F3',
          100: '#F1F0ED',
          200: '#E3E2E0',
          300: '#D1D0CE',
          400: '#AEADAB',
          500: '#9B9A97',
          600: '#6B6966',
          700: '#4F4D4A',
          800: '#37352F',
          900: '#2F2D29',
        },
        notion: {
          bg: '#F7F6F3',
          text: '#37352F',
          textLight: '#6B6966',
          border: '#E3E2E0',
          hover: '#F1F0ED',
        },
      },
    },
  },
  plugins: [],
}

