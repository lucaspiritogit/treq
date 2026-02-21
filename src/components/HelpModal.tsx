type HelpModalProps = {
  isOpen: boolean;
};

export function HelpModal(props: HelpModalProps) {
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
      zIndex={20}
      justifyContent="center"
      alignItems="center"
      backgroundColor="#0f172a"
    >
      <box border width="70%" height="70%" padding={1} title="Command Help" backgroundColor="#111827">
        <scrollbox focused flexGrow={1}>
          <text><strong>Vim-style Commands</strong></text>
          <text> </text>
          <text><span fg="#93c5fd">:send</span> or <span fg="#93c5fd">:s</span> - Send current HTTP request</text>
          <text><span fg="#93c5fd">:quit</span> or <span fg="#93c5fd">:q</span> - Exit app</text>
          <text><span fg="#93c5fd">:url</span> or <span fg="#93c5fd">:i</span> - Focus URL input</text>
          <text><span fg="#93c5fd">:headers</span> or <span fg="#93c5fd">:h</span> - Focus headers input</text>
          <text><span fg="#93c5fd">:r</span> or <span fg="#93c5fd">:req</span> - Focus request body input</text>
          <text><span fg="#93c5fd">:b</span> or <span fg="#93c5fd">:res</span> - Focus response body panel</text>
          <text><span fg="#93c5fd">:get</span> or <span fg="#93c5fd">:g</span> - Set method to GET</text>
          <text><span fg="#93c5fd">:post</span> or <span fg="#93c5fd">:p</span> - Set method to POST</text>
          <text><span fg="#93c5fd">:put</span> or <span fg="#93c5fd">:u</span> - Set method to PUT</text>
          <text><span fg="#93c5fd">:patch</span> or <span fg="#93c5fd">:t</span> - Set method to PATCH</text>
          <text><span fg="#93c5fd">:delete</span> or <span fg="#93c5fd">:d</span> - Set method to DELETE</text>
          <text><span fg="#93c5fd">:save</span> - Save current request (overwrites loaded request)</text>
          <text><span fg="#93c5fd">:list</span> - Focus request list sidebar</text>
          <text><span fg="#93c5fd">:toggle-list</span> or <span fg="#93c5fd">:tl</span> - Toggle request list sidebar</text>
          <text><span fg="#93c5fd">:reload</span> - Reload requests from requests.json</text>
          <text><span fg="#93c5fd">:help</span> - Open this help modal</text>
          <text><span fg="#93c5fd">:debug</span> - Open full-screen request/response debug modal</text>
          <text> </text>
          <text><span fg="#94a3b8">Interactive:</span> <span fg="#93c5fd">l</span> or <span fg="#93c5fd">Left</span> focuses request list</text>
          <text><span fg="#94a3b8">Request list:</span> <span fg="#93c5fd">Up/Down</span> navigate, <span fg="#93c5fd">Enter</span> loads request, <span fg="#93c5fd">Ctrl/Cmd+d</span> deletes request</text>
          <text> </text>
          <text><span fg="#94a3b8">Close help: Esc, Enter, or q</span></text>
        </scrollbox>
      </box>
    </box>
  );
}
