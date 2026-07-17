import { useState } from "react";
import { History } from "./components/History";
import { NewResearch } from "./components/NewResearch";
import { RunView } from "./components/RunView";
import { Settings } from "./components/Settings";

type View =
  | { name: "new" }
  | { name: "history" }
  | { name: "settings" }
  | { name: "run"; runId: string };

export function App() {
  const [view, setView] = useState<View>({ name: "new" });

  return (
    <>
      <header className="masthead">
        <h1>
          Syn<span>thora</span>
        </h1>
        <nav>
          <button
            className={view.name === "new" ? "active" : ""}
            onClick={() => setView({ name: "new" })}
          >
            New research
          </button>
          <button
            className={view.name === "history" ? "active" : ""}
            onClick={() => setView({ name: "history" })}
          >
            History
          </button>
          <button
            className={view.name === "settings" ? "active" : ""}
            onClick={() => setView({ name: "settings" })}
          >
            Settings
          </button>
        </nav>
      </header>
      <main>
        {view.name === "new" && (
          <NewResearch onStarted={(runId) => setView({ name: "run", runId })} />
        )}
        {view.name === "history" && (
          <History onOpen={(runId) => setView({ name: "run", runId })} />
        )}
        {view.name === "settings" && <Settings />}
        {view.name === "run" && <RunView runId={view.runId} />}
      </main>
    </>
  );
}
