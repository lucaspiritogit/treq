type CommandPanelProps = {
  commandMode: boolean;
  commandLine: string;
  commandFeedback: string;
};

export function CommandPanel(props: CommandPanelProps) {
  return (
    <box border padding={1} title="Command">
      {props.commandMode ? (
        <text>
          <span fg="#f8fafc">:</span>
          <span fg="#93c5fd">{props.commandLine || " "}</span>
        </text>
      ) : (
        <text>
          <span fg="#94a3b8">{props.commandFeedback || "Press :"}</span>
          <span fg="#64748b">{props.commandFeedback ? "" : " to open command mode (e.g. :send, :save, :debug, :help)"}</span>
        </text>
      )}
    </box>
  );
}
