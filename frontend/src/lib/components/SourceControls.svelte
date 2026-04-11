<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let leftX: number;
  export let leftY: number;
  export let leftZ: number;
  export let rightX: number;
  export let rightY: number;
  export let rightZ: number;

  const dispatch = createEventDispatcher();

  function updateLeftX(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateLeftX', { value: num });
  }

  function updateLeftY(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateLeftY', { value: num });
  }

  function updateLeftZ(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateLeftZ', { value: num });
  }

  function updateRightX(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateRightX', { value: num });
  }

  function updateRightY(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateRightY', { value: num });
  }

  function updateRightZ(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateRightZ', { value: num });
  }
</script>

<div class="source-controls">
  <h3>Speaker Positions</h3>
  <div class="control-group">
    <h4>Left Speaker</h4>
    <label>
      X (m)
      <input type="number" step="0.1" min="0" value={leftX} on:input={(e) => updateLeftX(e.currentTarget.value)} />
    </label>
    <label>
      Y (m)
      <input type="number" step="0.1" min="0" value={leftY} on:input={(e) => updateLeftY(e.currentTarget.value)} />
    </label>
    <label>
      Z (m)
      <input type="number" step="0.1" min="0" value={leftZ} on:input={(e) => updateLeftZ(e.currentTarget.value)} />
    </label>
  </div>
  <div class="control-group">
    <h4>Right Speaker</h4>
    <label>
      X (m)
      <input type="number" step="0.1" min="0" value={rightX} on:input={(e) => updateRightX(e.currentTarget.value)} />
    </label>
    <label>
      Y (m)
      <input type="number" step="0.1" min="0" value={rightY} on:input={(e) => updateRightY(e.currentTarget.value)} />
    </label>
    <label>
      Z (m)
      <input type="number" step="0.1" min="0" value={rightZ} on:input={(e) => updateRightZ(e.currentTarget.value)} />
    </label>
  </div>
</div>

<style>
  .source-controls h3 {
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    color: #fff;
  }
  .source-controls h4 {
    font-size: 0.8rem;
    margin: 0.5rem 0 0.25rem;
    color: #bbb;
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
  input {
    padding: 0.5rem;
    background: #333;
    border: 1px solid #444;
    border-radius: 4px;
    color: #fff;
    font-size: 0.9rem;
  }
  input:focus {
    outline: none;
    border-color: #2196F3;
  }
</style>
