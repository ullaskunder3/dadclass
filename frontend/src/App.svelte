<script>
  import { onMount } from "svelte";
  import {
    StartNotifications,
    StopNotifications,
  } from "../wailsjs/go/main/App";
  import {
    EventsOn,
    WindowGetSize,
    WindowSetSize,
  } from "../wailsjs/runtime/runtime";

  let interval = 10;
  // let joke = "";
  let listenerAttached = false;

  function start() {
    StartNotifications(interval);

    // Attach listener only once
    if (!listenerAttached) {
      EventsOn("joke", (data) => {
        // joke = data;
        if (Notification.permission === "granted") {
          new Notification("🤣 Dad Joke", { body: data });
        }
      });
      listenerAttached = true;
    }
  }

  function stop() {
    StopNotifications();
    // Optional: You could also set listenerAttached = false here
    // if StopNotifications also stops the Go event from firing
  }

  onMount(() => {
    if (Notification.permission !== "granted") {
      Notification.requestPermission();
    }
  });
  // Wait for DOM to render
  window.addEventListener("DOMContentLoaded", () => {
    const width = document.body.offsetWidth;
    const height = document.body.offsetHeight;

    // Add some padding if needed
    WindowSetSize(width + 20, height);
  });
</script>

<main>
  <h1>Dad Joke Notifier</h1>

  <div class="controls">
    <p>
      Send me a new dad joke every
      <input type="number" bind:value={interval} min="5" /> seconds.
    </p>
    <div class="button-group">
      <button class="mac-btn start" on:click={start}>Start</button>
      <button class="mac-btn stop" on:click={stop}>Stop</button>
    </div>
  </div>

  <!-- {#if joke}
    <div class="joke-glass">{joke}</div>
  {/if} -->
</main>

<style>
  @import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap");

  main {
    font-family:
      "Inter",
      -apple-system,
      BlinkMacSystemFont,
      sans-serif;
    padding-block: 2rem;
    text-align: center;
    color: #ccc;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-width: 200px;
  }

  h1 {
    font-size: 2rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
  }

  .controls {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  input {
    padding: 0.4rem 0.6rem;
    font-size: 1rem;
    border-radius: 0.6rem;
    border: 1px solid #ccc;
    backdrop-filter: blur(4px);
  }

  .button-group {
    display: flex;
    gap: 1rem;
    margin-top: 1rem;
  }

  .mac-btn {
    padding: 0.6rem 1.2rem;
    border-radius: 12px;
    font-size: 1rem;
    font-weight: 500;
    border: none;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition:
      transform 0.15s ease,
      box-shadow 0.2s ease;
    backdrop-filter: blur(10px);
  }

  .start {
    background: #007aff;
    color: white;
  }

  .stop {
    background: #d0d3d8;
    color: #333;
  }

  .mac-btn:hover {
    transform: scale(1.05);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  }

  .joke-glass {
    background: rgba(255, 255, 255, 0.15);
    backdrop-filter: blur(20px) saturate(150%);
    -webkit-backdrop-filter: blur(20px) saturate(150%);
    border-radius: 16px;
    border: 1px solid rgba(207, 3, 3, 0.3);
    padding: 2rem;
    color: #111;
    font-size: 1.2rem;
  }
</style>
