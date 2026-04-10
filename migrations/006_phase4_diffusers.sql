-- Diffuser types catalog
CREATE TABLE IF NOT EXISTS diffuser_types (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('QRD_1D', 'QRD_2D', 'CUSTOM')),
    panel_width_m REAL NOT NULL,
    panel_height_m REAL NOT NULL,
    panel_depth_m REAL NOT NULL,
    min_effective_hz REAL NOT NULL,
    max_effective_hz REAL NOT NULL,
    scatter_axis TEXT NOT NULL CHECK (scatter_axis IN ('HORIZONTAL', 'VERTICAL', 'BOTH')),
    mount_orientation_json TEXT,
    cost_class TEXT,
    is_preset INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Placement candidates
CREATE TABLE IF NOT EXISTS placement_candidates (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    analysis_run_id TEXT NOT NULL,
    surface_id TEXT NOT NULL,
    bounding_box_json TEXT NOT NULL,
    center_x REAL NOT NULL,
    center_y REAL NOT NULL,
    center_z REAL NOT NULL,
    orientation TEXT NOT NULL,
    score REAL NOT NULL,
    confidence REAL NOT NULL,
    decision TEXT NOT NULL CHECK (decision IN ('DIFFUSER', 'ABSORBER_RECOMMENDED', 'REJECT')),
    diffuser_type_id TEXT,
    reasons_json TEXT NOT NULL,
    warnings_json TEXT NOT NULL,
    target_bands_json TEXT NOT NULL,
    source_receiver_pairs_json TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (analysis_run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE,
    FOREIGN KEY (surface_id) REFERENCES surfaces(id) ON DELETE CASCADE,
    FOREIGN KEY (diffuser_type_id) REFERENCES diffuser_types(id) ON DELETE SET NULL
);

-- Indexes for placement candidates
CREATE INDEX IF NOT EXISTS idx_placement_candidates_project_id ON placement_candidates(project_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_analysis_run_id ON placement_candidates(analysis_run_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_surface_id ON placement_candidates(surface_id);
CREATE INDEX IF NOT EXISTS idx_placement_candidates_decision ON placement_candidates(decision);

-- Seed preset diffuser types
INSERT OR IGNORE INTO diffuser_types (id, name, category, panel_width_m, panel_height_m, panel_depth_m, min_effective_hz, max_effective_hz, scatter_axis, is_preset) VALUES
('qrd-1d-600x600', 'QRD-1D 600x600', 'QRD_1D', 0.6, 0.6, 0.1, 500, 4000, 'BOTH', 1),
('qrd-1d-1200x600', 'QRD-1D 1200x600', 'QRD_1D', 1.2, 0.6, 0.1, 400, 3000, 'BOTH', 1),
('qrd-2d-600x600', 'QRD-2D 600x600', 'QRD_2D', 0.6, 0.6, 0.15, 600, 5000, 'BOTH', 1),
('qrd-2d-1200x1200', 'QRD-2D 1200x1200', 'QRD_2D', 1.2, 1.2, 0.15, 500, 4000, 'BOTH', 1),
('custom-diffuser-600', 'Custom Diffuser 600', 'CUSTOM', 0.6, 0.6, 0.2, 400, 6000, 'BOTH', 1);
