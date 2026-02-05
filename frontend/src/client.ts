import { TodoService } from "./gen/todo/v1/todo_pb";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

// Use environment variable or default to backend service
const baseUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";

const transport = createConnectTransport({
  baseUrl,
});

export const todoClient = createClient(TodoService, transport);
