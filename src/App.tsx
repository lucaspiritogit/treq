import { type InputRenderable, type TextareaRenderable } from "@opentui/core";
import { useKeyboard, useRenderer, useTerminalDimensions } from "@opentui/react";
import { useCallback, useEffect, useRef, useState } from "react";
import { BodyPanels } from "./components/BodyPanels";
import { CommandPanel } from "./components/CommandPanel";
import { CommandSuggestionsPanel } from "./components/CommandSuggestionsPanel";
import { DeleteRequestModal } from "./components/DeleteRequestModal";
import { DebugModal } from "./components/DebugModal";
import { HelpModal } from "./components/HelpModal";
import { MethodPanel } from "./components/MethodPanel";
import { ModePanel } from "./components/ModePanel";
import { RequestListPanel } from "./components/RequestListPanel";
import { SaveRequestModal } from "./components/SaveRequestModal";
import { focusOrder, methodColors } from "./constants";
import type { DebugInfo, FocusField, HttpMethod, SavedRequest, UiMode } from "./types";
import {
  formatResponseBody,
  getCommandChar,
  isCommandStarterKey,
  isTextEntryField,
  maskSensitiveHeaderPairs,
  maskSensitiveHeadersText,
  parseHeaders,
} from "./utils";

const REQUESTS_FILE_PATH = `${process.cwd()}/treq-requests.json`;

export function App() {
  const renderer = useRenderer();
  const terminalDimensions = useTerminalDimensions();
  const [url, setUrl] = useState("");
  const [method, setMethod] = useState<HttpMethod>("GET");
  const [requestBodyText, setRequestBodyText] = useState("");
  const [headersText, setHeadersText] = useState("");
  const [responseBody, setResponseBody] = useState("");
  const [responseStatus, setResponseStatus] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [uiMode, setUiMode] = useState<UiMode>("input");
  const [focusField, setFocusField] = useState<FocusField>("url");
  const [commandMode, setCommandMode] = useState(false);
  const [commandLine, setCommandLine] = useState("");
  const [commandFeedback, setCommandFeedback] = useState("");
  const [helpModalOpen, setHelpModalOpen] = useState(false);
  const [savedRequests, setSavedRequests] = useState<SavedRequest[]>([]);
  const [activeRequestId, setActiveRequestId] = useState<string | null>(null);
  const [requestListOpen, setRequestListOpen] = useState(true);
  const [requestListCursorIndex, setRequestListCursorIndex] = useState(0);
  const [headersEditorVersion, setHeadersEditorVersion] = useState(0);
  const [requestBodyEditorVersion, setRequestBodyEditorVersion] = useState(0);
  const [saveModalOpen, setSaveModalOpen] = useState(false);
  const [saveRequestName, setSaveRequestName] = useState("");
  const [debugModalOpen, setDebugModalOpen] = useState(false);
  const [debugInfo, setDebugInfo] = useState<DebugInfo | null>(null);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [requestIdPendingDelete, setRequestIdPendingDelete] = useState<string | null>(null);
  const urlRef = useRef<InputRenderable>(null);
  const headersRef = useRef<TextareaRenderable>(null);
  const requestBodyRef = useRef<TextareaRenderable>(null);

  const loadSavedRequests = useCallback(async () => {
    try {
      const file = Bun.file(REQUESTS_FILE_PATH);
      if (!(await file.exists())) {
        setSavedRequests([]);
        return;
      }

      const content = await file.text();
      if (!content.trim()) {
        setSavedRequests([]);
        return;
      }

      const parsed = JSON.parse(content);
      if (!Array.isArray(parsed)) {
        setCommandFeedback("Invalid treq-requests.json format");
        return;
      }

      const normalized = parsed.filter((item): item is SavedRequest => {
        return (
          typeof item === "object" &&
          item !== null &&
          typeof item.id === "string" &&
          typeof item.name === "string" &&
          typeof item.method === "string" &&
          typeof item.url === "string" &&
          typeof item.headers === "string" &&
          typeof item.body === "string" &&
          typeof item.createdAt === "string"
        );
      });

      setSavedRequests(normalized);
    } catch {
      setCommandFeedback("Failed to load treq-requests.json");
    }
  }, []);

  useEffect(() => {
    void loadSavedRequests();
  }, [loadSavedRequests]);

  const saveCurrentRequest = useCallback(async (name: string) => {
    const trimmedUrl = (urlRef.current?.value ?? url).trim();
    if (!trimmedUrl) {
      setCommandFeedback("Cannot save: URL is empty");
      return;
    }

    const trimmedName = name.trim();
    if (!trimmedName) {
      setCommandFeedback("Cannot save: name is empty");
      return;
    }

    const isEditingHeaders = uiMode === "input" && focusField === "headers";
    const isEditingRequestBody = uiMode === "input" && focusField === "requestBody";
    const headersValue = isEditingHeaders ? (headersRef.current?.plainText ?? headersText) : headersText;
    const bodyValue = isEditingRequestBody ? (requestBodyRef.current?.plainText ?? requestBodyText) : requestBodyText;
    const existingRequestIndex = activeRequestId
      ? savedRequests.findIndex((savedRequest) => savedRequest.id === activeRequestId)
      : -1;
    let nextRequests: SavedRequest[] = [];
    let savedRequestId = activeRequestId;

    if (existingRequestIndex >= 0) {
      const existingRequest = savedRequests[existingRequestIndex];
      if (!existingRequest) {
        return;
      }

      const updatedRequest: SavedRequest = {
        id: existingRequest.id,
        name: trimmedName,
        method,
        url: trimmedUrl,
        headers: maskSensitiveHeadersText(headersValue),
        body: bodyValue,
        createdAt: existingRequest.createdAt,
      };
      nextRequests = savedRequests.map((savedRequest, index) => {
        if (index === existingRequestIndex) {
          return updatedRequest;
        }
        return savedRequest;
      });
      savedRequestId = existingRequest.id;
    } else {
      const nextRequest: SavedRequest = {
        id: `${Date.now()}`,
        name: trimmedName,
        method,
        url: trimmedUrl,
        headers: maskSensitiveHeadersText(headersValue),
        body: bodyValue,
        createdAt: new Date().toISOString(),
      };
      nextRequests = [...savedRequests, nextRequest];
      savedRequestId = nextRequest.id;
    }

    try {
      await Bun.write(REQUESTS_FILE_PATH, `${JSON.stringify(nextRequests, null, 2)}\n`);
      setSavedRequests(nextRequests);
      if (savedRequestId) {
        setActiveRequestId(savedRequestId);
      }
      setCommandFeedback(existingRequestIndex >= 0 ? `Updated request (${trimmedName})` : `Saved request (${trimmedName})`);
      setSaveModalOpen(false);
      setSaveRequestName("");
    } catch {
      setCommandFeedback("Failed to write treq-requests.json");
    }
  }, [activeRequestId, focusField, headersText, method, requestBodyText, savedRequests, uiMode, url]);

  const openSaveModal = useCallback(() => {
    const trimmedUrl = (urlRef.current?.value ?? url).trim();
    if (!trimmedUrl) {
      setCommandFeedback("Cannot save: URL is empty");
      return;
    }

    const activeRequest = activeRequestId
      ? savedRequests.find((savedRequest) => savedRequest.id === activeRequestId)
      : null;

    if (activeRequest) {
      void saveCurrentRequest(activeRequest.name);
      return;
    }

    setSaveRequestName(`${method} ${trimmedUrl}`);
    setSaveModalOpen(true);
  }, [activeRequestId, method, saveCurrentRequest, savedRequests, url]);

  const loadRequestIntoForm = useCallback((request: SavedRequest) => {
    setMethod(request.method);
    setUrl(request.url);
    setHeadersText(request.headers);
    setRequestBodyText(request.body);
    setActiveRequestId(request.id);
    setHeadersEditorVersion((value) => value + 1);
    setRequestBodyEditorVersion((value) => value + 1);
    setUiMode("input");
    setFocusField("url");
    setCommandFeedback(`Loaded: ${request.name}`);
  }, []);

  const openDeleteRequestModal = useCallback(() => {
    const requestToDelete = savedRequests[requestListCursorIndex];
    if (!requestToDelete) {
      return;
    }

    setRequestIdPendingDelete(requestToDelete.id);
    setDeleteModalOpen(true);
  }, [requestListCursorIndex, savedRequests]);

  const confirmDeleteRequest = useCallback(async () => {
    if (!requestIdPendingDelete) {
      setDeleteModalOpen(false);
      return;
    }

    const nextRequests = savedRequests.filter((savedRequest) => savedRequest.id !== requestIdPendingDelete);
    try {
      await Bun.write(REQUESTS_FILE_PATH, `${JSON.stringify(nextRequests, null, 2)}\n`);
      await loadSavedRequests();
      if (activeRequestId === requestIdPendingDelete) {
        setActiveRequestId(null);
      }
      setCommandFeedback(`Deleted request (${nextRequests.length})`);
    } catch {
      setCommandFeedback("Failed to delete request");
    } finally {
      setDeleteModalOpen(false);
      setRequestIdPendingDelete(null);
    }
  }, [activeRequestId, loadSavedRequests, requestIdPendingDelete, savedRequests]);

  useEffect(() => {
    if (!activeRequestId) {
      return;
    }

    const hasActiveRequest = savedRequests.some((savedRequest) => savedRequest.id === activeRequestId);
    if (!hasActiveRequest) {
      setActiveRequestId(null);
    }
  }, [activeRequestId, savedRequests]);

  useEffect(() => {
    if (!requestListOpen && focusField === "requestList") {
      setFocusField("method");
    }
  }, [focusField, requestListOpen]);

  useEffect(() => {
    if (savedRequests.length === 0) {
      setRequestListCursorIndex(0);
      return;
    }

    if (requestListCursorIndex >= savedRequests.length) {
      setRequestListCursorIndex(savedRequests.length - 1);
    }
  }, [requestListCursorIndex, savedRequests.length]);

  const sendRequest = useCallback(async () => {
    setCommandMode(false);
    setCommandLine("");
    setCommandFeedback("");

    const trimmedUrl = (urlRef.current?.value ?? url).trim();
    const isEditingHeaders = uiMode === "input" && focusField === "headers";
    const isEditingRequestBody = uiMode === "input" && focusField === "requestBody";
    const headersValue = isEditingHeaders ? (headersRef.current?.plainText ?? headersText) : headersText;
    const requestBodyValue = isEditingRequestBody ? (requestBodyRef.current?.plainText ?? requestBodyText) : requestBodyText;

    if (!trimmedUrl) {
      setResponseStatus("URL is required");
      setResponseBody("");
      return;
    }

    const normalizedMethod = method.trim().toUpperCase() || "GET";
    const parsedRequestHeaders = parseHeaders(headersValue);
    const requestInit: RequestInit = {
      method: normalizedMethod,
      headers: parsedRequestHeaders,
    };

    const bodyIncluded = normalizedMethod !== "GET" && normalizedMethod !== "HEAD" && !!requestBodyValue.trim();
    if (bodyIncluded) {
      requestInit.body = requestBodyValue;
    }

    let parsedUrl: URL | null = null;
    try {
      parsedUrl = new URL(trimmedUrl);
    } catch {
      parsedUrl = null;
    }

    const requestParams = parsedUrl
      ? Array.from(parsedUrl.searchParams.entries()).map((entry) => ({ key: entry[0], value: entry[1] }))
      : [];
    const requestHeaderPairs = Object.entries(parsedRequestHeaders).map((entry) => ({ key: entry[0], value: entry[1] }));
    const maskedRequestHeaderPairs = maskSensitiveHeaderPairs(requestHeaderPairs);

    setIsLoading(true);
    setResponseStatus("Sending request...");
    const startedAtMs = Date.now();
    const startedAtIso = new Date(startedAtMs).toISOString();
    try {
      const response = await fetch(trimmedUrl, requestInit);
      const responseText = await response.text();
      const finishedAtMs = Date.now();
      const responseHeaders = Array.from(response.headers.entries()).map((entry) => ({ key: entry[0], value: entry[1] }));
      setResponseStatus(`${response.status} ${response.statusText}`);
      setResponseBody(formatResponseBody(responseText));
      setDebugInfo({
        startedAt: startedAtIso,
        finishedAt: new Date(finishedAtMs).toISOString(),
        durationMs: finishedAtMs - startedAtMs,
        request: {
          method: normalizedMethod,
          url: trimmedUrl,
          origin: parsedUrl?.origin ?? "",
          pathname: parsedUrl?.pathname ?? "",
          search: parsedUrl?.search ?? "",
          params: requestParams,
          headers: maskedRequestHeaderPairs,
          headerCount: requestHeaderPairs.length,
          bodyIncluded,
        },
        response: {
          status: response.status,
          statusText: response.statusText,
          ok: response.ok,
          url: response.url,
          redirected: response.redirected,
          type: response.type,
          headers: responseHeaders,
          headerCount: responseHeaders.length,
        },
        errorMessage: null,
      });
    } catch (error) {
      const finishedAtMs = Date.now();
      const message = error instanceof Error ? error.message : "Unknown request error";
      setResponseStatus("Request failed");
      setResponseBody(message);
      setDebugInfo({
        startedAt: startedAtIso,
        finishedAt: new Date(finishedAtMs).toISOString(),
        durationMs: finishedAtMs - startedAtMs,
        request: {
          method: normalizedMethod,
          url: trimmedUrl,
          origin: parsedUrl?.origin ?? "",
          pathname: parsedUrl?.pathname ?? "",
          search: parsedUrl?.search ?? "",
          params: requestParams,
          headers: maskedRequestHeaderPairs,
          headerCount: requestHeaderPairs.length,
          bodyIncluded,
        },
        response: null,
        errorMessage: message,
      });
    } finally {
      setIsLoading(false);
    }
  }, [focusField, headersText, method, requestBodyText, uiMode, url]);

  const runCommand = useCallback((rawCommand: string) => {
    setCommandFeedback("");
    const command = rawCommand.trim().toLowerCase();

    if (!command) {
      return;
    }

    if (command === "q" || command === "quit" || command === "exit") {
      renderer.destroy();
      return;
    }

    if (command === "s" || command === "send" || command === "run") {
      void sendRequest();
      return;
    }

    if (command === "save") {
      openSaveModal();
      return;
    }

    if (command === "list" || command === "sidebar") {
      setRequestListOpen((value) => !value);
      return;
    }

    if (command === "reload" || command === "load") {
      void loadSavedRequests();
      return;
    }

    if (command === "i" || command === "input" || command === "url") {
      setUiMode("input");
      setFocusField("url");
      return;
    }

    if (command === "h" || command === "headers") {
      setUiMode("input");
      setFocusField("headers");
      return;
    }

    if (command === "r" || command === "req" || command === "request") {
      setUiMode("input");
      setFocusField("requestBody");
      return;
    }

    if (command === "b" || command === "res" || command === "response" || command === "body") {
      setUiMode("interactive");
      setFocusField("responseBody");
      return;
    }

    if (command === "g" || command === "get") {
      setMethod("GET");
      return;
    }

    if (command === "p" || command === "post") {
      setMethod("POST");
      return;
    }

    if (command === "u" || command === "put") {
      setMethod("PUT");
      return;
    }

    if (command === "d" || command === "delete") {
      setMethod("DELETE");
      return;
    }

    if (command === "help") {
      setHelpModalOpen(true);
      return;
    }

    if (command === "debug" || command === "dbg") {
      setDebugModalOpen(true);
      return;
    }

    setCommandFeedback("");
  }, [loadSavedRequests, openSaveModal, renderer, sendRequest]);

  useKeyboard((key) => {
    if (deleteModalOpen) {
      if (key.name === "escape" || key.name === "n" || key.name === "q") {
        setDeleteModalOpen(false);
        setRequestIdPendingDelete(null);
        return;
      }

      if (key.name === "enter" || key.name === "return" || key.name === "y") {
        void confirmDeleteRequest();
        return;
      }

      return;
    }

    if (debugModalOpen) {
      if (key.name === "escape" || key.name === "q" || key.name === "enter" || key.name === "return") {
        setDebugModalOpen(false);
      }
      return;
    }

    if (helpModalOpen) {
      if (key.name === "escape" || key.name === "q" || key.name === "enter" || key.name === "return") {
        setHelpModalOpen(false);
      }
      return;
    }

    if (saveModalOpen) {
      if (key.name === "escape") {
        setSaveModalOpen(false);
        setSaveRequestName("");
      } else if (key.name === "enter" || key.name === "return") {
        void saveCurrentRequest(saveRequestName);
      }
      return;
    }

    if (!commandMode && (key.ctrl || key.meta || key.option) && (key.name === "enter" || key.name === "return")) {
      void sendRequest();
      return;
    }

    if (!commandMode && uiMode === "interactive" && focusField !== "requestList" && (key.name === "enter" || key.name === "return")) {
      void sendRequest();
      return;
    }

    if (key.ctrl && key.name === "c") {
      renderer.destroy();
      return;
    }

    if (uiMode === "input" && isTextEntryField(focusField) && commandMode) {
      setCommandMode(false);
      setCommandLine("");
      setCommandFeedback("");
    }

    if (commandMode) {
      if (key.name === "escape") {
        setCommandMode(false);
        setCommandLine("");
        return;
      }

      if (key.name === "backspace") {
        setCommandLine((value) => value.slice(0, Math.max(0, value.length - 1)));
        return;
      }

      if (key.name === "enter" || key.name === "return") {
        const value = commandLine;
        setCommandMode(false);
        setCommandFeedback("");
        setCommandLine("");

        runCommand(value);
        return;
      }

      const commandChar = getCommandChar(key.name, key.sequence);
      if (commandChar) {
        setCommandLine((value) => value + commandChar);
      }
      return;
    }

    if (key.name === "escape") {
      setCommandMode(false);
      setCommandLine("");
      setCommandFeedback("");
      setUiMode("interactive");
      setFocusField("method");
      return;
    }

    if (key.name === "tab") {
      const activeFocusOrder = requestListOpen
        ? focusOrder
        : focusOrder.filter((field) => field !== "requestList");
      const currentIndex = activeFocusOrder.indexOf(focusField);
      const direction = key.shift ? -1 : 1;
      const nextIndex = (currentIndex + direction + activeFocusOrder.length) % activeFocusOrder.length;
      const nextFocus = activeFocusOrder[nextIndex];
      if (nextFocus) {
        setFocusField(nextFocus);
        setUiMode(isTextEntryField(nextFocus) ? "input" : "interactive");
      }
      return;
    }

    if (uiMode === "interactive") {
      if (isCommandStarterKey(key.name, key.sequence, key.shift)) {
        setCommandMode(true);
        setCommandLine("");
        setCommandFeedback("");
        return;
      }

      if (key.name === "i") {
        setUiMode("input");
        setFocusField("url");
        return;
      }

      if (key.name === "h") {
        setUiMode("input");
        setFocusField("headers");
        return;
      }

      if (key.name === "g") {
        setMethod("GET");
        return;
      }

      if (key.name === "p") {
        setMethod("POST");
        return;
      }

      if (key.name === "u") {
        setMethod("PUT");
        return;
      }

      if (key.ctrl && key.name === "d" && focusField === "requestList") {
        openDeleteRequestModal();
        return;
      }

      if (key.name === "d") {
        setMethod("DELETE");
        return;
      }

      if (key.name === "l" || key.name === "left") {
        if (!requestListOpen) {
          setRequestListOpen(true);
        }
        setFocusField("requestList");
        return;
      }

      if (key.name === "r") {
        setUiMode("input");
        setFocusField("requestBody");
        return;
      }

      if (key.name === "b") {
        setFocusField("responseBody");
        return;
      }
    }

    if (uiMode === "interactive" && focusField === "requestList") {
      if (savedRequests.length === 0) {
        if (key.name === "right") {
          setFocusField("method");
        }
        return;
      }

      if (key.name === "up" || key.name === "k") {
        setRequestListCursorIndex((value) => {
          const nextValue = value - 1;
          return nextValue < 0 ? savedRequests.length - 1 : nextValue;
        });
        return;
      }

      if (key.name === "down" || key.name === "j") {
        setRequestListCursorIndex((value) => {
          const nextValue = value + 1;
          return nextValue >= savedRequests.length ? 0 : nextValue;
        });
        return;
      }

      if (key.name === "enter" || key.name === "return") {
        const selectedRequest = savedRequests[requestListCursorIndex];
        if (selectedRequest) {
          loadRequestIntoForm(selectedRequest);
        }
        return;
      }

      if (key.name === "right") {
        setFocusField("method");
        return;
      }

    }

    if (uiMode === "input" && isTextEntryField(focusField)) {
      return;
    }

    if (focusField === "method" && (key.name === "enter" || key.name === "return")) {
      void sendRequest();
    }
  });

  const stackBodyPanels = terminalDimensions.width < (requestListOpen ? 90 : 30);
  const maskedHeadersText = maskSensitiveHeadersText(headersRef.current?.plainText ?? headersText);
  const maskedHeaderLines = maskedHeadersText.split("\n");
  const headersPanelContentHeight = 6;
  const headersPanelHeight = headersPanelContentHeight + 2;

  return (
    <box flexDirection="column" width="100%" height="100%" padding={1} gap={1}>
      <box flexDirection="row" flexGrow={1} minHeight={0} minWidth={0} gap={1}>
        {requestListOpen ? (
          <RequestListPanel
            requests={savedRequests}
            activeRequestId={activeRequestId}
            focused={uiMode === "interactive" && focusField === "requestList"}
            cursorIndex={requestListCursorIndex}
          />
        ) : null}

        <box flexDirection="column" flexGrow={1} minHeight={0} minWidth={0} gap={1}>
          <box flexDirection="row" minWidth={0} gap={1}>
            <MethodPanel
              focusField={focusField}
              uiMode={uiMode}
              method={method}
              methodColors={methodColors}
            />
            <box border padding={1} title="URL" flexGrow={1} flexShrink={1} minWidth={0}>
              <input
                ref={urlRef}
                value={url}
                onChange={setUrl}
                onInput={setUrl}
                onSubmit={() => {
                  void sendRequest();
                }}
                focused={uiMode === "input" && focusField === "url"}
              />
            </box>
          </box>

          <box border padding={1} title="Headers (Key: Value)" height={headersPanelHeight} minHeight={headersPanelHeight} flexShrink={0}>
            {uiMode === "input" && focusField === "headers" ? (
              <textarea
                key={`headers-${headersEditorVersion}`}
                ref={headersRef}
                initialValue={headersText}
                onContentChange={() => {
                  setHeadersText(headersRef.current?.plainText ?? "");
                }}
                onSubmit={() => {
                  void sendRequest();
                }}
                keyBindings={[
                  { name: "return", ctrl: true, action: "submit" },
                  { name: "enter", ctrl: true, action: "submit" },
                ]}
                focused
                wrapMode="word"
                width="100%"
                minWidth={0}
                flexGrow={1}
                minHeight={0}
                height={headersPanelContentHeight}
              />
            ) : (
              <scrollbox focused={focusField === "headers"} flexGrow={1} minHeight={0} height={headersPanelContentHeight}>
                {maskedHeaderLines.map((line, lineIndex) => (
                  <text key={`header-preview-${lineIndex}`}>
                    <span fg="#e2e8f0">{line || " "}</span>
                  </text>
                ))}
              </scrollbox>
            )}
          </box>

          <BodyPanels
            key={`body-${requestBodyEditorVersion}`}
            uiMode={uiMode}
            focusField={focusField}
            requestBodyText={requestBodyText}
            setRequestBodyText={setRequestBodyText}
            requestBodyRef={requestBodyRef}
            sendRequest={sendRequest}
            isLoading={isLoading}
            responseBody={responseBody}
            responseStatus={responseStatus}
            stacked={stackBodyPanels}
          />

          <CommandSuggestionsPanel commandMode={commandMode} commandLine={commandLine} />
          <CommandPanel commandMode={commandMode} commandLine={commandLine} commandFeedback={commandFeedback} />
          <ModePanel uiMode={uiMode} />
        </box>
      </box>
      <HelpModal isOpen={helpModalOpen} />
      <DebugModal isOpen={debugModalOpen} debugInfo={debugInfo} />
      <DeleteRequestModal
        isOpen={deleteModalOpen}
        requestToDelete={savedRequests.find((savedRequest) => savedRequest.id === requestIdPendingDelete) ?? null}
        onConfirm={() => {
          void confirmDeleteRequest();
        }}
        onCancel={() => {
          setDeleteModalOpen(false);
          setRequestIdPendingDelete(null);
        }}
      />
      <SaveRequestModal
        isOpen={saveModalOpen}
        requestName={saveRequestName}
        setRequestName={setSaveRequestName}
        onConfirm={() => {
          void saveCurrentRequest(saveRequestName);
        }}
      />
    </box>
  );
}
