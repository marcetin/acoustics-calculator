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

  function resetCamera() {
    if (viewer) {
      viewer.resetCamera();
    }
  }

  onMount(() => {
    viewer = new RoomViewer(container);
    viewer.animate();

    if (derivedScene) {
      viewer.updateScene(derivedScene, showGrid, showLabels, showAxes);
    }
  });

  onDestroy(() => {
    if (viewer) {
      viewer.dispose();
    }
  });

  $: if (viewer && derivedScene) {
    viewer.updateScene(derivedScene, showGrid, showLabels, showAxes);
  }
</script>

<div class="viewer-wrapper">
  <div class="viewer-container" bind:this={container}></div>
  <div class="viewer-help">
    <div class="help-title">3D Controls</div>
    <div class="help-item">Rotate: Left drag</div>
    <div class="help-item">Pan: Right drag</div>
    <div class="help-item">Zoom: Mouse wheel</div>
    <button class="reset-btn" on:click={resetCamera}>Reset View</button>
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

  .reset-btn {
    margin-top: 8px;
    padding: 6px 12px;
    background: #2196F3;
    border: none;
    border-radius: 4px;
    color: white;
    font-size: 0.75rem;
    cursor: pointer;
    transition: background 0.2s;
  }

  .reset-btn:hover {
    background: #1976D2;
  }
</style>
