# Acoustics Calculator - Technical Assumptions

## Overview

This document captures the technical assumptions and design decisions made during development of the Acoustics Calculator project.

---

## Database Assumptions

### SQLite Choice
- **Assumption**: SQLite with modernc.org/sqlite driver provides sufficient performance and reliability
- **Rationale**: CGO-free implementation simplifies deployment and cross-platform compatibility
- **Trade-offs**: Limited concurrent write performance, but acceptable for single-user application

### Schema Design
- **Assumption**: JSON storage for complex data structures (polygons, band maps) is sufficient for MVP
- **Rationale**: Simplifies schema evolution and provides flexibility
- **Trade-offs**: Limited query capabilities on JSON fields, but acceptable for current use case

### Migration Strategy
- **Assumption**: Simple file-based migration system with version tracking is adequate
- **Rationale**: Keeps deployment simple and provides clear upgrade path
- **Trade-offs**: No rollback mechanism, but migrations are designed to be forward-only

---

## Architecture Assumptions

### Go + HTMX Stack
- **Assumption**: Server-side rendering with HTMX provides adequate UX without SPA complexity
- **Rationale**: Reduces JavaScript complexity while maintaining interactivity
- **Trade-offs**: Limited offline capability, but acceptable for web-based tool

### Package Organization
- **Assumption**: Standard Go project structure with clear separation of concerns
- **Rationale**: Follows Go best practices and maintains maintainability
- **Trade-offs**: Slightly more boilerplate, but provides clear boundaries

### Dependency Injection
- **Assumption**: Simple constructor-based DI through App struct is sufficient
- **Rationale**: Avoids complex DI frameworks while maintaining testability
- **Trade-offs**: Manual wiring, but keeps dependencies explicit

---

## Acoustic Engineering Assumptions

### Rectangular Room Model
- **Assumption**: Shoebox (rectangular) room model covers majority of use cases for MVP
- **Rationale**: Simplifies calculations while providing valuable insights
- **Trade-offs**: Limited to rectangular spaces, polyhedral support planned for future phases

### Frequency Bands
- **Assumption**: Standard octave bands (63, 125, 250, 500, 1000, 2000, 4000 Hz) are sufficient
- **Rationale**: Industry standard for acoustic analysis
- **Trade-offs**: Limited frequency resolution, but adequate for most applications

### Speed of Sound
- **Assumption**: Default speed of sound (343 m/s at 20°C) with temperature adjustment is adequate
- **Rationale**: Provides reasonable accuracy for most room acoustics applications
- **Trade-offs**: Doesn't account for humidity effects in detail, but impact is minimal

### Material Properties
- **Assumption**: Pre-defined material library with absorption coefficients covers common cases
- **Rationale**: Reduces user complexity while providing realistic results
- **Trade-offs**: Limited customization, but sufficient for MVP

### 2D Room Editor
- **Assumption**: Top-down 2D SVG editor with drag-and-drop provides adequate positioning UX
- **Rationale**: Simplifies implementation while maintaining visual feedback
- **Trade-offs**: Z-coordinate edited via forms only, no 3D visualization

### Coordinate System
- **Assumption**: Origin at front-left-floor corner (X=width, Y=length, Z=height) is intuitive
- **Rationale**: Standard coordinate system for room acoustics
- **Trade-offs**: May require user education, but consistent with industry practice

### Geometry Validation
- **Assumption**: Backend validation for positions within room bounds prevents invalid states
- **Rationale**: Ensures data integrity for acoustic calculations
- **Trade-offs**: Requires server round-trips for validation, but provides robustness

### Acoustic Calculations
- **Assumption**: Rectangular modal formula provides adequate mode prediction for shoebox rooms
- **Rationale**: Standard analytical solution for rectangular spaces
- **Trade-offs**: Limited to shoebox geometry, polyhedral not supported

- **Assumption**: Image-source method up to 2nd order captures most significant early reflections
- **Rationale**: Higher-order reflections have diminishing energy contribution
- **Trade-offs**: Misses some late-arriving reflections, but adequate for early reflection analysis

- **Assumption**: Sabine and Eyring formulas provide reasonable RT60 estimates
- **Rationale**: Industry-standard formulas with known limitations
- **Trade-offs**: Assumes diffuse field, may not be accurate for small rooms

- **Assumption**: Default absorption coefficients (0.10-0.40) provide fallback for missing materials
- **Rationale**: Prevents calculation failures when material data incomplete
- **Trade-offs**: Less accurate than actual material data, but allows analysis to proceed

- **Assumption**: Speed of sound formula 331.3 + 0.606*tempC is adequate for temperature adjustment
- **Rationale**: Simple linear approximation for typical temperature ranges
- **Trade-offs**: Humidity effects ignored, impact minimal for most applications

### Placement Generation
- **Assumption**: Reflection aggregation by surface identifies high-interest treatment zones
- **Rationale**: Surfaces with more reflections are likely acoustic problem areas
- **Trade-offs**: Simple aggregation may miss nuanced acoustic patterns

- **Assumption**: Weighted scoring components provide adequate candidate ranking
- **Rationale**: Explicit weights make behavior transparent and tunable
- **Trade-offs**: Weight selection is heuristic, may not capture all acoustic factors

- **Assumption**: Grid-based layout synthesis is sufficient for V1
- **Rationale**: Simple deterministic placement without complex optimization
- **Trade-offs**: May not achieve optimal coverage or aesthetic arrangement

- **Assumption**: Veto rules prevent inappropriate placements
- **Rationale**: Nearfield and first-reflection zones are better treated with absorption in STUDIO
- **Trade-offs**: Conservative approach may miss viable diffusion opportunities

- **Assumption**: Diffuser catalog compatibility check ensures feasible recommendations
- **Rationale**: Depth, area, and band constraints prevent impossible placements
- **Trade-offs**: Limited catalog may restrict recommendation quality

- **Assumption**: PUBLIC profile heuristics provide provisional guidance
- **Rationale**: Public-space acoustics require different approach, but simplified for V1
- **Trade-offs**: Recommendations may not be optimal for complex public spaces

---

## User Interface Assumptions

### Browser Compatibility
- **Assumption**: Modern browsers with ES6+ support provide adequate coverage
- **Rationale**: HTMX and modern CSS features work reliably in current browsers
- **Trade-offs**: Excludes very old browsers, but acceptable for professional tool

### Screen Size
- **Assumption**: Desktop-first design with mobile responsiveness is adequate
- **Rationale**: Most acoustic work done on larger screens
- **Trade-offs**: Mobile experience limited, but functional

### Form Validation
- **Assumption**: Server-side validation with client-side feedback provides good UX
- **Rationale**: Ensures data integrity while providing responsive interface
- **Trade-offs**: Slightly more server round-trips, but provides robust validation

---

## Performance Assumptions

### Single-User Application
- **Assumption**: Application designed for single user or small team usage
- **Rationale**: Simplifies architecture and reduces complexity
- **Trade-offs**: Not optimized for high concurrency, but not required for use case

### Calculation Complexity
- **Assumption**: Acoustic calculations complete in reasonable time for room-sized problems
- **Rationale**: Image-source and modal analysis are computationally feasible for typical rooms
- **Trade-offs**: May not scale to very large venues, but adequate for target use cases

### Database Performance
- **Assumption**: SQLite provides adequate performance for project data storage
- **Rationale**: Project data size is relatively small and read-heavy
- **Trade-offs**: Limited write performance, but acceptable for application patterns

---

## Security Assumptions

### Authentication
- **Assumption**: No authentication required for MVP (local application)
- **Rationale**: Simplifies development and deployment for single-user tool
- **Trade-offs**: No multi-user support, but acceptable for initial release

### Data Privacy
- **Assumption**: Local data storage provides adequate privacy
- **Rationale**: Sensitive acoustic data stays on user's machine
- **Trade-offs**: No cloud sync or collaboration features

### Input Validation
- **Assumption**: Comprehensive input validation prevents security issues
- **Rationale**: Server-side validation ensures data integrity
- **Trade-offs**: Slightly more complex validation logic, but provides security

---

## Development Assumptions

### Go Version
- **Assumption**: Go 1.23+ features are available and stable
- **Rationale**: Modern Go features improve development experience
- **Trade-offs**: Requires relatively recent Go version, but reasonable for new project

### Testing Strategy
- **Assumption**: Unit tests for business logic provide adequate coverage
- **Rationale**: Focus testing on critical acoustic calculations and business rules
- **Trade-offs**: Limited integration testing, but sufficient for MVP

### Documentation
- **Assumption**: Inline documentation and phase docs provide adequate guidance
- **Rationale**: Balances documentation effort with development velocity
- **Trade-offs**: Less comprehensive than enterprise documentation, but appropriate for project size

---

## Future Considerations

### Scalability
- **Current assumption**: Single-user, single-machine deployment
- **Future consideration**: Multi-user support with proper authentication
- **Migration path**: Add authentication layer, move to PostgreSQL if needed

### Advanced Acoustics
- **Current assumption**: Rectangular rooms with basic materials
- **Future consideration**: Polyhedral geometry, advanced materials, ray tracing
- **Migration path**: Extend geometry engine, enhance material models

### Cloud Features
- **Current assumption**: Local-only application
- **Future consideration**: Cloud sync, collaboration, shared libraries
- **Migration path**: Add cloud storage backend, implement sync logic

---

## Assumption Review Process

These assumptions should be reviewed:
1. **At each phase transition** - Validate assumptions still hold
2. **When requirements change** - Reassess technical decisions
3. **During user testing** - Validate UX assumptions
4. **Performance testing** - Verify performance assumptions
5. **Security review** - Confirm security assumptions remain valid

---

## Notes

This document is living and should be updated as assumptions change or new ones are identified during development.
