import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { TodoService } from "./gen/todo/v1/todo_pb";
import { EventService } from "./gen/event/v1/event_pb";

const baseUrl = import.meta.env.VITE_API_URL || "";

const transport = createConnectTransport({
  baseUrl,
});

export const todoClient = createClient(TodoService, transport);
export const eventClient = createClient(EventService, transport);
