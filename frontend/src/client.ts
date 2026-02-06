import { TodoService } from "./gen/todo/v1/todo_pb";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

// Use empty string to let Vite proxy handle requests in dev, or use env variable for production
const baseUrl = import.meta.env.VITE_API_URL || "";

const transport = createConnectTransport({
  baseUrl,
});

export const todoClient = createClient(TodoService, transport);
