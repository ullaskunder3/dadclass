<script>
  import { onMount } from "svelte";
  import {
    StartNotifications,
    StopNotifications,
    IsRunning,
  } from "../wailsjs/go/main/App";
  import { EventsOn, WindowHide } from "../wailsjs/runtime/runtime";

  const hideWindow = typeof WindowHide === "function" ? WindowHide : () => {};

  let interval = 30;
  let isRunning = false;
  let currentJoke = "";
  let notifPermission = "default";
  let listenerAttached = false;
  let jokeCount = 0;
  let statusMsg = "Ready for some laughs? 😄";

  function attachJokeListener() {
    if (listenerAttached) return;
    EventsOn("joke", (data) => {
      currentJoke = data;
      jokeCount++;
      statusMsg = `Latest joke delivered!`;

      if (notifPermission === "granted") {
        new Notification("🤣 Dad Joke", { body: data });
      }
    });
    listenerAttached = true;
  }

  async function start() {
    if (isRunning) return;
    if (interval < 5) interval = 5;
    statusMsg = "Starting jokes...";
    await StartNotifications(interval);
    isRunning = true;
    attachJokeListener();
    statusMsg = `Jokes every ${interval}s`;
  }

  async function stop() {
    if (!isRunning) return;
    await StopNotifications();
    isRunning = false;
    statusMsg = "Notifications paused";
  }

  function toggle() {
    isRunning ? stop() : start();
  }

  function minimizeToTray() {
    hideWindow();
  }

  onMount(async () => {
    notifPermission = Notification.permission;
    if (notifPermission === "default") {
      const result = await Notification.requestPermission();
      notifPermission = result;
    }
    
    try {
      isRunning = await IsRunning();
      if (isRunning) {
        statusMsg = `Active: Jokes every ${interval}s`;
        attachJokeListener();
      }
    } catch (_) {}
  });
</script>

<main>
  <header>
    <div class="logo">
      <span class="logo-emoji">⚡</span>
      <div class="logo-text">
        <span class="app-name">dadclass</span>
        <span class="app-sub">Premium Humor</span>
      </div>
    </div>
    <div class="header-right">
      <div class="status-pill" class:active={isRunning}>
        <span class="status-dot"></span>
        <span>{isRunning ? "Live" : "Idle"}</span>
      </div>
      <button class="tray-btn" title="Minimize to system tray" on:click={minimizeToTray}>
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="17 11 12 6 7 11"/>
          <polyline points="17 18 12 13 7 18"/>
        </svg>
      </button>
    </div>
  </header>

  <div class="joke-card" class:has-joke={currentJoke}>
    {#if currentJoke}
      <p class="joke-text">{currentJoke}</p>
      {#if jokeCount > 0}
        <span class="joke-badge">JOKE #{jokeCount}</span>
      {/if}
    {:else}
      <div class="joke-placeholder">
        <span class="placeholder-icon">🎭</span>
        <p>Your laugh is loading...<br/>Start notifications to begin.</p>
      </div>
    {/if}
  </div>

  <div class="controls">
    <div class="input-section">
      <span class="label">INTERVAL</span>
      <div class="input-row">
        <input
          type="number"
          bind:value={interval}
          min="5"
          disabled={isRunning}
        />
        <span class="unit">seconds</span>
      </div>
    </div>

    <button
      class="primary-btn"
      class:running={isRunning}
      on:click={toggle}
    >
      {#if isRunning}
        <svg viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>
        PAUSE
      {:else}
        <svg viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
        ACTIVATE
      {/if}
    </button>
  </div>

  <footer>
    <span class="status-msg">{statusMsg}</span>
    <span class="hint">Will minimize to tray on close</span>
  </footer>
</main>

<style>
  @import url("https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&display=swap");

  :global(*) {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  :global(html, body) {
    height: 100%;
    width: 100%;
    overflow: hidden;
    background: #0f0f12;
    font-family: 'Outfit', sans-serif;
    color: #fff;
  }

  main {
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding: 20px;
    background: linear-gradient(180deg, rgba(255,255,255,0.03) 0%, rgba(0,0,0,0) 100%);
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo-emoji {
    font-size: 24px;
    background: linear-gradient(135deg, #7c3aed, #9d50bb);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .logo-text {
    display: flex;
    flex-direction: column;
    line-height: 1;
  }

  .app-name {
    font-size: 18px;
    font-weight: 700;
    letter-spacing: -0.5px;
  }

  .app-sub {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    color: #6b7280;
    letter-spacing: 1px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .tray-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    transition: all 0.2s ease;
    padding: 0;
    flex-shrink: 0;
  }

  .tray-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
    border-color: rgba(255, 255, 255, 0.2);
  }

  .tray-btn:active {
    transform: scale(0.92);
  }

  .status-pill {
    padding: 6px 12px;
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 100px;
    font-size: 11px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 6px;
    color: #6b7280;
  }

  .status-pill.active {
    color: #10b981;
    background: rgba(16, 185, 129, 0.1);
    border-color: rgba(16, 185, 129, 0.2);
  }

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
  }

  .status-pill.active .status-dot {
    box-shadow: 0 0 8px #10b981;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0% { opacity: 1; transform: scale(1); }
    50% { opacity: 0.5; transform: scale(1.2); }
    100% { opacity: 1; transform: scale(1); }
  }

  .joke-card {
    flex: 1;
    background: rgba(255,255,255,0.03);
    border: 1px solid rgba(255,255,255,0.06);
    border-radius: 20px;
    padding: 30px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    position: relative;
    margin-bottom: 24px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .joke-card.has-joke {
    background: rgba(124, 58, 237, 0.05);
    border-color: rgba(124, 58, 237, 0.1);
  }

  .joke-text {
    font-size: 18px;
    line-height: 1.6;
    font-weight: 500;
    color: #e5e7eb;
  }

  .joke-badge {
    position: absolute;
    bottom: 20px;
    font-size: 10px;
    font-weight: 700;
    color: #4b5563;
    letter-spacing: 1.5px;
  }

  .placeholder-icon {
    font-size: 40px;
    margin-bottom: 16px;
    opacity: 0.2;
  }

  .joke-placeholder p {
    color: #4b5563;
    font-size: 14px;
  }

  .controls {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;
  }

  .input-section {
    flex: 1;
    background: rgba(255,255,255,0.03);
    border: 1px solid rgba(255,255,255,0.06);
    border-radius: 16px;
    padding: 12px 16px;
  }

  .label {
    font-size: 9px;
    font-weight: 700;
    color: #6b7280;
    display: block;
    margin-bottom: 4px;
  }

  .input-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  input[type="number"] {
    background: transparent;
    border: none;
    color: #fff;
    font-family: inherit;
    font-size: 18px;
    font-weight: 600;
    width: 60px;
    outline: none;
  }

  .unit {
    font-size: 12px;
    color: #4b5563;
    font-weight: 500;
  }

  .primary-btn {
    flex: 1;
    background: #7c3aed;
    color: #fff;
    border: none;
    border-radius: 16px;
    font-weight: 700;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(124, 58, 237, 0.3);
  }

  .primary-btn.running {
    background: #1f1f23;
    box-shadow: none;
    border: 1px solid rgba(255,255,255,0.1);
  }

  .primary-btn:hover {
    transform: translateY(-2px);
    background: #8b5cf6;
  }

  .primary-btn.running:hover {
    background: #27272a;
  }

  footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .status-msg {
    font-size: 12px;
    color: #4b5563;
    font-weight: 500;
  }

  .hint {
    font-size: 11px;
    color: #374151;
    font-weight: 600;
    font-style: italic;
  }
</style>
