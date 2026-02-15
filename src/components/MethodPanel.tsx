import type { FocusField, HttpMethod, UiMode } from "../types";

type MethodPanelProps = {
  focusField: FocusField;
  uiMode: UiMode;
  method: HttpMethod;
  methodColors: Record<HttpMethod, string>;
};

export function MethodPanel(props: MethodPanelProps) {
  return (
    <box border padding={1} title="Method" width={14}>
      <text>
        <span fg={props.focusField === "method" && props.uiMode === "interactive" ? "#f8fafc" : "#94a3b8"}>
          {props.focusField === "method" && props.uiMode === "interactive" ? "> " : "  "}
        </span>
        <span fg={props.methodColors[props.method]}><strong>{props.method}</strong></span>
      </text>
    </box>
  );
}
