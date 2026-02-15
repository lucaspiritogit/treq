import type { DebugPair, FocusField } from "./types";

export type HighlightToken = {
  text: string;
  color?: string;
  bold?: boolean;
  italic?: boolean;
};

export function isTextEntryField(field: FocusField): boolean {
  return field === "url" || field === "headers" || field === "requestBody";
}

export function isCommandStarterKey(name: string, sequence: string, shift: boolean): boolean {
  return name === ":" || name === "colon" || sequence === ":" || (name === "semicolon" && shift);
}

export function getCommandChar(name: string, sequence: string): string | null {
  if (sequence.length === 1 && sequence >= " " && sequence <= "~") {
    return sequence;
  }
  if (name === "space") {
    return " ";
  }
  return null;
}

export function parseHeaders(rawHeaders: string): Record<string, string> {
  const result: Record<string, string> = {};
  const lines = rawHeaders.split("\n");
  for (const line of lines) {
    const trimmedLine = line.trim();
    if (!trimmedLine) {
      continue;
    }
    const separatorIndex = trimmedLine.indexOf(":");
    if (separatorIndex === -1) {
      continue;
    }
    const key = trimmedLine.slice(0, separatorIndex).trim();
    const value = trimmedLine.slice(separatorIndex + 1).trim();
    if (key) {
      result[key] = value;
    }
  }
  return result;
}

export function formatResponseBody(text: string): string {
  try {
    const parsed = JSON.parse(text);
    return JSON.stringify(parsed, null, 2);
  } catch {
    return text;
  }
}

export function isSensitiveHeaderName(name: string): boolean {
  const normalized = name.trim().toLowerCase();
  return normalized === "authorization" || normalized === "x-api-key" || normalized === "cookie" || normalized === "set-cookie";
}

export function maskSensitiveHeaderValue(name: string, value: string): string {
  if (!isSensitiveHeaderName(name)) {
    return value;
  }

  const trimmedValue = value.trim();
  if (!trimmedValue) {
    return "********";
  }

  if (name.trim().toLowerCase() === "authorization") {
    const separatorIndex = trimmedValue.indexOf(" ");
    if (separatorIndex > 0) {
      const scheme = trimmedValue.slice(0, separatorIndex);
      return `${scheme} ********`;
    }
  }

  return "********";
}

export function maskSensitiveHeadersText(rawHeaders: string): string {
  return rawHeaders
    .split("\n")
    .map((line) => {
      const separatorIndex = line.indexOf(":");
      if (separatorIndex === -1) {
        return line;
      }
      const key = line.slice(0, separatorIndex).trim();
      const value = line.slice(separatorIndex + 1).trim();
      const maskedValue = maskSensitiveHeaderValue(key, value);
      return `${key}: ${maskedValue}`;
    })
    .join("\n");
}

export function maskSensitiveHeaderPairs(headerPairs: DebugPair[]): DebugPair[] {
  return headerPairs.map((headerPair) => ({
    key: headerPair.key,
    value: maskSensitiveHeaderValue(headerPair.key, headerPair.value),
  }));
}

export function highlightJsonLine(line: string): HighlightToken[] {
  const tokens: HighlightToken[] = [];
  const pattern = /"(?:\\.|[^"\\])*"\s*:|"(?:\\.|[^"\\])*"|\b-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?\b|\btrue\b|\bfalse\b|\bnull\b|[{}\[\],:]/g;
  let lastIndex = 0;

  for (const match of line.matchAll(pattern)) {
    const value = match[0];
    const index = match.index ?? 0;

    if (index > lastIndex) {
      tokens.push({ text: line.slice(lastIndex, index), color: "#e2e8f0" });
    }

    if (value.endsWith(":")) {
      tokens.push({ text: value.slice(0, -1), color: "#93c5fd" });
      tokens.push({ text: ":", color: "#94a3b8" });
    } else if (value.startsWith("\"")) {
      tokens.push({ text: value, color: "#34d399" });
    } else if (value === "true" || value === "false") {
      tokens.push({ text: value, color: "#60a5fa", bold: true });
    } else if (value === "null") {
      tokens.push({ text: value, color: "#a78bfa", italic: true });
    } else if (value === "{" || value === "}" || value === "[" || value === "]" || value === "," || value === ":") {
      tokens.push({ text: value, color: "#94a3b8" });
    } else {
      tokens.push({ text: value, color: "#f59e0b" });
    }

    lastIndex = index + value.length;
  }

  if (lastIndex < line.length) {
    tokens.push({ text: line.slice(lastIndex), color: "#e2e8f0" });
  }

  if (tokens.length === 0) {
    return [{ text: line || " ", color: "#e2e8f0" }];
  }

  return tokens;
}
