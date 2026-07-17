import { useCallback, useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import {
  api,
  Citation,
  KnowledgeEdge,
  KnowledgeNode,
  RunDetail,
  TERMINAL_STATUSES,
} from "../api";
import { useRunEvents } from "../hooks/useRunEvents";
import { EventFeed } from "./EventFeed";
import { KnowledgeMapView } from "./KnowledgeMapView";

export function RunView({ runId }: { runId: string }) {
  const [run, setRun] = useState<RunDetail | null>(null);
  const [report, setReport] = useState<string | null>(null);
  const [citations, setCitations] = useState<Citation[]>([]);
  const [kmap, setKmap] = useState<{
    nodes: KnowledgeNode[];
    edges: KnowledgeEdge[];
  } | null>(null);
  const [steer, setSteer] = useState("");
  const { events, finished } = useRunEvents(runId);

  const refresh = useCallback(async () => {
    const detail = await api.getRun(runId);
    setRun(detail);
    if (detail.status === "completed") {
      try {
        const r = await api.getReport(runId);
        setReport(r.report_markdown);
        setCitations(r.citations);
      } catch {
        /* report may not exist for cancelled runs */
      }
      try {
        const m = await api.getKnowledgeMap(runId);
        if (m.nodes.length) setKmap(m);
      } catch {
        /* knowledge map optional */
      }
    }
  }, [runId]);

  useEffect(() => {
    refresh();
  }, [refresh, finished]);

  const running = run != null && !TERMINAL_STATUSES.includes(run.status);

  return (
    <>
      <section className="panel">
        <h2>{run?.question ?? "Loading…"}</h2>
        {run && (
          <p>
            <span className={`status-badge status-${run.status}`}>
              {run.status}
            </span>{" "}
            <code>{run.pipeline_id}</code>
          </p>
        )}
        {run?.brief && <p>{run.brief}</p>}
        {run?.error && <p className="error-text">{run.error}</p>}
        {running && (
          <>
            <button className="ghost" onClick={() => api.cancelRun(runId)}>
              Cancel run
            </button>
            <div className="steer-row">
              <input
                type="text"
                placeholder="Steer the research (e.g. 'focus on costs')"
                value={steer}
                onChange={(e) => setSteer(e.target.value)}
                aria-label="steering message"
              />
              <button
                className="primary"
                disabled={!steer.trim()}
                onClick={() => {
                  api.steerRun(runId, steer.trim());
                  setSteer("");
                }}
              >
                Steer
              </button>
            </div>
          </>
        )}
      </section>

      <section className="panel">
        <h2>Progress</h2>
        <EventFeed events={events} />
      </section>

      {report && (
        <section className="panel">
          <h2>Report</h2>
          <div className="report-body">
            <ReactMarkdown>{report}</ReactMarkdown>
          </div>
          {citations.length > 0 && (
            <details>
              <summary>{citations.length} citations</summary>
              <ol>
                {citations
                  .filter((c) => c.index != null)
                  .sort((a, b) => (a.index ?? 0) - (b.index ?? 0))
                  .map((c) => (
                    <li key={c.id} value={c.index ?? undefined}>
                      <a href={c.url} target="_blank" rel="noreferrer">
                        {c.title || c.url}
                      </a>
                      {!c.verified && " (unverified)"}
                    </li>
                  ))}
              </ol>
            </details>
          )}
        </section>
      )}

      {kmap && (
        <section className="panel">
          <h2>Knowledge map</h2>
          <KnowledgeMapView nodes={kmap.nodes} />
        </section>
      )}
    </>
  );
}
