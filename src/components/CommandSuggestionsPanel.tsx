import { commandSuggestions } from "../constants";

type CommandSuggestionsPanelProps = {
  commandMode: boolean;
  commandLine: string;
};

export function CommandSuggestionsPanel(props: CommandSuggestionsPanelProps) {
  if (!props.commandMode) {
    return null;
  }

  const query = props.commandLine.trim().toLowerCase().replace(/^:/, "");
  const matches = commandSuggestions.filter((item) => item.name.startsWith(query)).slice(0, 5);

  return (
    <box border padding={1} title="Commands" backgroundColor="#0b1220">
      {matches.length === 0 ? (
        <text>
          <span fg="#64748b">No matching commands</span>
        </text>
      ) : (
        <box flexDirection="column" gap={0}>
          {matches.map((item) => (
            <box key={item.name} marginBottom={1}>
              <box flexDirection="row" width="100%" minWidth={0}>
                <text>
                  <span fg="#93c5fd">:{item.name}</span>
                </text>
                <box flexGrow={1} />
                <text>
                  <span fg="#64748b">{item.description}</span>
                </text>
              </box>
              <text>
                <span fg="#0b1220"> </span>
              </text>
            </box>
          ))}
        </box>
      )}
    </box>
  );
}
