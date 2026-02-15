import type { DebugInfo } from "../types";

type DebugModalProps = {
  isOpen: boolean;
  debugInfo: DebugInfo | null;
};

export function DebugModal(props: DebugModalProps) {
  if (!props.isOpen) {
    return null;
  }

  return (
    <box
      position="absolute"
      top={0}
      right={0}
      bottom={0}
      left={0}
      zIndex={40}
      backgroundColor="#020617"
    >
      <box border width="100%" height="100%" padding={1} title="Debug Request/Response" backgroundColor="#020617">
        <scrollbox focused flexGrow={1}>
          {props.debugInfo ? (
            <>
              <text><strong>Request</strong></text>
              <text><span fg="#94a3b8">Method:</span> <span fg="#f8fafc">{props.debugInfo.request.method}</span></text>
              <text><span fg="#94a3b8">URL:</span> <span fg="#f8fafc">{props.debugInfo.request.url}</span></text>
              <text><span fg="#94a3b8">Origin:</span> <span fg="#f8fafc">{props.debugInfo.request.origin || "(unknown)"}</span></text>
              <text><span fg="#94a3b8">Path:</span> <span fg="#f8fafc">{props.debugInfo.request.pathname || "/"}</span></text>
              <text><span fg="#94a3b8">Query:</span> <span fg="#f8fafc">{props.debugInfo.request.search || "(none)"}</span></text>
              <text><span fg="#94a3b8">Body sent:</span> <span fg="#f8fafc">{props.debugInfo.request.bodyIncluded ? "yes" : "no"}</span></text>
              <text> </text>
              <text><strong>Params</strong></text>
              {props.debugInfo.request.params.length > 0 ? (
                props.debugInfo.request.params.map((param, index) => (
                  <text key={`param-${index}`}>
                    <span fg="#93c5fd">{param.key}</span>
                    <span fg="#94a3b8"> = </span>
                    <span fg="#f8fafc">{param.value}</span>
                  </text>
                ))
              ) : (
                <text><span fg="#64748b">(none)</span></text>
              )}
              <text> </text>
              <text><strong>Sent Headers ({props.debugInfo.request.headerCount})</strong></text>
              {props.debugInfo.request.headers.length > 0 ? (
                props.debugInfo.request.headers.map((header, index) => (
                  <text key={`request-header-${index}`}>
                    <span fg="#93c5fd">{header.key}</span>
                    <span fg="#94a3b8">: </span>
                    <span fg="#f8fafc">{header.value}</span>
                  </text>
                ))
              ) : (
                <text><span fg="#64748b">(none)</span></text>
              )}
              <text> </text>
              <text><strong>Response</strong></text>
              {props.debugInfo.response ? (
                <>
                  <text><span fg="#94a3b8">Status:</span> <span fg="#f8fafc">{props.debugInfo.response.status} {props.debugInfo.response.statusText}</span></text>
                  <text><span fg="#94a3b8">OK:</span> <span fg="#f8fafc">{props.debugInfo.response.ok ? "true" : "false"}</span></text>
                  <text><span fg="#94a3b8">Final URL:</span> <span fg="#f8fafc">{props.debugInfo.response.url}</span></text>
                  <text><span fg="#94a3b8">Redirected:</span> <span fg="#f8fafc">{props.debugInfo.response.redirected ? "true" : "false"}</span></text>
                  <text><span fg="#94a3b8">Type:</span> <span fg="#f8fafc">{props.debugInfo.response.type}</span></text>
                  <text> </text>
                  <text><strong>Received Headers ({props.debugInfo.response.headerCount})</strong></text>
                  {props.debugInfo.response.headers.length > 0 ? (
                    props.debugInfo.response.headers.map((header, index) => (
                      <text key={`response-header-${index}`}>
                        <span fg="#93c5fd">{header.key}</span>
                        <span fg="#94a3b8">: </span>
                        <span fg="#f8fafc">{header.value}</span>
                      </text>
                    ))
                  ) : (
                    <text><span fg="#64748b">(none)</span></text>
                  )}
                </>
              ) : (
                <>
                  <text><span fg="#f87171">Request failed before a response was received.</span></text>
                  <text><span fg="#94a3b8">Error:</span> <span fg="#f8fafc">{props.debugInfo.errorMessage || "Unknown request error"}</span></text>
                </>
              )}
              <text> </text>
              <text><strong>Timing</strong></text>
              <text><span fg="#94a3b8">Started:</span> <span fg="#f8fafc">{props.debugInfo.startedAt}</span></text>
              <text><span fg="#94a3b8">Finished:</span> <span fg="#f8fafc">{props.debugInfo.finishedAt}</span></text>
              <text><span fg="#94a3b8">Duration:</span> <span fg="#f8fafc">{props.debugInfo.durationMs} ms</span></text>
            </>
          ) : (
            <>
              <text><strong>No debug data yet</strong></text>
              <text><span fg="#94a3b8">Send a request first, then run </span><span fg="#93c5fd">:debug</span><span fg="#94a3b8">.</span></text>
            </>
          )}
          <text> </text>
          <text><span fg="#64748b">Close debug: Esc, Enter, or q</span></text>
        </scrollbox>
      </box>
    </box>
  );
}
