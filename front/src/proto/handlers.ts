export interface WsEventHandler {
  handler: CallableFunction;
  decodeFn: CallableFunction;
  kwargs: Record<string, CallableFunction>;
}

// Event IDs
export const WS_CHALLENGE_ATTEMPT = 0x01;
export const WS_CHALLENGE_REVEAL = 0x02;

let handlers: Record<number, WsEventHandler> = {};

export function registerEventHandler(
  eventId: number,
  handler: CallableFunction,
  decodeFn: CallableFunction,
  kwargs: Record<string, any>
) {
  handlers[eventId] = { handler, kwargs, decodeFn };
}

export function unregisterEventHandler(eventId: number) {
  delete handlers[eventId];
}

export function clearEventHandlers() {
  handlers = {};
}

export function callHandler(eventId: number, data: Uint8Array) {
  const handler = handlers[eventId];
  if (handler) {
    handler.handler(handler.decodeFn(data.slice(1)), handler.kwargs);
  } else {
    console.error(`[event] no handler found for event id 0x${eventId.toString(16)}`);
  }
}
