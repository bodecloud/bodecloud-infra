import { useEffect, useState } from "react";
import { api, RunSummary } from "../api";

export function History({ onOpen }: { onOpen: (runId: string) => void }) {
  const [runs, setRuns] = useState<RunSummary[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.listRuns().then(setRuns).catch((e) => setError(String(e)));
  }, []);

  return (
    <section className="panel">
      <h2>Research history</h2>
      {error && <p className="error-text">{error}</p>}
      {runs.length === 0 && !error && <p>No research yet.</p>}
      {runs.length > 0 && (
        <table className="runs">
          <thead>
            <tr>
              <th>Question</th>
              <th>Pipeline</th>
              <th>Status</th>
              <th>Started</th>
            </tr>
          </thead>
          <tbody>
            {runs.map((r) => (
              <tr key={r.id} onClick={() => onOpen(r.id)}>
                <td>{r.question}</td>
                <td>
                  <code>{r.pipeline_id}</code>
                </td>
                <td>
                  <span className={`status-badge status-${r.status}`}>
                    {r.status}
                  </span>
                </td>
                <td>{new Date(r.created_at).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}
