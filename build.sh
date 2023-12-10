#!/bin/bash -xe
templ generate
./tailwindcss -c ./tailwind.config.js -i ./assets/styles/input.css -o ./assets/styles/output.css
