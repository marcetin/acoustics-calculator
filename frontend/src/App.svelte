<script lang="ts">
  import { onMount } from 'svelte';
  import ControlPanel from './lib/components/ControlPanel.svelte';
  import RoomControls from './lib/components/RoomControls.svelte';
  import MaterialControls from './lib/components/MaterialControls.svelte';
  import SourceControls from './lib/components/SourceControls.svelte';
  import ReceiverControls from './lib/components/ReceiverControls.svelte';
  import ViewerControls from './lib/components/ViewerControls.svelte';
  import ErrorPanel from './lib/components/ErrorPanel.svelte';
  import Viewer3D from './lib/components/Viewer3D.svelte';
  import { getExample, previewScene, type SceneConfig, type DerivedScene } from './lib/api/client';

  let config: SceneConfig | null = null;
  let derivedScene: DerivedScene | null = null;
  let lastValidScene: DerivedScene | null = null;
  let error: string | null = null;
  let fieldErrors: Record<string, string> = {};
  let materials = ['gypsum_board', 'painted_concrete', 'wood_panel', 'carpet', 'glass', 'curtain_heavy'];
  let debounceTimer: number | null = null;

  onMount(async () => {
    try {
      config = await getExample();
      await requestPreview();
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to load example scene';
    }
  });

  async function requestPreview() {
    if (!config) return;

    try {
      const response = await previewScene(config);
      derivedScene = response.derived;
      lastValidScene = response.derived;
      error = null;
      fieldErrors = {};
    } catch (e) {
      // Preserve last valid scene
      if (lastValidScene) {
        derivedScene = lastValidScene;
      }
      error = e instanceof Error ? e.message : 'Failed to preview scene';
      
      // Extract field errors from error object if available
      if (e instanceof Error && (e as any).fields) {
        fieldErrors = (e as any).fields;
      }
    }
  }

  function debouncedPreview() {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    debounceTimer = window.setTimeout(() => {
      requestPreview();
    }, 300);
  }

  function updateRoomDimension(field: 'widthM' | 'lengthM' | 'heightM', value: number) {
    if (!config) return;
    config = {
      ...config,
      room: {
        ...config.room,
        [field]: value
      }
    };
    debouncedPreview();
  }

  function updateSurfaceMaterial(field: string, value: string) {
    if (!config) return;
    config = {
      ...config,
      surfaces: {
        ...config.surfaces,
        [field]: { materialId: value }
      }
    };
    debouncedPreview();
  }

  function updateSourcePosition(source: 'left' | 'right', axis: 'x' | 'y' | 'z', value: number) {
    if (!config) return;
    config = {
      ...config,
      sources: {
        ...config.sources,
        [source]: {
          ...config.sources[source],
          [axis]: value
        }
      }
    };
    debouncedPreview();
  }

  function updateReceiverPosition(axis: 'x' | 'y' | 'z', value: number) {
    if (!config) return;
    config = {
      ...config,
      receiver: {
        ...config.receiver,
        [axis]: value
      }
    };
    debouncedPreview();
  }

  function updateViewerOption(field: 'showGrid' | 'showLabels' | 'showAxes', value: boolean) {
    if (!config) return;
    config = {
      ...config,
      viewer: {
        ...config.viewer,
        [field]: value
      }
    };
    debouncedPreview();
  }

  function resetToExample() {
    getExample().then((example: SceneConfig) => {
      config = example;
      fieldErrors = {};
      error = null;
      requestPreview();
    });
  }
</script>

<div class="app">
  <div class="sidebar">
    <ControlPanel>
      <div class="header">
        <h2>Room Viewer</h2>
        <button on:click={resetToExample}>Reset to Example</button>
      </div>

      {#if config}
        <RoomControls
          width={config.room.widthM}
          length={config.room.lengthM}
          height={config.room.heightM}
          {fieldErrors}
          on:updateWidth={(e) => updateRoomDimension('widthM', e.detail.value)}
          on:updateLength={(e) => updateRoomDimension('lengthM', e.detail.value)}
          on:updateHeight={(e) => updateRoomDimension('heightM', e.detail.value)}
        />

        <MaterialControls
          {materials}
          frontWall={config.surfaces.frontWall.materialId}
          rearWall={config.surfaces.rearWall.materialId}
          leftWall={config.surfaces.leftWall.materialId}
          rightWall={config.surfaces.rightWall.materialId}
          ceiling={config.surfaces.ceiling.materialId}
          floor={config.surfaces.floor.materialId}
          on:updateFrontWall={(e) => updateSurfaceMaterial('frontWall', e.detail.value)}
          on:updateRearWall={(e) => updateSurfaceMaterial('rearWall', e.detail.value)}
          on:updateLeftWall={(e) => updateSurfaceMaterial('leftWall', e.detail.value)}
          on:updateRightWall={(e) => updateSurfaceMaterial('rightWall', e.detail.value)}
          on:updateCeiling={(e) => updateSurfaceMaterial('ceiling', e.detail.value)}
          on:updateFloor={(e) => updateSurfaceMaterial('floor', e.detail.value)}
        />

        <SourceControls
          leftX={config.sources.left.x}
          leftY={config.sources.left.y}
          leftZ={config.sources.left.z}
          rightX={config.sources.right.x}
          rightY={config.sources.right.y}
          rightZ={config.sources.right.z}
          {fieldErrors}
          on:updateLeftX={(e) => updateSourcePosition('left', 'x', e.detail.value)}
          on:updateLeftY={(e) => updateSourcePosition('left', 'y', e.detail.value)}
          on:updateLeftZ={(e) => updateSourcePosition('left', 'z', e.detail.value)}
          on:updateRightX={(e) => updateSourcePosition('right', 'x', e.detail.value)}
          on:updateRightY={(e) => updateSourcePosition('right', 'y', e.detail.value)}
          on:updateRightZ={(e) => updateSourcePosition('right', 'z', e.detail.value)}
        />

        <ReceiverControls
          x={config.receiver.x}
          y={config.receiver.y}
          z={config.receiver.z}
          {fieldErrors}
          on:updateX={(e) => updateReceiverPosition('x', e.detail.value)}
          on:updateY={(e) => updateReceiverPosition('y', e.detail.value)}
          on:updateZ={(e) => updateReceiverPosition('z', e.detail.value)}
        />

        <ViewerControls
          showGrid={config.viewer.showGrid}
          showLabels={config.viewer.showLabels}
          showAxes={config.viewer.showAxes}
          on:toggleGrid={() => updateViewerOption('showGrid', !config.viewer.showGrid)}
          on:toggleLabels={() => updateViewerOption('showLabels', !config.viewer.showLabels)}
          on:toggleAxes={() => updateViewerOption('showAxes', !config.viewer.showAxes)}
        />
      {/if}

      <ErrorPanel {error} />
    </ControlPanel>
  </div>

  <div class="viewer">
    {#if derivedScene}
      <Viewer3D
        {derivedScene}
        showGrid={config?.viewer.showGrid ?? true}
        showLabels={config?.viewer.showLabels ?? true}
        showAxes={config?.viewer.showAxes ?? true}
      />
    {:else}
      <div class="loading">Loading 3D viewer...</div>
    {/if}
  </div>
</div>

<style>
  .app {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .sidebar {
    width: 350px;
    min-width: 350px;
    overflow-y: auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .header h2 {
    font-size: 1.1rem;
    color: #fff;
    margin: 0;
  }

  button {
    padding: 0.5rem 1rem;
    background: #2196F3;
    border: none;
    border-radius: 4px;
    color: white;
    font-size: 0.85rem;
    cursor: pointer;
    transition: background 0.2s;
  }

  button:hover {
    background: #1976D2;
  }

  .viewer {
    flex: 1;
    position: relative;
  }

  .loading {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #888;
    font-size: 1rem;
  }
</style>
