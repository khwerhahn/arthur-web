/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/**/*.{html,js,go,templ}", "./components/*.{html,js,go,templ}","./views/**/*.{html,js,go,templ}", "./views/*.{html,js,go,templ}", "./layouts/**/*.{html,js,go,templ}", "./layouts/*.{html,js,go,templ}"],
  theme: {
    extend: {
      colors: {
        black: "#1C2434",
        primary: "#7b16ff",
        primaryalt: "#6e13e5",
        secondary: "#565554",
        tertiary: "#16DB93",
      },
    },
  },
  plugins: [],
}

