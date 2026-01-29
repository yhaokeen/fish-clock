<script lang="ts">
  import WorkCountdown from "./components/WorkCountdown.svelte";
  import PomodoroTimer from "./components/PomodoroTimer.svelte";
  import { configStore } from "./stores/config";

  // 当前视图：'pomodoro' | 'work'，先写死显示番茄钟
  const currentView = 'pomodoro';

  // Svelte 5: 使用 $derived 派生值
  let opacity = $derived($configStore?.opacity ?? 0.6);

  // Svelte 5: 使用 $effect 替代 onMount
  $effect(() => {
    configStore.load();
  });
  
  // Svelte 5: 使用 $effect 监听配置变化
  $effect(() => {
    if ($configStore) {
      console.log('App 收到配置更新:', $configStore);
    }
  });
</script>

<main>
  <div 
    class="card" 
    style="--wails-draggable:drag; background: rgba(255, 255, 255, {opacity});"
  >
    {#if currentView === 'pomodoro'}
      <PomodoroTimer />
    {:else}
      <WorkCountdown />
    {/if}
  </div>
</main>

<style>
  main {
    width: 100vw;
    height: 100vh;
    margin: 0;
    padding: 0;
    border-radius: 16px;
  }

  .card {
    width: 100%;
    height: 100%;
    border-radius: 16px;
    padding: 0 0;
    box-sizing: border-box;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    color: #333;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    /* 毛玻璃模糊效果 */
    /* backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px); */
  }

</style>
