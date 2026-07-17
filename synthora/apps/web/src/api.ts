// Typed client for the Synthora REST + WebSocket API.

export interface PipelineSpec {
  id: string;
  name: string;
  description: string;
  tags: string[];
}

export interface RunSummary {
  id: string;
  question: string;
  pipeline_id: string;
  status: string;
  created_at: string;
  finished_at: string | null;
}

export interface RunDetail extends RunSummary {
  brief: string | null;
  error: string | null;
  config: Record<string, unknown>;
  started_at: string | null;
}

export interface RunEvent {
  run_id: string;
  type: string;
  message: string;
  node: string | null;
  payload: Record<string, unknown>;
  timestamp: string;
}

export interface Citation {
  id: string;
  url: string;
  title: string;
  snippet: string;
  confidence: number;
  index: number | null;
  verified: boolean;
}

export interface KnowledgeNode {
  id: string;
  name: string;
  summary: string;
  parent_id: string | null;
  infos: Citation[];
}

export interface KnowledgeEdge {
  id: string;
  source_id: string;
  target_id: string;
  relation: string;
}

export interface Providers {
  llm_providers: string[];
  search_engines: string[];
  search_strategies: string[];
}

let authToken: string | null = null;
export function setToken(token: string | null) {
  authToken = token;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init?.headers as Record<string, string>),
  };
  if (authToken) headers["Authorization"] = `Bearer ${authToken}`;
  const resp = await fetch(path, { ...init, headers });
  if (!resp.ok) {
    const body = await resp.text();
    throw new Error(`${resp.status}: ${body}`);
  }
  return resp.json() as Promise<T>;
}

export const api = {
  listPipelines: () =>
    request<{ pipelines: PipelineSpec[] }>("/api/v1/pipelines").then(
      (d) => d.pipelines,
    ),
  listProviders: () => request<Providers>("/api/v1/providers"),
  listRuns: () =>
    request<{ runs: RunSummary[] }>("/api/v1/research").then((d) => d.runs),
  getRun: (id: string) => request<RunDetail>(`/api/v1/research/${id}`),
  startResearch: (
    question: string,
    pipelineId: string,
    config?: Record<string, unknown>,
  ) =>
    request<{ run_id: string; status: string }>("/api/v1/research", {
      method: "POST",
      body: JSON.stringify({ question, pipeline_id: pipelineId, config }),
    }),
  cancelRun: (id: string) =>
    request(`/api/v1/research/${id}/cancel`, { method: "POST", body: "{}" }),
  steerRun: (id: string, message: string) =>
    request(`/api/v1/research/${id}/steer`, {
      method: "POST",
      body: JSON.stringify({ message }),
    }),
  getReport: (id: string) =>
    request<{
      report_markdown: string;
      citations: Citation[];
      status: string;
    }>(`/api/v1/research/${id}/report`),
  getEvents: (id: string) =>
    request<{ events: RunEvent[] }>(`/api/v1/research/${id}/events`).then(
      (d) => d.events,
    ),
  getKnowledgeMap: (id: string) =>
    request<{ nodes: KnowledgeNode[]; edges: KnowledgeEdge[] }>(
      `/api/v1/research/${id}/knowledge-map`,
    ),
};

export function eventsSocketUrl(runId: string): string {
  const proto = window.location.protocol === "https:" ? "wss" : "ws";
  return `${proto}://${window.location.host}/api/v1/research/${runId}/events/ws`;
}

export const TERMINAL_STATUSES = ["completed", "failed", "cancelled"];
