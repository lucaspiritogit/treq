#!/usr/bin/env sh

set -eu

REPO_URL="${TREQ_REPO_URL:-https://github.com/lpirito/treq.git}"
INSTALL_DIR="${TREQ_INSTALL_DIR:-$HOME/.local/share/treq}"
BIN_DIR="${TREQ_BIN_DIR:-$HOME/.local/bin}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf "Error: missing required command: %s\n" "$1" >&2
    exit 1
  fi
}

require_cmd git
require_cmd curl

BUN_CMD=""
if command -v bun >/dev/null 2>&1; then
  BUN_CMD="$(command -v bun)"
else
  printf "Bun not found. Installing Bun...\n"
  curl -fsSL https://bun.sh/install | sh

  if command -v bun >/dev/null 2>&1; then
    BUN_CMD="$(command -v bun)"
  elif [ -x "$HOME/.bun/bin/bun" ]; then
    BUN_CMD="$HOME/.bun/bin/bun"
  else
    printf "Error: Bun was installed but is not available in PATH.\n" >&2
    printf "Try running: export PATH=\"$HOME/.bun/bin:$PATH\"\n" >&2
    exit 1
  fi
fi

mkdir -p "$(dirname "$INSTALL_DIR")"
if [ -d "$INSTALL_DIR/.git" ]; then
  printf "Updating existing treq installation...\n"
  git -C "$INSTALL_DIR" pull --ff-only
else
  printf "Cloning treq into %s...\n" "$INSTALL_DIR"
  git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"
fi

printf "Installing dependencies...\n"
"$BUN_CMD" install --cwd "$INSTALL_DIR"

mkdir -p "$BIN_DIR"
LAUNCHER="$BIN_DIR/treq"

cat > "$LAUNCHER" <<EOF
#!/usr/bin/env sh
set -eu
exec "$BUN_CMD" run "$INSTALL_DIR/src/index.tsx" "\$@"
EOF

chmod +x "$LAUNCHER"

printf "\nDone. treq installed to %s\n" "$INSTALL_DIR"
printf "Launcher created at %s\n" "$LAUNCHER"
printf "Run it with: treq\n"

case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    printf "\nNote: %s is not in your PATH.\n" "$BIN_DIR"
    printf "Add this line to your shell profile:\n"
    printf "export PATH=\"%s:\$PATH\"\n" "$BIN_DIR"
    ;;
esac
