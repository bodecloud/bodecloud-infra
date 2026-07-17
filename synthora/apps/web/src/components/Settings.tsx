import { useEffect, useState } from "react";
import { api, Providers } from "../api";

export function Settings() {
  const [providers, setProviders] = useState<Providers | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.listProviders().then(setProviders).catch((e) => setError(String(e)));
  }, []);

  return (
    <section className="panel">
      <h2>Providers</h2>
      {error && <p className="error-text">{error}</p>}
      {providers && (
        <>
          <h3>LLM providers</h3>
          <div className="provider-list">
            {providers.llm_providers.map((p) => (
              <code key={p}>{p}</code>
            ))}
          </div>
          <h3>Search engines</h3>
          <div className="provider-list">
            {providers.search_engines.map((p) => (
              <code key={p}>{p}</code>
            ))}
          </div>
          <h3>Search strategies</h3>
          <div className="provider-list">
            {providers.search_strategies.map((p) => (
              <code key={p}>{p}</code>
            ))}
          </div>
          <p>
            Model and engine defaults are set per deployment via environment
            variables, and per run via the API <code>config</code> payload.
          </p>
        </>
      )}
    </section>
  );
}
