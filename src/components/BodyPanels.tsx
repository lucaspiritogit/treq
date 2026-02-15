import type { TextareaRenderable } from "@opentui/core";
import type { FocusField, UiMode } from "../types";
import { highlightJsonLine } from "../utils";

type BodyPanelsProps = {
  uiMode: UiMode;
  focusField: FocusField;
  requestBodyText: string;
  setRequestBodyText: (value: string) => void;
  requestBodyRef: React.RefObject<TextareaRenderable | null>;
  sendRequest: () => Promise<void>;
  isLoading: boolean;
  responseBody: string;
  responseStatus: string;
  stacked: boolean;
};

export function BodyPanels(props: BodyPanelsProps) {
  const requestLines = props.requestBodyText.split("\n");
  const responseContent = props.isLoading ? "Waiting for response..." : props.responseBody;
  const responseLines = responseContent.split("\n");

  return (
    <box flexDirection={props.stacked ? "column" : "row"} flexGrow={1} minHeight={0} gap={1}>
      <box border padding={1} title="Request Body" flexGrow={1} flexShrink={1} flexBasis={0} minWidth={0} minHeight={0}>
        {props.uiMode === "input" && props.focusField === "requestBody" ? (
          <textarea
            ref={props.requestBodyRef}
            initialValue={props.requestBodyText}
            onContentChange={() => {
              props.setRequestBodyText(props.requestBodyRef.current?.plainText ?? "");
            }}
            onSubmit={() => {
              void props.sendRequest();
            }}
            keyBindings={[
              { name: "return", ctrl: true, action: "submit" },
              { name: "enter", ctrl: true, action: "submit" },
            ]}
            focused
            wrapMode="word"
            width="100%"
            flexGrow={1}
            flexShrink={1}
            minWidth={0}
            minHeight={0}
            height="100%"
          />
        ) : (
          <scrollbox flexGrow={1} minHeight={0} height="100%">
            {requestLines.map((line, lineIndex) => (
              <text key={`request-${lineIndex}`}>
                {highlightJsonLine(line).map((token, tokenIndex) => {
                  if (token.bold) {
                    return (
                      <span key={`request-${lineIndex}-${tokenIndex}`} fg={token.color}>
                        <strong>{token.text}</strong>
                      </span>
                    );
                  }

                  if (token.italic) {
                    return (
                      <span key={`request-${lineIndex}-${tokenIndex}`} fg={token.color}>
                        <em>{token.text}</em>
                      </span>
                    );
                  }

                  return <span key={`request-${lineIndex}-${tokenIndex}`} fg={token.color}>{token.text}</span>;
                })}
              </text>
            ))}
          </scrollbox>
        )}
      </box>
      <box border padding={1} title={`Response Body (${props.responseStatus})`} flexGrow={1} flexShrink={1} flexBasis={0} minWidth={0} minHeight={0}>
        <scrollbox focused={props.focusField === "responseBody"} flexGrow={1} minHeight={0} height="100%">
          {responseLines.map((line, lineIndex) => (
            <text key={`response-${lineIndex}`}>
              {highlightJsonLine(line).map((token, tokenIndex) => {
                if (token.bold) {
                  return (
                    <span key={`response-${lineIndex}-${tokenIndex}`} fg={token.color}>
                      <strong>{token.text}</strong>
                    </span>
                  );
                }

                if (token.italic) {
                  return (
                    <span key={`response-${lineIndex}-${tokenIndex}`} fg={token.color}>
                      <em>{token.text}</em>
                    </span>
                  );
                }

                return <span key={`response-${lineIndex}-${tokenIndex}`} fg={token.color}>{token.text}</span>;
              })}
            </text>
          ))}
        </scrollbox>
      </box>
    </box>
  );
}
