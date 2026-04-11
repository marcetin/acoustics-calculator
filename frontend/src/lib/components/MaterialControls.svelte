<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let materials: string[];
  export let frontWall: string;
  export let rearWall: string;
  export let leftWall: string;
  export let rightWall: string;
  export let ceiling: string;
  export let floor: string;

  const dispatch = createEventDispatcher();

  function updateFrontWall(value: string) {
    dispatch('updateFrontWall', { value });
  }

  function updateRearWall(value: string) {
    dispatch('updateRearWall', { value });
  }

  function updateLeftWall(value: string) {
    dispatch('updateLeftWall', { value });
  }

  function updateRightWall(value: string) {
    dispatch('updateRightWall', { value });
  }

  function updateCeiling(value: string) {
    dispatch('updateCeiling', { value });
  }

  function updateFloor(value: string) {
    dispatch('updateFloor', { value });
  }
</script>

<div class="material-controls">
  <h3>Surface Materials</h3>
  <div class="control-group">
    <label>
      Front Wall
      <select value={frontWall} on:change={(e) => updateFrontWall(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
    <label>
      Rear Wall
      <select value={rearWall} on:change={(e) => updateRearWall(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
    <label>
      Left Wall
      <select value={leftWall} on:change={(e) => updateLeftWall(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
    <label>
      Right Wall
      <select value={rightWall} on:change={(e) => updateRightWall(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
    <label>
      Ceiling
      <select value={ceiling} on:change={(e) => updateCeiling(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
    <label>
      Floor
      <select value={floor} on:change={(e) => updateFloor(e.currentTarget.value)}>
        {#each materials as material}
          <option value={material}>{material}</option>
        {/each}
      </select>
    </label>
  </div>
</div>

<style>
  .material-controls h3 {
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
    flex-direction: column;
    gap: 0.25rem;
    font-size: 0.85rem;
    color: #aaa;
  }
  select {
    padding: 0.5rem;
    background: #333;
    border: 1px solid #444;
    border-radius: 4px;
    color: #fff;
    font-size: 0.9rem;
  }
  select:focus {
    outline: none;
    border-color: #2196F3;
  }
</style>
