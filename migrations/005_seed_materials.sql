-- Seed preset materials
-- These are system-wide materials available to all projects

INSERT INTO materials (id, project_id, name, category, absorption_bands_json, scattering_bands_json, is_preset, created_at, updated_at)
VALUES
('mat_preset_painted_concrete', NULL, 'Painted Concrete', 'Wall',
 '[{"frequency":63,"value":0.01},{"frequency":125,"value":0.01},{"frequency":250,"value":0.01},{"frequency":500,"value":0.02},{"frequency":1000,"value":0.02},{"frequency":2000,"value":0.02},{"frequency":4000,"value":0.03}]',
 '[{"frequency":63,"value":0.05},{"frequency":125,"value":0.05},{"frequency":250,"value":0.05},{"frequency":500,"value":0.05},{"frequency":1000,"value":0.05},{"frequency":2000,"value":0.05},{"frequency":4000,"value":0.05}]',
 1, datetime('now'), datetime('now')),
 
('mat_preset_gypsum_board', NULL, 'Gypsum Board', 'Wall',
 '[{"frequency":63,"value":0.05},{"frequency":125,"value":0.05},{"frequency":250,"value":0.06},{"frequency":500,"value":0.07},{"frequency":1000,"value":0.09},{"frequency":2000,"value":0.10},{"frequency":4000,"value":0.12}]',
 '[{"frequency":63,"value":0.05},{"frequency":125,"value":0.05},{"frequency":250,"value":0.05},{"frequency":500,"value":0.05},{"frequency":1000,"value":0.05},{"frequency":2000,"value":0.05},{"frequency":4000,"value":0.05}]',
 1, datetime('now'), datetime('now')),
 
('mat_preset_carpet', NULL, 'Carpet', 'Floor',
 '[{"frequency":63,"value":0.01},{"frequency":125,"value":0.05},{"frequency":250,"value":0.10},{"frequency":500,"value":0.20},{"frequency":1000,"value":0.35},{"frequency":2000,"value":0.50},{"frequency":4000,"value":0.60}]',
 '[{"frequency":63,"value":0.02},{"frequency":125,"value":0.02},{"frequency":250,"value":0.03},{"frequency":500,"value":0.03},{"frequency":1000,"value":0.04},{"frequency":2000,"value":0.04},{"frequency":4000,"value":0.05}]',
 1, datetime('now'), datetime('now')),
 
('mat_preset_glass', NULL, 'Glass', 'Window',
 '[{"frequency":63,"value":0.35},{"frequency":125,"value":0.20},{"frequency":250,"value":0.12},{"frequency":500,"value":0.07},{"frequency":1000,"value":0.04},{"frequency":2000,"value":0.03},{"frequency":4000,"value":0.02}]',
 '[{"frequency":63,"value":0.10},{"frequency":125,"value":0.10},{"frequency":250,"value":0.10},{"frequency":500,"value":0.10},{"frequency":1000,"value":0.10},{"frequency":2000,"value":0.10},{"frequency":4000,"value":0.10}]',
 1, datetime('now'), datetime('now')),
 
('mat_preset_heavy_curtain', NULL, 'Heavy Curtain', 'Window',
 '[{"frequency":63,"value":0.05},{"frequency":125,"value":0.10},{"frequency":250,"value":0.25},{"frequency":500,"value":0.55},{"frequency":1000,"value":0.65},{"frequency":2000,"value":0.70},{"frequency":4000,"value":0.70}]',
 '[{"frequency":63,"value":0.05},{"frequency":125,"value":0.05},{"frequency":250,"value":0.05},{"frequency":500,"value":0.05},{"frequency":1000,"value":0.05},{"frequency":2000,"value":0.05},{"frequency":4000,"value":0.05}]',
 1, datetime('now'), datetime('now')),
 
('mat_preset_wood_panel', NULL, 'Wood Panel', 'Wall',
 '[{"frequency":63,"value":0.10},{"frequency":125,"value":0.15},{"frequency":250,"value":0.20},{"frequency":500,"value":0.15},{"frequency":1000,"value":0.10},{"frequency":2000,"value":0.08},{"frequency":4000,"value":0.07}]',
 '[{"frequency":63,"value":0.15},{"frequency":125,"value":0.15},{"frequency":250,"value":0.15},{"frequency":500,"value":0.15},{"frequency":1000,"value":0.15},{"frequency":2000,"value":0.15},{"frequency":4000,"value":0.15}]',
 1, datetime('now'), datetime('now'));
