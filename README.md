# Acoustics Calculator

Professional acoustic analysis and diffuser placement recommendations for studios, hi-fi rooms, and public spaces.

## Overview

This is a web application for acoustic analysis that provides:

- **Room modal analysis** for rectangular spaces
- **Early reflection calculations** using image-source method
- **Reverberation time estimation** (Sabine & Eyring)
- **Diffuser placement recommendations** with reasoning
- **Material absorption coefficients** database
- **Interactive 2D room editor** for source/receiver positioning

## Technology Stack

- **Backend**: Go 1.23+ with chi router
- **Database**: SQLite (modernc.org/sqlite, CGO-free)
- **Frontend**: HTML templates with HTMX
- **Styling**: Minimal CSS (no framework dependencies)
- **JavaScript**: Minimal, only for 2D editor interactions

## Quick Start

### Prerequisites

- Go 1.23 or later
- Make (for build commands)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd acoustics-calculator
```

2. Set up dependencies:
```bash
make setup
```

3. Run the application:
```bash
make run
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Development

### Build Commands

- `make run` - Start the development server
- `make build` - Build the application binary
- `make test` - Run all tests
- `make fmt` - Format Go code
- `make clean` - Remove build artifacts
- `make reset-db` - Reset the database

### Project Structure

```
cmd/server/          - Application entrypoint
internal/
  app/              - Application configuration and bootstrap
  http/             - HTTP handlers and routing
  domain/           - Domain models and business logic
  repo/             - Database repositories and migrations
  service/          - Business service layer
  acoustics/        - Acoustic calculation engines
  jobs/             - Background job processing
templates/          - HTML templates
static/            - CSS and static assets
migrations/        - Database migration files
```

## Features

### Space Types Supported

1. **Studio / Control Room**
   - Optimized for mixing and mastering
   - Symmetry-focused placement
   - Near-field diffusion avoidance

2. **Hi-Fi Listening Room**
   - Critical listening environments
   - Rear wall and rear-side optimization
   - Balanced treatment approach

3. **Public Spaces**
   - Larger venue analysis
   - Audience zone considerations
   - Coverage and clarity metrics

### Analysis Pipeline

1. **Geometry Input** - Room dimensions and surface materials
2. **Source/Receiver Setup** - Speaker and listener positioning
3. **Modal Analysis** - Room mode calculations and pressure mapping
4. **Reflection Analysis** - Early reflection paths and timing
5. **RT Estimation** - Reverberation time by frequency bands
6. **Placement Generation** - Diffuser/absorber candidate zones
7. **Scoring & Ranking** - Intelligent placement recommendations

## MVP Scope

The current MVP (Phase 6 Complete - MVP Done) implements:

- ✅ Project management with space type configuration
- ✅ Shoebox room geometry with automatic surface generation
- ✅ Material library with preset absorption/scattering coefficients
- ✅ Surface material assignment interface
- ✅ Interactive 2D room editor (top-down, drag-and-drop)
- ✅ Source positioning with CRUD operations
- ✅ Receiver positioning with CRUD operations
- ✅ Treatment constraint configuration
- ✅ Real-time validation and error handling
- ✅ Summary tab with configuration readiness tracking
- ✅ Rectangular modal frequency calculations
- ✅ First and second order early reflections
- ✅ Sabine and Eyring RT calculations
- ✅ Acoustic analysis execution with persistence
- ✅ Results visualization (modes, RT, reflections)
- ✅ Prerequisites validation and warnings
- ✅ Unit tests for acoustics engine
- ✅ Diffuser placement candidate generation
- ✅ Scoring and ranking system
- ✅ Treatment recommendations (DIFFUSER, ABSORBER_RECOMMENDED, REJECT)
- ✅ Diffuser type assignment (QRD_1D, QRD_2D, CUSTOM)
- ✅ Layout synthesis for panel coordinates
- ✅ Placement persistence linked to AnalysisRun
- ✅ Placements UI with results
- ✅ Profile-specific rules (STUDIO, HIFI, PUBLIC)
- ✅ Unit tests for scoring, veto, catalog, layout
- ✅ In-process job tracking with persistence
- ✅ Job-backed analysis execution
- ✅ Job-backed placement generation
- ✅ HTMX polling for job status updates
- ✅ Job status panel with progress bar
- ✅ Stale-state detection for analysis and placements
- ✅ PUBLIC profile warnings in analysis and placements
- ✅ Empty states for no analysis/placements
- ✅ Unit tests for job lifecycle
- ✅ JSON export for full project report
- ✅ CSV export for placement recommendations
- ✅ Printable HTML report page
- ✅ Export actions in project workspace
- ✅ Demo project creation
- ✅ Empty/stale/PUBLIC state handling in reports
- ✅ Print-friendly CSS
- ✅ Unit tests for report/export

**MVP Status**: COMPLETE

## Future Enhancements

- Polyhedral room geometry support
- Advanced ray tracing analysis
- FEM/BEM acoustic modeling
- Real-time analysis updates
- Export to CAD formats
- Mobile responsive interface

## Contributing

1. Follow the existing code style and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure all builds and tests pass before submitting

## License

[License information to be added]
