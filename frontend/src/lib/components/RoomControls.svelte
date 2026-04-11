<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let width: number;
  export let length: number;
  export let height: number;
  export let fieldErrors: Record<string, string> = {};

  const dispatch = createEventDispatcher();

  function updateWidth(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateWidth', { value: num });
  }

  function updateLength(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateLength', { value: num });
  }

  function updateHeight(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateHeight', { value: num });
  }

  $: widthError = fieldErrors['room.widthM'] || '';
  $: lengthError = fieldErrors['room.lengthM'] || '';
  $: heightError = fieldErrors['room.heightM'] || '';
</script>

<div class="room-controls">
  <h3>Room Dimensions</h3>
  <div class="control-group">
    <label>
      Width (m)
      <input 
        type="number" 
        step="0.1" 
        min="0.1" 
        value={width} 
        on:input={(e) => updateWidth(e.currentTarget.value)}
        class:error={widthError}
      />
      {#if widthError}
        <span class="error-text">{widthError}</span>
      {/if}
    </label>
    <label>
      Length (m)
      <input 
        type="number" 
        step="0.1" 
        min="0.1" 
        value={length} 
        on:input={(e) => updateLength(e.currentTarget.value)}
        class:error={lengthError}
      />
      {#if lengthError}
        <span class="error-text">{lengthError}</span>
      {/if}
    </label>
    <label>
      Height (m)
      <input 
        type="number" 
        step="0.1" 
        min="0.1" 
        value={height} 
        on:input={(e) => updateHeight(e.currentTarget.value)}
        class:error={heightError}
      />
      {#if heightError}
        <span class="error-text">{heightError}</span>
      {/if}
    </label>
  </div>
</div>

<style>
  .room-controls h3 {
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

  input.error {
    border-color: #f44336;
  }

  .error-text {
    font-size: 0.75rem;
    color: #f44336;
    margin-top: 2px;
  }
</style>
