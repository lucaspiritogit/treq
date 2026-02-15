import type { UiMode } from "../types";

type ModePanelProps = {
  uiMode: UiMode;
};

export function ModePanel(props: ModePanelProps) {
  return (
    <text>
      <span fg="#64748b">mode: </span>
      <span fg={props.uiMode === "interactive" ? "#cbd5e1" : "#94a3b8"}>
        {props.uiMode === "interactive" ? "INTERACTIVE" : "INPUT"}
      </span>
    </text>
  );
}
