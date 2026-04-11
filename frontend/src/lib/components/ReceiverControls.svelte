<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let x: number;
  export let y: number;
  export let z: number;
  export let fieldErrors: Record<string, string> = {};

  const dispatch = createEventDispatcher();

  function updateX(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateX', { value: num });
  }

  function updateY(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateY', { value: num });
  }

  function updateZ(value: string) {
    const num = parseFloat(value) || 0;
    dispatch('updateZ', { value: num });
  }

  $: receiverError = fieldErrors['receiver'] || '';
</script>

<div class="receiver-controls">
  <h3>Listening Position</h3>
  <div class="control-group">
    <label>
      X (m)
      <input 
        type="number" 
        step="0.1" 
        min="0" 
        value={x} 
        on:input={(e) => updateX(e.currentTarget.value)}
        class:error={receiverError}
      />
    </label>
    <label>
      Y (m)
      <input 
        type="number" 
        step="0.1" 
        min="0" 
        value={y} 
        on:input={(e) => updateY(e.currentTarget.value)}
        class:error={receiverError}
      />
    </label>
    <label>
      Z (m)
      <input 
        type="number" 
        step="0.1" 
        min="0" 
        value={z} 
        on:input={(e) => updateZ(e.currentTarget.value)}
        class:error={receiverError}
      />
    </label>
    {#if receiverError}
      <span class="error-text">{receiverError}</span>
    {/if}
  </div>
</div>

<style>
  .receiver-controls h3 {
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
