<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let showGrid: boolean;
  export let showLabels: boolean;
  export let showAxes: boolean;

  const dispatch = createEventDispatcher();

  function toggleGrid() {
    dispatch('toggleGrid');
  }

  function toggleLabels() {
    dispatch('toggleLabels');
  }

  function toggleAxes() {
    dispatch('toggleAxes');
  }
</script>

<div class="viewer-controls">
  <h3>Viewer Options</h3>
  <div class="control-group">
    <label>
      <input type="checkbox" checked={showGrid} on:change={toggleGrid} />
      Show Grid
    </label>
    <label>
      <input type="checkbox" checked={showLabels} on:change={toggleLabels} />
      Show Labels
    </label>
    <label>
      <input type="checkbox" checked={showAxes} on:change={toggleAxes} />
      Show Axes
    </label>
  </div>
</div>

<style>
  .viewer-controls h3 {
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    color: #fff;
  }
  .control-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    color: #aaa;
    cursor: pointer;
  }
  input[type="checkbox"] {
    width: 16px;
    height: 16px;
    cursor: pointer;
  }
</style>
