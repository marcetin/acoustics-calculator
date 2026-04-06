-- Initial database schema for Acoustics Calculator

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    space_type TEXT NOT NULL CHECK (space_type IN ('STUDIO', 'HIFI', 'PUBLIC')),
    goal_type TEXT NOT NULL CHECK (goal_type IN ('MIXING', 'MASTERING', 'HIFI_STEREO', 'SPEECH', 'MUSIC', 'MULTIPURPOSE')),
    units TEXT NOT NULL DEFAULT 'METERS' CHECK (units IN ('METERS', 'FEET')),
    temperature_c REAL DEFAULT 20.0,
    humidity_percent REAL DEFAULT 50.0,
    speed_of_sound REAL DEFAULT 343.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Room geometries table
CREATE TABLE IF NOT EXISTS room_geometries (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    geometry_type TEXT NOT NULL CHECK (geometry_type IN ('SHOEBOX', 'POLYHEDRAL')),
    width REAL NOT NULL,
    length REAL NOT NULL,
    height REAL NOT NULL,
    volume_m3 REAL NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Materials table
CREATE TABLE IF NOT EXISTS materials (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    absorption_json TEXT NOT NULL DEFAULT '{}',
    scattering_json TEXT NOT NULL DEFAULT '{}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Surfaces table
CREATE TABLE IF NOT EXISTS surfaces (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('WALL', 'CEILING', 'FLOOR', 'OBJECT_FACE')),
    polygon3d TEXT NOT NULL DEFAULT '[]',
    normal TEXT NOT NULL DEFAULT '{}',
    area_m2 REAL NOT NULL,
    material_id TEXT,
    is_mountable BOOLEAN DEFAULT TRUE,
    mount_exclusion_zones TEXT NOT NULL DEFAULT '[]',
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE SET NULL
);

-- Sources table
CREATE TABLE IF NOT EXISTS sources (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('SPEAKER', 'POINT_SOURCE', 'LINE_ARRAY')),
    position TEXT NOT NULL DEFAULT '{}',
    aim TEXT NOT NULL DEFAULT '{}',
    gain_db REAL DEFAULT 0.0,
    group_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Receivers table
CREATE TABLE IF NOT EXISTS receivers (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('LISTENER', 'MIC', 'AUDIENCE_ZONE')),
    position TEXT NOT NULL DEFAULT '{}',
    weight REAL DEFAULT 1.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Diffuser types table
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
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_projects_space_type ON projects(space_type);
CREATE INDEX IF NOT EXISTS idx_projects_goal_type ON projects(goal_type);
CREATE INDEX IF NOT EXISTS idx_room_geometries_project_id ON room_geometries(project_id);
CREATE INDEX IF NOT EXISTS idx_materials_project_id ON materials(project_id);
CREATE INDEX IF NOT EXISTS idx_surfaces_project_id ON surfaces(project_id);
CREATE INDEX IF NOT EXISTS idx_surfaces_material_id ON surfaces(material_id);
CREATE INDEX IF NOT EXISTS idx_sources_project_id ON sources(project_id);
CREATE INDEX IF NOT EXISTS idx_receivers_project_id ON receivers(project_id);
