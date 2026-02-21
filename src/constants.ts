import type { FocusField, HttpMethod } from "./types";

export const focusOrder: FocusField[] = ["requestList", "method", "url", "headers", "requestBody", "responseBody"];

export const methodColors: Record<HttpMethod, string> = {
  GET: "#22c55e",
  POST: "#0ea5e9",
  PUT: "#f59e0b",
  PATCH: "#f97316",
  DELETE: "#ef4444",
};

export type CommandSuggestion = {
  name: string;
  shortcuts: string[];
  description: string;
};

export const commandSuggestions: CommandSuggestion[] = [
  { name: "send", shortcuts: ["s", "run"], description: "Send current HTTP request" },
  { name: "save", shortcuts: [], description: "Save current request" },
  { name: "list", shortcuts: [], description: "Focus request list sidebar" },
  { name: "toggle-list", shortcuts: ["tl", "sidebar"], description: "Toggle request list sidebar" },
  { name: "reload", shortcuts: ["load"], description: "Reload saved requests file" },
  { name: "url", shortcuts: ["i", "input"], description: "Focus URL input" },
  { name: "headers", shortcuts: ["h"], description: "Focus headers input" },
  { name: "request", shortcuts: ["r", "req"], description: "Focus request body" },
  { name: "response", shortcuts: ["b", "res", "body"], description: "Focus response body" },
  { name: "get", shortcuts: ["g"], description: "Set method to GET" },
  { name: "post", shortcuts: ["p"], description: "Set method to POST" },
  { name: "put", shortcuts: ["u"], description: "Set method to PUT" },
  { name: "patch", shortcuts: ["t"], description: "Set method to PATCH" },
  { name: "delete", shortcuts: ["d"], description: "Set method to DELETE" },
  { name: "help", shortcuts: ["?"], description: "Open help modal" },
  { name: "debug", shortcuts: ["dbg"], description: "Open request/response debug modal" },
  { name: "quit", shortcuts: ["q", "exit"], description: "Exit application" },
];

export function resolveCommandAlias(rawCommand: string): string {
  const command = rawCommand.trim().toLowerCase();
  const match = commandSuggestions.find((item) => item.name === command || item.shortcuts.includes(command));
  return match?.name ?? command;
}
