# Parametric Room Viewer

A simple parametric room visualization application with Go backend and Svelte frontend.

## What This Is

A single-screen application for visualizing room configurations in 3D:
- **Parameters are the source of truth** - edit values on the left, see results on the right
- **3D viewer is read-only** - no CAD editing, no dragging in the scene
- **Simple workflow** - one screen, no tabs, no database, no auth

## Stack

**Backend:**
- Go 1.21+
- chi router
- JSON API only

**Frontend:**
- Svelte
- Vite
- TypeScript
- Three.js

## Project Structure

```
backend/
  cmd/server/main.go
  go.mod
  internal/
    api/
      dto.go
      handlers.go
      router.go
      handlers_test.go
    materials/
      catalog.go
    scene/
      model.go
      derive.go
      example.go
      example_test.go
    validation/
      validate.go
      validate_test.go

frontend/
  index.html
  package.json
  vite.config.ts
  tsconfig.json
  tsconfig.node.json
  svelte.config.js
  src/
    main.ts
    App.svelte
    lib/
      api/
        client.ts
      scene/
        types.ts
      viewer/
        room-viewer.ts
      components/
        ControlPanel.svelte
        RoomControls.svelte
        MaterialControls.svelte
        SourceControls.svelte
        ReceiverControls.svelte
        ViewerControls.svelte
        ErrorPanel.svelte
        Viewer3D.svelte
```

## Running the Backend

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

## Running the Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend dev server runs on `http://localhost:5173` (proxies API to backend)

## Building for Production

**Backend:**
```bash
cd backend
go build -o server cmd/server/main.go
./server
```

**Frontend:**
```bash
cd frontend
npm run build
npm run preview
```

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/scene/example` - Get default example scene configuration
- `POST /api/scene/preview` - Validate and derive render-ready scene from configuration

## Supported Parameters

**Room:**
- Width, length, height (meters)

**Surface Materials:**
- Front wall, rear wall, left wall, right wall, ceiling, floor
- Materials: gypsum_board, painted_concrete, wood_panel, carpet, glass, curtain_heavy

**Sources:**
- Left speaker (x, y, z)
- Right speaker (x, y, z)

**Receiver:**
- Listening position (x, y, z)

**Viewer Options:**
- Show grid
- Show labels
- Show axes

## Preview Flow

1. On load, frontend fetches example scene from `/api/scene/example`
2. Frontend populates form with example values
3. Frontend POSTs to `/api/scene/preview` to get derived scene
4. Backend validates configuration and derives render-ready data
5. Frontend renders 3D scene with Three.js
6. On parameter change, frontend debounces and re-requests preview

## Important Notes

- **This is not a CAD application** - the 3D viewer is for visualization only
- **Parameters are the source of truth** - edit the form, not the scene
- **No database** - all state is transient
- **No authentication** - open API
- **No acoustic analysis** - this is a viewer only
