css:
	npx tailwindcss -i ./views/css/app.css -o ./public/css/styles.css --watch

templ:
	templ generate --watch --proxy="http://localhost:4000"

js:
	node esbuild.config.js
