import type { SavedRequest } from "../types";

type DeleteRequestModalProps = {
  isOpen: boolean;
  requestToDelete: SavedRequest | null;
  onConfirm: () => void;
  onCancel: () => void;
};

export function DeleteRequestModal(props: DeleteRequestModalProps) {
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
      zIndex={35}
      justifyContent="center"
      alignItems="center"
      backgroundColor="#0f172a"
    >
      <box border width={64} padding={1} title="Delete Request" backgroundColor="#111827">
        <text><span fg="#fca5a5">Delete this request?</span></text>
        <text> </text>
        <text>
          <span fg="#94a3b8">Name: </span>
          <span fg="#e2e8f0">{props.requestToDelete?.name ?? "(unknown)"}</span>
        </text>
        <text>
          <span fg="#94a3b8">Method: </span>
          <span fg="#e2e8f0">{props.requestToDelete?.method ?? ""}</span>
        </text>
        <text>
          <span fg="#94a3b8">URL: </span>
          <span fg="#93c5fd">{props.requestToDelete?.url ?? ""}</span>
        </text>
        <text> </text>
        <text><span fg="#64748b">Enter/Y to confirm, Esc/N to cancel</span></text>
      </box>
    </box>
  );
}
