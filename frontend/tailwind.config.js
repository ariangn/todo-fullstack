export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#D69E2E",
          light: "#F6E05E",
          dark: "#B7791F",
        },
      },
      fontFamily: {
        header: ["joystix", "sans-serif"],
      },
    },
  },
  plugins: [],
};
