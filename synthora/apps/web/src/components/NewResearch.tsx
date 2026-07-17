import { useEffect, useState } from "react";
import { api, PipelineSpec } from "../api";

export function NewResearch({
  onStarted,
}: {
  onStarted: (runId: string) => void;
}) {
  const [pipelines, setPipelines] = useState<PipelineSpec[]>([]);
  const [pipelineId, setPipelineId] = useState("deep_research");
  const [question, setQuestion] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.listPipelines().then(setPipelines).catch((e) => setError(String(e)));
  }, []);

  async function start() {
    setBusy(true);
    setError(null);
    try {
      const { run_id } = await api.startResearch(question.trim(), pipelineId);
      onStarted(run_id);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="panel">
      <h2>Ask a research question</h2>
      <textarea
        placeholder="What would you like to understand deeply?"
        value={question}
        onChange={(e) => setQuestion(e.target.value)}
        aria-label="research question"
      />
      <div className="pipeline-grid" role="radiogroup" aria-label="pipeline">
        {pipelines.map((p) => (
          <div
            key={p.id}
            role="radio"
            aria-checked={pipelineId === p.id}
            tabIndex={0}
            className={`pipeline-option ${pipelineId === p.id ? "selected" : ""}`}
            onClick={() => setPipelineId(p.id)}
            onKeyDown={(e) => e.key === "Enter" && setPipelineId(p.id)}
          >
            <h3>{p.name}</h3>
            <p>{p.description}</p>
          </div>
        ))}
      </div>
      {error && <p className="error-text">{error}</p>}
      <button
        className="primary"
        disabled={busy || question.trim().length < 3}
        onClick={start}
      >
        {busy ? "Starting…" : "Start research"}
      </button>
    </section>
  );
}
