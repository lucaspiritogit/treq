# treq 🦖
An HTTP client made with OpenTUI with a kinda vimified experience.

![image](./assets/demo.png)

## Run locally

```bash
bun install
bun run dev
```

## Features

- Vim-style command mode (`:`) with command suggestions (for example `:send`, `:save`, `:list`, `:help`)
- Fast keyboard-driven flow for method, URL, headers, request body, and response body
- Request list sidebar (toggleable), with keyboard navigation
- Save requests locally to `treq-requests.json`