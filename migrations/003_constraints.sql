-- Constraints and seed data

-- Constraint sets table
CREATE TABLE IF NOT EXISTS constraint_sets (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    max_panel_depth_m REAL DEFAULT 0.3,
    allowed_surface_kinds TEXT NOT NULL DEFAULT '["WALL","CEILING"]',
    forbidden_surface_ids TEXT NOT NULL DEFAULT '[]',
    symmetry_required BOOLEAN DEFAULT FALSE,
    budget_class TEXT DEFAULT 'STANDARD' CHECK (budget_class IN ('BASIC', 'STANDARD', 'PREMIUM')),
    preferred_treatment_mode TEXT DEFAULT 'BALANCED' CHECK (preferred_treatment_mode IN ('DIFFUSION_FOCUSED', 'ABSORPTION_FOCUSED', 'BALANCED')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Seed data for diffuser types
INSERT OR IGNORE INTO diffuser_types (id, name, category, panel_width_m, panel_height_m, panel_depth_m, min_effective_hz, max_effective_hz, scatter_axis) VALUES
('QRD1D_600x1200', 'QRD 1D 600x1200', 'QRD_1D', 0.6, 1.2, 0.15, 500, 4000, 'HORIZONTAL'),
('QRD2D_600x600', 'QRD 2D 600x600', 'QRD_2D', 0.6, 0.6, 0.12, 800, 5000, 'BOTH'),
('QRD2D_1200x1200', 'QRD 2D 1200x1200', 'QRD_2D', 1.2, 1.2, 0.2, 400, 3000, 'BOTH'),
('CUSTOM_GENERIC', 'Generic Custom Diffuser', 'CUSTOM', 0.6, 0.6, 0.1, 600, 8000, 'BOTH');

-- Create index for constraints
CREATE INDEX IF NOT EXISTS idx_constraint_sets_project_id ON constraint_sets(project_id);
