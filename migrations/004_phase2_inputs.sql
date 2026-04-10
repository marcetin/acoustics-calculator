-- PHASE 2: Geometry + Materials + Inputs
-- Tables for room geometry, surfaces, materials, sources, receivers, and constraints

-- Drop existing tables that will be recreated
DROP TABLE IF EXISTS constraint_sets;
DROP TABLE IF EXISTS receivers;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS materials;
DROP TABLE IF EXISTS surfaces;
DROP TABLE IF EXISTS room_geometries;

-- Room geometries table (enhanced from PHASE 0)
CREATE TABLE room_geometries (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    geometry_type TEXT NOT NULL CHECK (geometry_type IN ('SHOEBOX', 'POLYHEDRAL')),
    width REAL NOT NULL CHECK (width > 0),
    length REAL NOT NULL CHECK (length > 0),
    height REAL NOT NULL CHECK (height > 0),
    volume_m3 REAL NOT NULL DEFAULT 0.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Surfaces table (enhanced from PHASE 0)
CREATE TABLE surfaces (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    geometry_id TEXT,
    name TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('WALL', 'CEILING', 'FLOOR', 'OBJECT_FACE')),
    polygon_3d_json TEXT,
    normal_x REAL NOT NULL DEFAULT 0.0,
    normal_y REAL NOT NULL DEFAULT 0.0,
    normal_z REAL NOT NULL DEFAULT 0.0,
    area_m2 REAL NOT NULL DEFAULT 0.0,
    material_id TEXT,
    is_mountable BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (geometry_id) REFERENCES room_geometries(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE SET NULL
);

-- Materials table (enhanced from PHASE 0)
CREATE TABLE materials (
    id TEXT PRIMARY KEY,
    project_id TEXT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    absorption_bands_json TEXT NOT NULL,
    scattering_bands_json TEXT NOT NULL,
    is_preset BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Sources table (enhanced from PHASE 0)
CREATE TABLE sources (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('SPEAKER', 'POINT_SOURCE', 'LINE_ARRAY')),
    position_x REAL NOT NULL DEFAULT 0.0,
    position_y REAL NOT NULL DEFAULT 0.0,
    position_z REAL NOT NULL DEFAULT 0.0,
    aim_x REAL NOT NULL DEFAULT 0.0,
    aim_y REAL NOT NULL DEFAULT 0.0,
    aim_z REAL NOT NULL DEFAULT 0.0,
    gain_db REAL NOT NULL DEFAULT 0.0,
    group_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Receivers table (enhanced from PHASE 0)
CREATE TABLE receivers (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('LISTENER', 'MIC', 'AUDIENCE_ZONE')),
    position_x REAL NOT NULL DEFAULT 0.0,
    position_y REAL NOT NULL DEFAULT 0.0,
    position_z REAL NOT NULL DEFAULT 0.0,
    weight REAL NOT NULL DEFAULT 1.0 CHECK (weight > 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Constraint sets table (enhanced from PHASE 0)
CREATE TABLE constraint_sets (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL UNIQUE,
    max_panel_depth_m REAL NOT NULL DEFAULT 0.3 CHECK (max_panel_depth_m > 0),
    symmetry_required BOOLEAN NOT NULL DEFAULT 0,
    allowed_surface_kinds_json TEXT NOT NULL,
    forbidden_surface_ids_json TEXT NOT NULL,
    preferred_treatment_mode TEXT NOT NULL DEFAULT 'BALANCED' CHECK (preferred_treatment_mode IN ('BALANCED', 'DIFFUSION_FIRST', 'ABSORPTION_FIRST')),
    budget_class TEXT NOT NULL DEFAULT 'MEDIUM' CHECK (budget_class IN ('LOW', 'MEDIUM', 'HIGH')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_room_geometries_project_id ON room_geometries(project_id);
CREATE INDEX idx_surfaces_project_id ON surfaces(project_id);
CREATE INDEX idx_surfaces_geometry_id ON surfaces(geometry_id);
CREATE INDEX idx_surfaces_material_id ON surfaces(material_id);
CREATE INDEX idx_materials_project_id ON materials(project_id);
CREATE INDEX idx_sources_project_id ON sources(project_id);
CREATE INDEX idx_receivers_project_id ON receivers(project_id);
CREATE INDEX idx_constraint_sets_project_id ON constraint_sets(project_id);
