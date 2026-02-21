# treq 🦖
`treq` is a keyboard-first HTTP client for the terminal, built with OpenTUI. It is designed around a vim-like workflow so you can move between URL, headers, request body, response, saved requests, and commands without touching the mouse.
![image](./assets/demo.png)

## Run locally

```bash
bun install
bun run dev
```

## Install

```bash
curl -fsSSl https://raw.githubusercontent.com/lucaspiritogit/treq/main/install.sh | sh
```

## Features

- Vim-style command mode (`:`) with command suggestions (for example `:send`, `:save`, `:list`, `:help`)
- Fast keyboard-driven flow for method, URL, headers, request body, and response body
- Request list sidebar (toggleable), with keyboard navigation
- Save requests locally to `treq-requests.json`

## Keyboard Shortcuts

- `:` open command mode
- `Esc` switch to interactive mode and interactive method
- `Tab` / `Shift+Tab` cycle focus between panels
- `Ctrl+Enter` / `Cmd+Enter` / `Alt+Enter` send request
- `Enter` send request (interactive mode, non-request-list focus)
- `i` focus URL input
- `h` focus headers input
- `r` focus request body input
- `b` focus response body panel
- `l` / `Left` focus request list (opens list if hidden)
- `g` set method `GET`
- `p` set method `POST`
- `u` set method `PUT`
- `t` set method `PATCH`
- `d` set method `DELETE`
- Request list: `Up/Down` or `k/j` navigate, `Enter` load selected request, `Ctrl+d` / `Cmd+d` delete selected request

## Commands

- `:send`, `:s`, `:run` send current request
- `:save` save current request (overwrites loaded request)
- `:list`, `:sidebar` toggle request list sidebar
- `:reload`, `:load` reload requests from `treq-requests.json`
- `:url`, `:i`, `:input` focus URL input
- `:headers`, `:h` focus headers input
- `:request`, `:req`, `:r` focus request body input
- `:response`, `:res`, `:body`, `:b` focus response body
- `:get`, `:g` set method `GET`
- `:post`, `:p` set method `POST`
- `:put`, `:u` set method `PUT`
- `:patch`, `:t` set method `PATCH`
- `:delete`, `:d` set method `DELETE`
- `:debug`, `:dbg` open debug modal to see additional context about the request and the response.
- `:help` open command help modal
- `:quit`, `:q`, `:exit` close app
