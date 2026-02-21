import type { FocusField, HttpMethod } from "./types";

export const focusOrder: FocusField[] = ["requestList", "method", "url", "headers", "requestBody", "responseBody"];

export const methodColors: Record<HttpMethod, string> = {
  GET: "#22c55e",
  POST: "#0ea5e9",
  PUT: "#f59e0b",
  PATCH: "#f97316",
  DELETE: "#ef4444",
};

export const commandSuggestions = [
  { name: "send", description: "Send current HTTP request" },
  { name: "save", description: "Save current request" },
  { name: "list", description: "Focus request list sidebar" },
  { name: "toggle-list", description: "Toggle request list sidebar" },
  { name: "reload", description: "Reload saved requests file" },
  { name: "url", description: "Focus URL input" },
  { name: "headers", description: "Focus headers input" },
  { name: "request", description: "Focus request body" },
  { name: "response", description: "Focus response body" },
  { name: "get", description: "Set method to GET" },
  { name: "post", description: "Set method to POST" },
  { name: "put", description: "Set method to PUT" },
  { name: "patch", description: "Set method to PATCH" },
  { name: "delete", description: "Set method to DELETE" },
  { name: "help", description: "Open help modal" },
  { name: "debug", description: "Open request/response debug modal" },
  { name: "quit", description: "Exit application" },
] as const;
