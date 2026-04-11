<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { DerivedScene } from '../api/client';
  import { RoomViewer } from '../viewer/room-viewer';

  export let derivedScene: DerivedScene | null = null;
  export let showGrid: boolean = true;
  export let showLabels: boolean = true;
  export let showAxes: boolean = true;

  let container: HTMLElement;
  let viewer: RoomViewer | null = null;
  let isFirstLoad = true;

  function resetCamera() {
    if (viewer) {
      viewer.resetCamera();
    }
  }

  function setPresetView(preset: 'isometric' | 'front' | 'left' | 'top') {
    if (viewer) {
      viewer.setPresetView(preset);
    }
  }

  onMount(() => {
    viewer = new RoomViewer(container);
    viewer.animate();

    if (derivedScene) {
      viewer.updateScene(derivedScene, showGrid, showLabels, showAxes, true);
      isFirstLoad = false;
    }
  });

  onDestroy(() => {
    if (viewer) {
      viewer.dispose();
    }
  });

  $: if (viewer && derivedScene) {
    viewer.updateScene(derivedScene, showGrid, showLabels, showAxes, isFirstLoad);
    if (isFirstLoad) isFirstLoad = false;
  }
</script>

<div class="viewer-wrapper">
  <div class="viewer-container" bind:this={container}></div>
  <div class="camera-toolbar">
    <div class="toolbar-section">
      <span class="toolbar-label">View</span>
      <button class="toolbar-btn" on:click={() => setPresetView('isometric')}>Iso</button>
      <button class="toolbar-btn" on:click={() => setPresetView('front')}>Front</button>
      <button class="toolbar-btn" on:click={() => setPresetView('left')}>Left</button>
      <button class="toolbar-btn" on:click={() => setPresetView('top')}>Top</button>
    </div>
    <div class="toolbar-section">
      <button class="toolbar-btn reset-btn" on:click={resetCamera}>Reset</button>
    </div>
  </div>
  <div class="viewer-help">
    <div class="help-title">3D Controls</div>
    <div class="help-item">Rotate: Left drag</div>
    <div class="help-item">Pan: Right drag</div>
    <div class="help-item">Zoom: Mouse wheel</div>
  </div>
</div>

<style>
  .viewer-wrapper {
    position: relative;
    width: 100%;
    height: 100%;
  }

  .viewer-container {
    width: 100%;
    height: 100%;
  }

  .camera-toolbar {
    position: absolute;
    top: 16px;
    right: 16px;
    display: flex;
    gap: 8px;
    background: rgba(0, 0, 0, 0.8);
    border-radius: 8px;
    padding: 8px 12px;
    pointer-events: auto;
    backdrop-filter: blur(4px);
    border: 1px solid #444;
  }

  .toolbar-section {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .toolbar-section:not(:last-child) {
    padding-right: 8px;
    border-right: 1px solid #444;
  }

  .toolbar-label {
    font-size: 0.75rem;
    color: #888;
    margin-right: 4px;
    font-weight: 500;
  }

  .toolbar-btn {
    padding: 4px 8px;
    background: #333;
    border: 1px solid #444;
    border-radius: 4px;
    color: #fff;
    font-size: 0.75rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .toolbar-btn:hover {
    background: #444;
    border-color: #555;
  }

  .reset-btn {
    background: #2196F3;
    border-color: #2196F3;
  }

  .reset-btn:hover {
    background: #1976D2;
    border-color: #1976D2;
  }

  .viewer-help {
    position: absolute;
    bottom: 16px;
    right: 16px;
    background: rgba(0, 0, 0, 0.8);
    border-radius: 8px;
    padding: 12px 16px;
    font-size: 0.8rem;
    color: #fff;
    pointer-events: auto;
    backdrop-filter: blur(4px);
    border: 1px solid #444;
  }

  .help-title {
    font-weight: 600;
    margin-bottom: 8px;
    color: #2196F3;
  }

  .help-item {
    margin: 4px 0;
    color: #ccc;
  }
</style>
