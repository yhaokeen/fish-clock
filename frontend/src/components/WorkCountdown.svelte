<script lang="ts">
  import CountdownTimer from "./stats/CountdownTimer.svelte";
  import PaydayCountdown from "./stats/PaydayCountdown.svelte";
  import WeekendCountdown from "./stats/WeekendCountdown.svelte";
  import TodayEarnings from "./stats/TodayEarnings.svelte";
  import FestivalCountdown from "./stats/FestivalCountdown.svelte";
  import { configStore } from "../stores/config";

  // Svelte 5: 使用 $derived 派生值
  let workStartHour = $derived($configStore?.workStartHour ?? 9);
  let workEndHour = $derived($configStore?.workEndHour ?? 18);
  let payday = $derived($configStore?.payday ?? 25);
  let monthlySalary = $derived($configStore?.monthlySalary ?? 15000);
</script>

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

<style>
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
