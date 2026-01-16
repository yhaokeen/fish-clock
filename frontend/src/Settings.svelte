<script lang="ts">
  import * as WindowService from "../bindings/fish-clock/app/windowservice";
  import type { Config } from "../bindings/fish-clock/app/models";
  import { configStore } from "./stores/config";

  // Svelte 5: 使用 $state 声明响应式状态
  let originalConfig = $state<Config | null>(null);
  let isReady = $state(false);
  
  // 使用本地变量来绑定表单
  let workStartHour = $state(9);
  let workEndHour = $state(18);
  let payday = $state(25);
  let monthlySalary = $state(15000);
  let opacityPercent = $state(90);

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
      payday = originalConfig.payday;
      monthlySalary = originalConfig.monthlySalary;
      opacityPercent = Math.round(originalConfig.opacity * 100);
    }
    WindowService.CloseCurrentWindow();
  }
</script>

<div class="settings-container">
  <div class="settings-content">
    <div class="header">
      <h2 class="title">设置</h2>
      <button class="close-btn" onclick={handleCancel} aria-label="关闭">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
          <path d="M15 5L5 15M5 5L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
    </div>

    {#if !isReady}
      <div class="loading">加载中...</div>
    {:else if $configStore}
      <div class="form-container">
        <div class="form-row">
          <label for="work-start-hour">工作时间</label>
          <div class="time-range">
            <select 
              id="work-start-hour" 
              bind:value={workStartHour}
              onchange={handleConfigChange}
            >
              {#each Array(24) as _, i}
                <option value={i}>{String(i).padStart(2, "0")}:00</option>
              {/each}
            </select>
            <span class="separator">-</span>
            <select 
              id="work-end-hour" 
              bind:value={workEndHour} 
              onchange={handleConfigChange}
              aria-label="下班时间"
            >
              {#each Array(24) as _, i}
                <option value={i}>{String(i).padStart(2, "0")}:00</option>
              {/each}
            </select>
          </div>
        </div>

        <div class="form-row">
          <label for="payday">发薪日</label>
          <div class="input-with-unit">
            <input
              id="payday"
              type="number"
              bind:value={payday}
              oninput={handleConfigChange}
              min="1"
              max="31"
            />
            <span class="unit">号</span>
          </div>
        </div>

        <div class="form-row">
          <label for="monthly-salary">月薪</label>
          <div class="input-with-unit">
            <input
              id="monthly-salary"
              type="number"
              bind:value={monthlySalary}
              oninput={handleConfigChange}
              min="0"
              step="100"
            />
            <span class="unit">¥</span>
          </div>
        </div>

        <div class="form-row opacity-row">
          <label for="opacity-slider">窗口透明度</label>
          <div class="opacity-control">
            <span class="opacity-value">{opacityPercent}%</span>
            <input
              id="opacity-slider"
              type="range"
              min="0"
              max="100"
              value={opacityPercent}
              oninput={handleOpacityChange}
              class="slider"
            />
            <div class="opacity-toggle">
              <input
                type="checkbox"
                id="opacity-toggle"
                checked={opacityPercent > 50}
                class="toggle-checkbox"
                disabled
              />
              <label for="opacity-toggle" class="toggle-label"></label>
            </div>
          </div>
        </div>
      </div>

      <div class="button-group">
        <button class="btn btn-cancel" onclick={handleCancel}>取消</button>
        <button class="btn btn-save" onclick={handleSave}>保存</button>
      </div>
    {:else}
      <div class="error">无法加载配置</div>
    {/if}
  </div>
</div>

<style>
  .settings-container {
    width: 100vw;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: rgba(0, 0, 0, 0.5);
    padding: 20px;
    box-sizing: border-box;
  }

  .settings-content {
    background: white;
    border-radius: 16px;
    padding: 28px 32px 32px;
    width: 100%;
    max-width: 480px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 28px;
  }

  .title {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: #1a1a1a;
  }

  .close-btn {
    background: transparent;
    border: none;
    width: 32px;
    height: 32px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #666;
    transition: all 0.2s;
  }

  .close-btn:hover {
    background: #f5f5f5;
    color: #333;
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
    gap: 18px;
  }

  .form-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 20px;
  }

  .form-row label {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    min-width: 100px;
    flex-shrink: 0;
  }

  .time-range {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
  }

  select {
    flex: 1;
    padding: 10px 14px;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    font-size: 14px;
    background: #fafafa;
    cursor: pointer;
    outline: none;
    transition: all 0.2s;
    color: #333;
  }

  select:hover {
    border-color: #d0d0d0;
    background: #f5f5f5;
  }

  select:focus {
    border-color: #2196f3;
    background: white;
  }

  .separator {
    color: #999;
    font-size: 14px;
  }

  .input-with-unit {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
  }

  .input-with-unit input {
    flex: 1;
    padding: 10px 14px;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    font-size: 14px;
    background: #fafafa;
    outline: none;
    transition: all 0.2s;
    color: #333;
  }

  .input-with-unit input:hover {
    border-color: #d0d0d0;
    background: #f5f5f5;
  }

  .input-with-unit input:focus {
    border-color: #2196f3;
    background: white;
  }

  .unit {
    color: #999;
    font-size: 14px;
    min-width: 20px;
  }

  .opacity-row {
    padding-top: 8px;
  }

  .opacity-control {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
  }

  .opacity-value {
    min-width: 42px;
    font-weight: 600;
    color: #333;
    font-size: 14px;
  }

  .slider {
    flex: 1;
    height: 4px;
    border-radius: 2px;
    background: #e0e0e0;
    outline: none;
    appearance: none;
    -webkit-appearance: none;
    cursor: pointer;
  }

  .slider::-webkit-slider-thumb {
    appearance: none;
    -webkit-appearance: none;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: #2196f3;
    cursor: pointer;
    transition: all 0.2s;
  }

  .slider::-webkit-slider-thumb:hover {
    transform: scale(1.2);
    box-shadow: 0 0 0 8px rgba(33, 150, 243, 0.1);
  }

  .slider::-moz-range-thumb {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: #2196f3;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
  }

  .slider::-moz-range-thumb:hover {
    transform: scale(1.2);
  }

  .opacity-toggle {
    position: relative;
  }

  .toggle-checkbox {
    display: none;
  }

  .toggle-label {
    display: block;
    width: 44px;
    height: 24px;
    background: #2196f3;
    border-radius: 12px;
    position: relative;
    cursor: pointer;
    transition: background 0.3s;
  }

  .toggle-label::after {
    content: '';
    position: absolute;
    top: 2px;
    left: 2px;
    width: 20px;
    height: 20px;
    background: white;
    border-radius: 50%;
    transition: transform 0.3s;
    transform: translateX(20px);
  }

  .toggle-checkbox:not(:checked) + .toggle-label {
    background: #ddd;
  }

  .toggle-checkbox:not(:checked) + .toggle-label::after {
    transform: translateX(0);
  }

  .button-group {
    display: flex;
    gap: 12px;
    margin-top: 32px;
  }

  .btn {
    flex: 1;
    padding: 12px 24px;
    border: none;
    border-radius: 10px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel {
    background: #f5f5f5;
    color: #666;
  }

  .btn-cancel:hover {
    background: #e8e8e8;
    color: #333;
  }

  .btn-save {
    background: #2196f3;
    color: white;
  }

  .btn-save:hover {
    background: #1976d2;
    box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
  }

  .btn:active {
    transform: scale(0.98);
  }
</style>
