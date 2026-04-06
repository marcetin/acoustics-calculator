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

## PHASE 2 - Geometry + Materials + Inputs (NEXT)

**Goal**: Implement complete input system for acoustic analysis.

### Planned Features
- ✅ Room geometry configuration (rectangular/shoebox)
- ✅ Surface material assignments
- ✅ 2D room editor with drag-and-drop
- ✅ Source positioning system
- ✅ Receiver positioning system
- ✅ Constraint configuration
- ✅ Material library with absorption coefficients
- ✅ Geometry validation

### Technical Implementation
- Interactive SVG-based room editor
- Drag-and-drop source/receiver positioning
- Real-time validation
- Material database with frequency-dependent properties
- Export/import functionality for room layouts

---

## PHASE 3 - Acoustics Engine V1

**Goal**: Implement core acoustic analysis capabilities.

### Planned Features
- ✅ Rectangular modal solver
- ✅ Image-source early reflections (1st/2nd order)
- ✅ RT60 calculation (Sabine & Eyring)
- ✅ Frequency band analysis
- ✅ Basic metrics calculation
- ✅ Background job system
- ✅ Analysis status tracking

### Analysis Pipeline
1. Input validation
2. Modal frequency calculation
3. Early reflection path tracing
4. Reverberation time estimation
5. Metric computation
6. Result storage and presentation

---

## PHASE 4 - Diffuser Candidate Engine V1

**Goal**: Generate intelligent treatment recommendations.

### Planned Features
- ✅ Candidate zone generation
- ✅ Scoring and ranking system
- ✅ Diffuser vs absorber decision logic
- ✅ Layout synthesis
- ✅ Placement reasoning
- ✅ Results visualization

### Decision Logic
- Space-type specific rules
- Acoustic problem identification
- Treatment type selection
- Position optimization
- Constraint satisfaction

---

## PHASE 5 - Job System + Results UI + Public Profile

**Goal**: Complete UX for analysis execution and results.

### Planned Features
- ✅ Background job processing
- ✅ Real-time status updates
- ✅ Results dashboard
- ✅ Export functionality
- ✅ Public space heuristics
- ✅ Report generation

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

**Phase 1 Complete** ✅

The application now provides a fully functional project management system with:
- Complete CRUD operations for projects
- Professional web interface with HTMX navigation
- Comprehensive validation and error handling
- Extensible foundation for future acoustic features

**Next**: Phase 2 - Geometry + Materials + Inputs

Ready to begin implementing the room configuration and input systems that will enable actual acoustic analysis.
