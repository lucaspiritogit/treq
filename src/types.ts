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

export type DebugPair = {
  key: string;
  value: string;
};

export type DebugInfo = {
  startedAt: string;
  finishedAt: string;
  durationMs: number;
  request: {
    method: string;
    url: string;
    origin: string;
    pathname: string;
    search: string;
    params: DebugPair[];
    headers: DebugPair[];
    headerCount: number;
    bodyIncluded: boolean;
  };
  response: {
    status: number | null;
    statusText: string;
    ok: boolean;
    url: string;
    redirected: boolean;
    type: string;
    headers: DebugPair[];
    headerCount: number;
  } | null;
  errorMessage: string | null;
};
