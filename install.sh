#!/usr/bin/env sh
set -eu

REPO_URL="https://github.com/lucaspiritogit/treq.git"
BIN_NAME="treq"
INSTALL_DIR="/usr/local/lib/treq"
BIN_DIR="/usr/local/bin"

fail() {
  printf "%s\n" "$1" >&2
  exit 1
}

if ! command -v git >/dev/null 2>&1; then
  fail "git is required to install treq"
fi

if command -v bun >/dev/null 2>&1; then
  BUN_CMD="$(command -v bun)"
elif command -v curl >/dev/null 2>&1; then
  curl -fsSL https://bun.sh/install | sh
  if command -v bun >/dev/null 2>&1; then
    BUN_CMD="$(command -v bun)"
  elif [ -x "$HOME/.bun/bin/bun" ]; then
    BUN_CMD="$HOME/.bun/bin/bun"
  else
    fail "Bun was installed but is not available. Add \$HOME/.bun/bin to PATH and retry."
  fi
else
  fail "bun is required and curl is needed to install it automatically"
fi

if [ -d "$INSTALL_DIR/.git" ]; then
  git -C "$INSTALL_DIR" pull --ff-only
else
  INSTALL_PARENT="$(dirname "$INSTALL_DIR")"
  if [ -w "$INSTALL_PARENT" ]; then
    mkdir -p "$INSTALL_PARENT"
    git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"
  elif command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p "$INSTALL_PARENT"
    sudo git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"
    sudo chown -R "$USER" "$INSTALL_DIR"
  else
    fail "Cannot write to $INSTALL_PARENT. Run with sufficient permissions."
  fi
fi

"$BUN_CMD" install --cwd "$INSTALL_DIR"

LAUNCHER_CONTENT="#!/usr/bin/env sh
set -eu
exec \"$BUN_CMD\" run \"$INSTALL_DIR/src/index.tsx\" \"\$@\""

TARGET_PATH="${BIN_DIR}/${BIN_NAME}"
BIN_PARENT="$(dirname "$BIN_DIR")"
if [ -d "$BIN_DIR" ] && [ -w "$BIN_DIR" ]; then
  printf "%s\n" "$LAUNCHER_CONTENT" > "$TARGET_PATH"
  chmod +x "$TARGET_PATH"
elif [ ! -d "$BIN_DIR" ] && [ -w "$BIN_PARENT" ]; then
  mkdir -p "$BIN_DIR"
  printf "%s\n" "$LAUNCHER_CONTENT" > "$TARGET_PATH"
  chmod +x "$TARGET_PATH"
elif command -v sudo >/dev/null 2>&1; then
  sudo mkdir -p "$BIN_DIR"
  printf "%s\n" "$LAUNCHER_CONTENT" | sudo tee "$TARGET_PATH" >/dev/null
  sudo chmod +x "$TARGET_PATH"
else
  fail "No write access to ${BIN_DIR} and sudo is unavailable."
fi

printf "Installed %s to %s\n" "$BIN_NAME" "$TARGET_PATH"
