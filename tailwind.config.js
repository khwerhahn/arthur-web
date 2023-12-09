/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/**/*.{html,js,go}", "./components/*.{html,js,go}"],
  theme: {
    extend: {
      colors: {
        primary: "#7b16ff",
        secondary: "#565554",
        tertiary: "#16DB93",
      },
    },
  },
  plugins: [],
}

