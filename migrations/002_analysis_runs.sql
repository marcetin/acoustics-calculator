-- Analysis runs and related tables

-- Analysis runs table
CREATE TABLE IF NOT EXISTS analysis_runs (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    profile TEXT NOT NULL DEFAULT 'STUDIO' CHECK (profile IN ('STUDIO', 'HIFI', 'PUBLIC')),
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'running', 'completed', 'failed')),
    inputs_hash TEXT NOT NULL,
    metrics_json TEXT NOT NULL DEFAULT '{}',
    started_at DATETIME,
    finished_at DATETIME,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Placement candidates table
CREATE TABLE IF NOT EXISTS placement_candidates (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    analysis_run_id TEXT,
    surface_id TEXT,
    center TEXT NOT NULL DEFAULT '{}',
    bounding_box TEXT NOT NULL DEFAULT '{}',
    orientation TEXT NOT NULL DEFAULT '{}',
    score REAL NOT NULL DEFAULT 0.0,
    confidence REAL NOT NULL DEFAULT 0.0,
    decision TEXT NOT NULL CHECK (decision IN ('DIFFUSER', 'ABSORBER_RECOMMENDED', 'REJECT')),
    diffuser_type_id TEXT,
    reasons_json TEXT NOT NULL DEFAULT '[]',
    warnings_json TEXT NOT NULL DEFAULT '[]',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (analysis_run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE,
    FOREIGN KEY (surface_id) REFERENCES surfaces(id) ON DELETE SET NULL,
    FOREIGN KEY (diffuser_type_id) REFERENCES diffuser_types(id) ON DELETE SET NULL
);

-- Create indexes for analysis tables
CREATE INDEX IF NOT EXISTS idx_analysis_runs_project_id ON analysis_runs(project_id);
CREATE INDEX IF NOT EXISTS idx_analysis_runs_status ON analysis_runs(status);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_project_id ON placement_candidates(project_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_analysis_run_id ON placement_candidates(analysis_run_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_surface_id ON placement_candidates(surface_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_decision ON placement_candidates(decision);
