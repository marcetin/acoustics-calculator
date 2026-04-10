# Acoustics Calculator - Development Phases

## Overview

This document outlines the phased development approach for the Acoustics Calculator project. Each phase builds incrementally on the previous one while maintaining a buildable, functional application at every step.

---

## PHASE 0 - Bootstrap ✅ COMPLETED

**Goal**: Establish project foundation and basic infrastructure.

### What Was Implemented
- ✅ Go project structure with proper package organization
- ✅ HTTP server using chi router
- ✅ SQLite database with modernc.org/sqlite (CGO-free)
- ✅ Database migration system
- ✅ Basic HTML template system
- ✅ CSS styling and static file serving
- ✅ Makefile with build, run, test commands
- ✅ README with setup instructions

### Key Files Added/Modified
- `cmd/server/main.go` - Application entrypoint
- `internal/app/` - Application configuration and bootstrap
- `internal/http/router.go` - HTTP routing
- `internal/repo/` - Database layer with migrations
- `templates/layouts/base.html` - Base template
- `templates/pages/dashboard.html` - Dashboard page
- `static/css/app.css` - Application styling
- `Makefile` - Build automation
- `README.md` - Project documentation

### Database Schema
- `projects` - Core project information
- `room_geometries` - Room dimensions
- `surfaces` - Surface materials
- `materials` - Material properties
- `sources` - Sound source positions
- `receivers` - Listening positions
- `constraint_sets` - Treatment constraints
- `diffuser_types` - Diffuser catalog
- `analysis_runs` - Analysis execution
- `placement_candidates` - Treatment recommendations

---

## PHASE 1 - Project Wizard + Core CRUD ✅ COMPLETED

**Goal**: Implement complete project management workflow with HTMX tab interface.

### What Was Implemented
- ✅ Project domain model with enums (SpaceType, GoalType, Units)
- ✅ Project repository with full CRUD operations
- ✅ Project service layer with validation
- ✅ HTTP handlers for project management
- ✅ Project creation form with validation
- ✅ Project workspace with HTMX tab navigation
- ✅ Dashboard with project listing
- ✅ Complete tab system with empty states
- ✅ Error handling and form validation
- ✅ Basic unit tests for domain and service layers

### Key Features
- **Project Creation**: Full form with space type, goal type, environmental conditions
- **Project Workspace**: Tab-based interface for different aspects
- **HTMX Navigation**: Partial rendering for tab content
- **Validation**: Comprehensive input validation with error messages
- **Empty States**: Professional "coming soon" content for unimplemented features

### Tab Structure
1. **Summary** - Project overview and setup progress
2. **Geometry** - Room dimensions (Phase 2)
3. **Surfaces** - Material assignments (Phase 2)
4. **Sources** - Sound source positions (Phase 2)
5. **Receivers** - Listening positions (Phase 2)
6. **Constraints** - Treatment limitations (Phase 2)
7. **Analysis** - Acoustic calculations (Phase 3)
8. **Placements** - Treatment recommendations (Phase 4)

### Routes Implemented
- `GET /` - Dashboard (redirects to /dashboard)
- `GET /dashboard` - Project listing
- `GET /projects/new` - New project form
- `POST /projects` - Create project
- `GET /projects/{id}` - Project workspace
- `GET /hx/projects/{id}/tabs/{tab}` - HTMX tab content

---

## PHASE 2 - Geometry + Materials + Inputs ✅ COMPLETED

**Goal**: Implement complete input system for acoustic analysis.

### What Was Implemented
- ✅ Room geometry configuration (shoebox only)
- ✅ Surface material assignments
- ✅ 2D room editor with drag-and-drop (X/Y plane)
- ✅ Source positioning system with CRUD
- ✅ Receiver positioning system with CRUD
- ✅ Constraint configuration
- ✅ Material library with preset materials (absorption/scattering coefficients)
- ✅ Geometry validation
- ✅ Surface auto-generation for shoebox rooms
- ✅ Summary tab with real-time readiness status
- ✅ Unit tests for geometry, source, and receiver services

### Key Features
- **Shoebox Geometry**: Define room dimensions (width, length, height) with automatic surface generation
- **Material Assignment**: Assign preset materials to surfaces via dropdown with real-time updates
- **2D Room Editor**: SVG-based top-down view with drag-and-drop for source/receiver positioning
- **Source/Receiver CRUD**: Create, list, and delete sound sources and listening positions
- **Constraints**: Configure treatment preferences (panel depth, symmetry, budget class)
- **Validation**: Comprehensive validation for geometry dimensions, positions within room, material assignments
- **Summary Tab**: Real-time progress tracking showing configuration status and readiness for analysis

### Technical Implementation
- SVG-based room editor with drag-and-drop on X/Y plane
- Z-coordinate edited via forms only
- Backend validation for all position changes
- Material database with frequency-dependent absorption and scattering coefficients
- HTMX partial rendering for all tab updates
- Repository and service layers with proper error handling

### Limitations
- **Shoebox Only**: Polyhedral geometry is scaffolded but not implemented
- **2D Editor**: Top-down view only, Z edited via forms
- **No Export/Import**: Room layout export/import not implemented
- **Preset Materials Only**: Custom material creation not exposed in UI

### Domain Models Added
- `RoomGeometry` - Room dimensions and type (SHOEBOX/POLYHEDRAL)
- `Surface` - Room surfaces with materials and properties
- `Material` - Acoustic materials with absorption/scattering bands
- `Source` - Sound source positions and properties
- `Receiver` - Listening positions and weights
- `ConstraintSet` - Treatment constraints and preferences

### Database Migrations
- `004_phase2_inputs.sql` - Tables for all phase 2 entities
- `005_seed_materials.sql` - Preset material data (drywall, carpet, wood, etc.)

### Routes Added
- `POST /projects/{id}/geometry` - Save room geometry
- `POST /projects/{id}/surfaces/{surfaceID}` - Update surface material
- `POST /projects/{id}/sources` - Create source
- `POST /projects/{id}/sources/{sourceID}/delete` - Delete source
- `POST /projects/{id}/receivers` - Create receiver
- `POST /projects/{id}/receivers/{receiverID}/delete` - Delete receiver
- `POST /projects/{id}/constraints` - Save constraints
- `GET /hx/projects/{id}/editor` - Render 2D editor
- `POST /hx/projects/{id}/editor/sources/{sourceID}/move` - Move source via editor
- `POST /hx/projects/{id}/editor/receivers/{receiverID}/move` - Move receiver via editor

---

## PHASE 3 - Acoustics Engine V1 ✅ COMPLETED

**Goal**: Implement core acoustic analysis capabilities.

### What Was Implemented
- ✅ Rectangular modal solver (shoebox only)
- ✅ Image-source early reflections (1st/2nd order)
- ✅ RT60 calculation (Sabine & Eyring)
- ✅ Frequency band analysis (octave bands)
- ✅ Basic metrics calculation
- ✅ Analysis run persistence
- ✅ Analysis UI with results display
- ✅ Prerequisites validation
- ✅ Warnings system
- ✅ Unit tests for acoustics engine

### Key Features
- **Modal Analysis**: Calculates axial, tangential, and oblique modes up to 300Hz for shoebox rooms
- **Speed of Sound**: Formula-based calculation with temperature adjustment (331.3 + 0.606*tempC)
- **RT Estimation**: Sabine and Eyring T60 calculations across 7 octave bands (63-4000Hz)
- **Material Fallback**: Default absorption coefficients when materials are missing
- **Image-Source**: First and second order reflection path generation with energy estimates
- **Analysis Persistence**: AnalysisRun model with status tracking (queued, running, completed, failed)
- **Results UI**: Tables for modes, RT bands, and reflections with warnings display
- **Prerequisites Check**: Blocks analysis if geometry, sources, or receivers missing

### Technical Implementation
- Modal solver uses standard rectangular room mode formula: f = (c/2) * sqrt((l/Lx)^2 + (m/Ly)^2 + (n/Lz)^2)
- Modes sorted by ascending frequency
- RT estimation uses equivalent absorption area calculation
- Image-source method with 6 surface reflections (walls, floor, ceiling)
- Second-order reflections via double image-source transformation
- Validation ensures only shoebox geometry is supported
- Warnings persisted in metrics_json for user visibility

### Limitations
- **Shoebox Only**: Polyhedral geometry not supported
- **1st/2nd Order Only**: Higher-order reflections not calculated
- **No Async Jobs**: Analysis runs synchronously
- **No Ray Tracing**: Simple image-source method only
- **Temperature Only**: Humidity ignored in speed-of-sound calculation
- **Preset Materials Only**: Uses preset material library

### Domain Models Added
- `AnalysisRun` - Analysis execution with status and metrics
- `AnalysisMetrics` - Complete analysis results (modes, RT, reflections, warnings)
- `Mode` - Room mode with (l,m,n), frequency, and classification
- `RTResults` - RT60 results per frequency band
- `RTBand` - Single frequency band RT data
- `ReflectionPath` - Early reflection path with sequence and energy

### Database Migrations
- No new migrations (analysis_runs table already existed in 002_analysis_runs.sql)

### Routes Added
- `POST /projects/{id}/analysis/run` - Run acoustic analysis

### Files Added
- `internal/domain/analysis_run.go` - Analysis run domain model
- `internal/acoustics/modal.go` - Modal solver
- `internal/acoustics/rt.go` - RT estimator
- `internal/acoustics/reflections.go` - Image-source solver
- `internal/repo/analysis_repo.go` - Analysis repository
- `internal/service/analysis_service.go` - Analysis service
- `internal/http/handlers/analysis.go` - Analysis HTTP handler
- `internal/acoustics/modal_test.go` - Modal solver tests
- `internal/acoustics/rt_test.go` - RT estimator tests
- `internal/acoustics/reflections_test.go` - Image-source tests

### Files Modified
- `internal/repo/repositories.go` - Added AnalysisRepository
- `internal/app/app.go` - Added AnalysisService and firstN template function
- `internal/http/router.go` - Added analysis handler and route
- `internal/http/handlers/projects.go` - Load analysis data for analysis tab
- `templates/partials/tab_analysis.html` - Replaced placeholder with functional UI

---

## PHASE 4 - Diffuser Candidate Engine V1 

**Goal**: Generate intelligent treatment recommendations.

### What Was Implemented
- Candidate zone generation from reflection analysis
- Scoring and ranking system with weighted components
- Veto logic for mountability, depth, nearfield, first-reflection
- Diffuser type assignment (QRD_1D, QRD_2D, CUSTOM)
- Simple layout synthesis for panel coordinates
- Placement persistence linked to AnalysisRun
- Placements UI with prerequisites and results
- Profile-specific rules (STUDIO, HIFI, PUBLIC)
- Unit tests for scoring, veto, catalog, layout

### Key Features
- **Candidate Extraction**: Aggregates reflections by surface, creates patches based on hit concentration
- **Scoring Engine**: Weighted components (energy, time, receiver importance, symmetry, distance, band match, path support) with profile-specific adjustments
- **Veto Logic**: Filters candidates based on mountability, patch size, nearfield distance, first-reflection sidewalls in STUDIO
- **Diffuser Assignment**: Selects compatible diffuser from catalog based on target bands, depth limit, patch area
- **Layout Synthesis**: Grid-based panel placement within patch bounds
- **Profile Rules**: STUDIO penalizes nearfield and first reflections, HIFI allows more diffusion, PUBLIC with provisional warning
- **Decisions**: DIFFUSER, ABSORBER_RECOMMENDED, REJECT with explanations
- **Persistence**: PlacementCandidates linked to AnalysisRun, deterministic regeneration

### Technical Implementation
- Candidate extraction aggregates reflections by surface ID
- Patch creation clamps to minimum 0.6m size
- Scoring uses explicit weights: energy (0.25), time (0.15), receiver (0.10), mount (0.10), symmetry (0.10-0.20), distance (0.10), band match (0.10), path support (0.10)
- Studio profile doubles symmetry weight and nearfield penalty
- Veto checks: mountability, patch area (>0.25m2), depth limit, nearfield (<15ms), first-reflection sidewall in STUDIO
- Catalog compatibility check: depth limit, patch area, band match, category preference (QRD_2D > QRD_1D)
- Layout synthesis: grid fill with panel dimensions, centered in patch
- PUBLIC profile includes provisional-support warning

### Limitations
- **SHOEBOX Only**: Polyhedral geometry not supported
- **1st/2nd Order Only**: Uses PHASE 3 reflection data only
- **Grid Layout Only**: No advanced packing optimization
- **No Re-simulation**: Does not re-run acoustic analysis after treatment
- **PUBLIC Provisional**: PUBLIC profile uses simplified heuristics
- **Synchronous**: Placement generation runs synchronously in request flow

### Domain Models Added
- `DiffuserType` - Diffuser catalog entry with category, dimensions, effective range
- `PlacementCandidate` - Treatment recommendation with decision, score, reasons
- `BoundingBox` - Patch/zone bounds
- `CandidateReason` - Explanation code and message
- `PlacementSummary` - Aggregate statistics
- `LayoutResult` - Panel centers and coverage
- `Vec3` - 3D coordinate

### Database Migrations
- `006_phase4_diffusers.sql` - diffuser_types table with seed data, placement_candidates table

### Routes Added
- `POST /projects/{id}/placements/generate` - Generate treatment placements

### Files Added
- `internal/domain/placement.go` - Placement candidate domain model
- `internal/domain/diffuser.go` - Diffuser type domain model
- `internal/acoustics/diffuser/catalog.go` - Diffuser catalog and compatibility
- `internal/acoustics/diffuser/layout.go` - Layout synthesis
- `internal/acoustics/scoring/candidates.go` - Candidate extraction
- `internal/acoustics/scoring/ranker.go` - Scoring engine
- `internal/acoustics/scoring/veto.go` - Veto logic
- `internal/repo/diffuser_repo.go` - Diffuser repository
- `internal/repo/placement_repo.go` - Placement repository
- `internal/service/placement_service.go` - Placement service
- `internal/http/handlers/placements.go` - Placements HTTP handler
- `internal/acoustics/diffuser/catalog_test.go` - Catalog tests
- `internal/acoustics/diffuser/layout_test.go` - Layout tests
- `internal/acoustics/scoring/candidates_test.go` - Extraction tests
- `internal/acoustics/scoring/ranker_test.go` - Scoring tests
- `internal/acoustics/scoring/veto_test.go` - Veto tests

### Files Modified
- `internal/repo/repositories.go` - Added DiffuserRepository and PlacementRepository
- `internal/app/app.go` - Added PlacementService
- `internal/http/router.go` - Added placements handler and route
- `internal/http/handlers/projects.go` - Load placements data for placements tab
- `templates/partials/tab_placements.html` - Replaced placeholder with functional UI

---

## PHASE 5 - Job System + Results UI + Public Profile Polish ✅ COMPLETED

**Goal**: Complete UX for analysis execution and results.

### What Was Implemented
- ✅ In-process job tracking with persistence
- ✅ Job-backed analysis execution
- ✅ Job-backed placement generation
- ✅ HTMX polling for job status updates
- ✅ Job status panel with progress bar
- ✅ Results UI consolidation
- ✅ Stale-state detection for analysis and placements
- ✅ PUBLIC profile warnings in both analysis and placements
- ✅ Empty states for no analysis/placements
- ✅ Unit tests for job lifecycle

### Key Features
- **Job Model**: JobKind (ANALYSIS, PLACEMENTS), JobStatus (QUEUED, RUNNING, COMPLETED, FAILED), progress tracking
- **Job Phases**: Analysis (validating_inputs, loading_project_data, computing_modes, computing_rt, computing_reflections, persisting_results), Placements (validating_inputs, loading_analysis_context, extracting_candidates, scoring_candidates, applying_veto, assigning_diffusers, synthesizing_layouts, persisting_results)
- **Job Service**: CreateJob, StartJob, UpdateJobProgress, CompleteJob, FailJob, GetJob, GetLatestByProject, GetLatestActiveByProject, GetLatestByProjectAndKind
- **HTMX Polling**: Job status panel polls every 1 second while running, auto-refreshes on completion
- **Stale Detection**: Analysis stale if geometry/sources/receivers updated after last analysis, Placements stale if new analysis exists after last placement generation
- **PUBLIC Profile**: Explicit warnings in readiness checks for both analysis and placements
- **Empty States**: Clear messaging when no analysis run exists or no placements generated
- **Readiness Checks**: CanRunAnalysis checks geometry, sources, receivers, SHOEBOX type; CanGeneratePlacements checks completed analysis with reflections

### Technical Implementation
- Job persistence in SQLite with indexes on project_id, status, kind, created_at
- Goroutine-based execution for analysis and placement jobs
- Progress updates at key phases (10%, 20%, 30%, 50%, 60%, 70%, 80%, 90%, 100%)
- Stale check compares UpdatedAt of geometry/sources/receivers with AnalysisRun.CreatedAt
- Stale check compares AnalysisRun.CreatedAt with latest PlacementCandidate.CreatedAt
- Job status panel shows status badge, phase, progress bar, message, error message
- HTMX hx-trigger="every 1s" for polling, hx-swap="outerHTML" for refresh

### Limitations
- **In-Process Only**: No distributed workers or external queues
- **No WebSockets**: HTMX polling only
- **Simple Stale Detection**: Timestamp-based, not input hash-based
- **PUBLIC Profile**: Still provisional with simplified heuristics
- **No Re-simulation**: Does not re-run analysis after treatment
- **Goroutine Lifecycle**: Simple goroutines without advanced error recovery

### Domain Models Added
- Job - Job tracking with ID, ProjectID, Kind, Status, Phase, ProgressPercent, Message, ErrorMessage, CreatedAt, StartedAt, FinishedAt
- AnalysisReadiness - Ready, Reason, Warnings

### Database Migrations
- 007_phase5_jobs.sql - jobs table with indexes

### Routes Added
- GET /hx/projects/{id}/jobs/{jobID}/status - Job status polling

### Files Added
- internal/domain/job.go - Job domain model
- internal/repo/job_repo.go - Job repository
- internal/service/job_service.go - Job service with job execution methods
- internal/http/handlers/jobs.go - Jobs handler
- internal/service/job_service_test.go - Job lifecycle tests
- templates/partials/job_status.html - Job status panel template
- migrations/007_phase5_jobs.sql - Jobs table migration

### Files Modified
- internal/repo/repositories.go - Added JobRepository
- internal/app/app.go - Added JobService
- internal/http/router.go - Added jobs handler and job status route
- internal/http/handlers/analysis.go - Updated RunAnalysis to use job system
- internal/http/handlers/placements.go - Updated GeneratePlacements to use job system
- internal/http/handlers/projects.go - Load active jobs and stale state for tabs
- internal/service/analysis_service.go - Added CanRunAnalysis, IsAnalysisStale, GetLatestByProject
- internal/service/placement_service.go - Added ArePlacementsStale, updated PUBLIC warning
- templates/partials/tab_analysis.html - Job status panel, stale warning, empty state, readiness warnings
- templates/partials/tab_placements.html - Job status panel, stale warning, readiness warnings functionality

---

## PHASE 6 - Reports + Polish

**Goal**: Finalize MVP with export capabilities and polish.

### Planned Features
- ✅ JSON/CSV export
- ✅ HTML reports
- ✅ Sample projects
- ✅ Performance optimization
- ✅ Documentation completion

---

## Technical Architecture

### Layer Structure
```
cmd/server/          - Application entrypoint
internal/
  app/              - Configuration and bootstrap
  http/             - HTTP handlers and routing
  domain/           - Business entities and logic
  repo/             - Data persistence
  service/          - Business services
  acoustics/        - Acoustic calculation engines
  jobs/             - Background processing
templates/          - HTML templates
static/            - CSS and static assets
migrations/        - Database schema
```

### Key Principles
- **Incremental Development**: Each phase is buildable and functional
- **HTMX-First**: Server-side rendering with minimal JavaScript
- **Clean Architecture**: Clear separation of concerns
- **Test-Driven**: Unit tests for business logic
- **User-Focused**: Professional UX at every phase

---

## Current Status

**Phase 5 Complete** ✅

The application now provides a complete job tracking and results experience with:
- In-process job tracking with persistence
- Job-backed analysis execution with progress updates
- Job-backed placement generation with progress updates
- HTMX polling for real-time job status
- Job status panel with progress bar and phases
- Stale-state detection for analysis and placements
- PUBLIC profile warnings in analysis and placements
- Empty states for no analysis/placements
- Unit tests for job lifecycle

**Ready for**: Phase 6 - Reports + Polish

The job tracking system is complete and ready for export capabilities and final polish. Users can now:
- Run analysis as tracked jobs with progress updates
- Generate placements as tracked jobs with progress updates
- See real-time job status with progress bar
- View stale warnings when inputs change
- See clear PUBLIC profile warnings
- Regenerate jobs safely without confusing stale UI state
