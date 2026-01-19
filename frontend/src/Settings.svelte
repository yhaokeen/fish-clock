<script lang="ts">
  import * as WindowService from "../bindings/fish-clock/app/windowservice";
  import type { Config } from "../bindings/fish-clock/app/models";
  import { configStore } from "./stores/config";
  import { Button, Input, Label, Range, Timepicker } from 'flowbite-svelte';

  // Svelte 5: 使用 $state 声明响应式状态
  let originalConfig = $state<Config | null>(null);
  let isReady = $state(false);
  
  // 使用本地变量来绑定表单
  let workStartHour = $state(9);
  let workEndHour = $state(18);
  let payday = $state(25);
  let monthlySalary = $state(15000);
  let opacityPercent = $state(90);

  // 时间范围对象用于 Timepicker range
  let selectedTimeRange = $state({ time: "09:00", endTime: "18:00" });

  // 格式化时间为 HH:MM
  function formatTime(hour: number): string {
    return `${String(hour).padStart(2, "0")}:00`;
  }

  // 从时间字符串解析小时
  function parseHour(timeStr: string): number {
    const hour = parseInt(timeStr.split(':')[0]);
    return isNaN(hour) ? 0 : hour;
  }

  // 处理时间范围变化
  function handleRangeChange(data: { time: string; endTime: string; [key: string]: string }): void {
    if (data) {
      selectedTimeRange = {
        time: data.time,
        endTime: data.endTime
      };
      // 更新小时数并更新配置
      workStartHour = parseHour(data.time);
      workEndHour = parseHour(data.endTime);
      if (isReady) {
        updateConfig();
      }
    }
  }

  // Svelte 5: 使用 $effect 替代 onMount
  $effect(() => {
    (async () => {
      // 如果配置还没加载，先加载
      if (!$configStore) {
        await configStore.load();
      }
      
      // 保存当前配置的副本并初始化本地变量
      if ($configStore) {
        originalConfig = JSON.parse(JSON.stringify($configStore));
        workStartHour = $configStore.workStartHour;
        workEndHour = $configStore.workEndHour;
        selectedTimeRange = {
          time: formatTime($configStore.workStartHour),
          endTime: formatTime($configStore.workEndHour)
        };
        payday = $configStore.payday;
        monthlySalary = $configStore.monthlySalary;
        opacityPercent = Math.round($configStore.opacity * 100);
      }
      
      isReady = true;
    })();
  });


  function updateConfig() {
    const newConfig = {
      workStartHour,
      workEndHour,
      payday,
      monthlySalary,
      opacity: opacityPercent / 100
    };
    
    console.log('更新配置:', newConfig);
    
    configStore.update(cfg => {
      if (cfg) {
        const updatedConfig = { ...cfg, ...newConfig };
        // 广播配置更新到其他窗口
        configStore.broadcast(updatedConfig);
        return updatedConfig;
      }
      return cfg;
    });
  }

  function handleOpacityChange(event: Event) {
    const target = event.target as HTMLInputElement;
    opacityPercent = parseInt(target.value);
    updateConfig();
  }

  function handleConfigChange() {
    updateConfig();
  }

  async function handleSave() {
    if (!$configStore) return;
    
    // 先更新一次确保最新值
    updateConfig();
    
    // 保存到后端
    await configStore.save();
    WindowService.CloseCurrentWindow();
  }

  function handleCancel() {
    // 恢复原始配置
    if (originalConfig) {
      configStore.set(originalConfig);
      // 同时恢复本地变量
      workStartHour = originalConfig.workStartHour;
      workEndHour = originalConfig.workEndHour;
      selectedTimeRange = {
        time: formatTime(originalConfig.workStartHour),
        endTime: formatTime(originalConfig.workEndHour)
      };
      payday = originalConfig.payday;
      monthlySalary = originalConfig.monthlySalary;
      opacityPercent = Math.round(originalConfig.opacity * 100);
    }
    WindowService.CloseCurrentWindow();
  }
</script>

<div class="settings-content" style="--wails-draggable:drag;">
  <h2 class="title">设置</h2>

  {#if !isReady}
    <div class="loading">加载中...</div>
  {:else if $configStore}
    <div class="form-container">
      <div class="form-row">
        <Label class="form-label">工作时间</Label>
        <Timepicker
          type="range"
          onselect={handleRangeChange}
          value={selectedTimeRange.time}
          endValue={selectedTimeRange.endTime}
          divClass="shadow-none timepicker-no-extra-icon"
        />
      </div>

      <div class="form-row">
        <Label for="payday" class="form-label">发薪日</Label>
        <div class="input-with-unit">
          <Input
            id="payday"
            type="number"
            bind:value={payday}
            oninput={handleConfigChange}
            min="1"
            max="31"
            size="sm"
            class="flex-1"
          />
          <span class="unit">号</span>
        </div>
      </div>

      <div class="form-row">
        <Label for="monthly-salary" class="form-label">月薪</Label>
        <div class="input-with-unit">
          <Input
            id="monthly-salary"
            type="number"
            bind:value={monthlySalary}
            oninput={handleConfigChange}
            min="0"
            step="100"
            size="sm"
            class="flex-1"
          />
          <span class="unit">¥</span>
        </div>
      </div>

      <div class="form-row">
        <Label for="opacity-slider" class="form-label">窗口透明度</Label>
        <div class="opacity-control">
          <span class="opacity-value">{opacityPercent}%</span>
          <Range
            id="opacity-slider"
            bind:value={opacityPercent}
            oninput={handleOpacityChange}
            min="0"
            max="100"
            size="sm"
            class="flex-1"
          />
        </div>
      </div>
    </div>

    <div class="button-group">
      <Button color="light" onclick={handleCancel} class="flex-1">取消</Button>
      <Button color="blue" onclick={handleSave} class="flex-1">保存</Button>
    </div>
  {:else}
    <div class="error">无法加载配置</div>
  {/if}
</div>

<style>
  .settings-content {
    width: 100vw;
    height: 100vh;
    background: white;
    border-radius: 16px;
    padding: 20px 20px 20px;
    box-sizing: border-box;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    display: flex;
    flex-direction: column;
  }

  .title {
    margin: 0 0 28px 0;
    font-size: 20px;
    font-weight: 600;
    color: #1a1a1a;
  }

  .loading {
    text-align: center;
    padding: 60px 20px;
    font-size: 15px;
    color: #999;
  }

  .error {
    text-align: center;
    padding: 60px 20px;
    font-size: 15px;
    color: #f44336;
  }

  .form-container {
    display: flex;
    flex-direction: column;
    gap: 24px;
    flex: 1;
    overflow-y: auto;
    padding-right: 8px;
  }
  
  .form-container::-webkit-scrollbar {
    width: 6px;
  }
  
  .form-container::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .form-container::-webkit-scrollbar-thumb {
    background: #ddd;
    border-radius: 3px;
  }
  
  .form-container::-webkit-scrollbar-thumb:hover {
    background: #bbb;
  }

  .form-row {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .form-row :global(.form-label) {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    min-width: 100px;
    flex-shrink: 0;
    margin-bottom: 0;
  }

  .input-with-unit {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  /* 输入框设置固定宽度 */
  .input-with-unit :global(.flex-1) {
    flex: none;
    width: 120px;
  }

  .unit {
    color: #999;
    font-size: 14px;
    min-width: 24px;
    font-weight: 500;
  }

  .opacity-control {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  /* Range 滑块设置固定宽度 */
  .opacity-control :global(.flex-1) {
    flex: none;
    width: 200px;
  }

  .opacity-value {
    min-width: 48px;
    font-weight: 600;
    color: #333;
    font-size: 14px;
  }

  .button-group {
    display: flex;
    gap: 12px;
    margin-top: 24px;
    padding-top: 24px;
    border-top: 1px solid #f0f0f0;
    flex-shrink: 0;
  }

  /* 自定义 Flowbite 组件样式 */
  :global(.flex-1) {
    flex: 1;
  }

  /* 隐藏 Timepicker range 模式中 Flowbite 添加的额外时钟按钮 */
  /* 浏览器原生 input[type="time"] 已经有时钟图标了 */
  .form-row :global(button[aria-label*="time picker"]) {
    display: none !important;
  }
</style>
