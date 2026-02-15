export type FocusField = "requestList" | "url" | "method" | "headers" | "requestBody" | "responseBody";

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

export type UiMode = "interactive" | "input";

export type MethodOption = {
  name: HttpMethod;
  description: string;
  value: HttpMethod;
};

export type SavedRequest = {
  id: string;
  name: string;
  method: HttpMethod;
  url: string;
  headers: string;
  body: string;
  createdAt: string;
};
