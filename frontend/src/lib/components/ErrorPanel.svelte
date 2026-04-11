<script lang="ts">
  export let error: string | null = null;
</script>

{#if error}
  <div class="error-panel">
    <div class="error-message">{error}</div>
  </div>
{/if}

<style>
  .error-panel {
    padding: 1rem;
    background: #3d1f1f;
    border: 1px solid #8b3a3a;
    border-radius: 4px;
    margin-top: 1rem;
  }
  .error-message {
    color: #ff6b6b;
    font-size: 0.85rem;
    line-height: 1.4;
  }
</style>
