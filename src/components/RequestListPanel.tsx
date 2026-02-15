import type { SavedRequest } from "../types";

type RequestListPanelProps = {
  requests: SavedRequest[];
  activeRequestId: string | null;
  focused: boolean;
  cursorIndex: number;
};

export function RequestListPanel(props: RequestListPanelProps) {
  return (
    <box border padding={1} title="Requests" width={28} minHeight={0}>
      <scrollbox focused={props.focused} flexGrow={1} minHeight={0} height="100%">
        {props.requests.length === 0 ? (
          <text>
            <span fg="#64748b">No saved requests yet.</span>
          </text>
        ) : (
          props.requests.map((request, index) => (
            <box key={request.id} marginBottom={1}>
              <text>
                <span fg={props.focused && index === props.cursorIndex ? "#f8fafc" : "#64748b"}>
                  {props.focused && index === props.cursorIndex ? "> " : "  "}
                </span>
                <span fg="#64748b">{index + 1}. </span>
                <span fg={request.id === props.activeRequestId ? "#22c55e" : "#64748b"}>{request.id === props.activeRequestId ? "* " : "  "}</span>
                <span fg="#e2e8f0">{request.name}</span>
              </text>
              <text>
                <span fg="#94a3b8">{request.method}</span>
                <span fg="#64748b"> </span>
                <span fg="#93c5fd">{request.url}</span>
              </text>
            </box>
          ))
        )}
      </scrollbox>
    </box>
  );
}
