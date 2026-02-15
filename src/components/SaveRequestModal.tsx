type SaveRequestModalProps = {
  isOpen: boolean;
  requestName: string;
  setRequestName: (value: string) => void;
  onConfirm: () => void;
};

export function SaveRequestModal(props: SaveRequestModalProps) {
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
      zIndex={30}
      justifyContent="center"
      alignItems="center"
      backgroundColor="#0f172a"
    >
      <box border width={56} padding={1} title="Save Request" backgroundColor="#111827">
        <text>
          <span fg="#94a3b8">Name</span>
        </text>
        <box border marginTop={1} padding={1}>
          <input
            value={props.requestName}
            onChange={props.setRequestName}
            onInput={props.setRequestName}
            onSubmit={props.onConfirm}
            focused
          />
        </box>
        <text>
          <span fg="#64748b">Enter to save, Esc to cancel</span>
        </text>
      </box>
    </box>
  );
}
