<script lang="ts">
  import CountdownTimer from "./components/CountdownTimer.svelte";
  import PaydayCountdown from "./components/stats/PaydayCountdown.svelte";
  import WeekendCountdown from "./components/stats/WeekendCountdown.svelte";
  import TodayEarnings from "./components/stats/TodayEarnings.svelte";
  import FestivalCountdown from "./components/stats/FestivalCountdown.svelte";
  import { configStore } from "./stores/config";

  // Svelte 5: 使用 $derived 派生值
  let workStartHour = $derived($configStore?.workStartHour ?? 9);
  let workEndHour = $derived($configStore?.workEndHour ?? 18);
  let payday = $derived($configStore?.payday ?? 25);
  let monthlySalary = $derived($configStore?.monthlySalary ?? 15000);
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
    {#if $configStore}
      {#key $configStore}
        <div class="content">
          <CountdownTimer 
            offWorkHour={workEndHour}
            offWorkMinute={0} 
            title="下班还有" 
          />
          <div class="stats">
            <PaydayCountdown {payday} />
            <WeekendCountdown />
            <FestivalCountdown />
            <TodayEarnings 
              {monthlySalary}
              {workStartHour}
              {workEndHour}
            />
          </div>
        </div>
      {/key}
    {:else}
      <div class="loading">加载中...</div>
    {/if}
  </div>
</main>

<style>
  main {
    width: 100vw;
    height: 100vh;
    margin: 0;
    padding: 0;
    background: transparent;
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
  }

  .loading {
    font-size: 14px;
    color: #999;
  }

  .stats {
    display: flex;
    flex-direction: row;
    gap: 15px;
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
  }
</style>
