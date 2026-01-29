<script lang="ts">
  import * as WindowService from "../bindings/fish-clock/app/windowservice";
  import type { Config } from "../bindings/fish-clock/app/models";
  import { configStore } from "./stores/config";
  import { Button, Input, Label, Range, Timepicker } from 'flowbite-svelte';
  import { ClockSolid, CalendarMonthSolid, DollarOutline, AdjustmentsHorizontalOutline } from 'flowbite-svelte-icons';

  // Svelte 5: 使用 $state 声明响应式状态
  let originalConfig: Config | null = null;  // 不用 $state，避免被 effect 重复赋值
  let isInitialized = false;  // 标记是否已初始化
  
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
      updateConfig();
    }
  }

  $effect(() => {
    (async () => {
      // 如果配置还没加载，先加载
      if (!$configStore) {
        await configStore.load();
      }
      
      // 只在首次加载时保存原始配置和初始化本地变量
      if ($configStore && !isInitialized) {
        isInitialized = true;
        originalConfig = JSON.parse(JSON.stringify($configStore));
        workStartHour = $configStore.workStartHour;
        workEndHour = $configStore.workEndHour;
        selectedTimeRange = {
          time: formatTime($configStore.workStartHour),
          endTime: formatTime($configStore.workEndHour)
        };
        payday = $configStore.payday;
        monthlySalary = $configStore.monthlySalary;
        opacityPercent = Math.round((1 - $configStore.opacity) * 100);
      }
      
    })();
  });

  function updateConfig() {
    const newConfig = {
      workStartHour,
      workEndHour,
      payday,
      monthlySalary,
      opacity: 1 - opacityPercent / 100
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
      // 广播原始配置到主窗口，让主窗口也恢复
      configStore.broadcast(originalConfig);
    }
    WindowService.CloseCurrentWindow();
  }
</script>

<div class="w-screen h-screen bg-gray-50 rounded-2xl p-6 flex flex-col shadow-xl" style="--wails-draggable:drag;">
  <div class="mb-6">
    <h2 class="text-lg font-semibold text-gray-900 mb-3">设置</h2>
    <div class="title-divider"></div>
  </div>

  <div class="flex flex-col gap-5 flex-1 overflow-y-auto pr-1" style="--wails-draggable: no-drag;">
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2 min-w-[100px] shrink-0">
        <ClockSolid class="w-4 h-4 text-indigo-500" />
        <Label class="text-sm font-medium text-gray-600 !mb-0">工作时间</Label>
      </div>
      <Timepicker
        type="range"
        onselect={handleRangeChange}
        value={selectedTimeRange.time}
        endValue={selectedTimeRange.endTime}
        divClass="shadow-none gap-2"
      />
    </div>

    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2 min-w-[100px] shrink-0">
        <CalendarMonthSolid class="w-4 h-4 text-indigo-500" />
        <Label for="payday" class="text-sm font-medium text-gray-600 !mb-0">发薪日</Label>
      </div>
      <div class="flex items-center gap-2">
        <Input
          id="payday"
          type="number"
          bind:value={payday}
          oninput={handleConfigChange}
          min="1"
          max="31"
          size="sm"
          class="w-24 !bg-gray-100 !border-gray-200 !rounded-lg focus:!border-indigo-400 focus:!ring-indigo-100"
        />
        <span class="text-sm text-gray-500 font-medium">号</span>
      </div>
    </div>

    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2 min-w-[100px] shrink-0">
        <DollarOutline class="w-4 h-4 text-indigo-500" />
        <Label for="monthly-salary" class="text-sm font-medium text-gray-600 !mb-0">月薪</Label>
      </div>
      <div class="flex items-center gap-2">
        <Input
          id="monthly-salary"
          type="number"
          bind:value={monthlySalary}
          oninput={handleConfigChange}
          min="0"
          step="100"
          size="sm"
          class="w-24 !bg-gray-100 !border-gray-200 !rounded-lg focus:!border-indigo-400 focus:!ring-indigo-100"
        />
        <span class="text-sm text-gray-500 font-medium">¥</span>
      </div>
    </div>

    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2 min-w-[100px] shrink-0">
        <AdjustmentsHorizontalOutline class="w-4 h-4 text-indigo-500" />
        <Label for="opacity-slider" class="text-sm font-medium text-gray-600 !mb-0">透明度</Label>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-sm font-semibold text-gray-600 min-w-[42px] text-right">{opacityPercent}%</span>
        <Range
          id="opacity-slider"
          bind:value={opacityPercent}
          oninput={handleOpacityChange}
          min="0"
          max="100"
          size="sm"
          class="w-40 accent-indigo-500"
        />
      </div>
    </div>
  </div>

  <div class="flex gap-3 mt-5 pt-5 border-t border-gray-100 justify-end shrink-0" style="--wails-draggable: no-drag;">
    <Button color="alternative" size="sm" onclick={handleCancel} class="!px-5">取消</Button>
    <Button color="primary" size="sm" onclick={handleSave} class="!px-5 !bg-indigo-500 hover:!bg-indigo-600">保存</Button>
  </div>
</div>

<style>
  /* 标题下方渐变分隔线 - Tailwind 不支持渐变背景到透明 */
  .title-divider {
    height: 2px;
    background: linear-gradient(90deg, #6366f1 0%, #a855f7 50%, transparent 100%);
    border-radius: 1px;
    width: 60%;
  }

  /* 隐藏 Timepicker 多余的时钟按钮 */
  :global(button[aria-label*="time picker"]) {
    display: none !important;
  }

  /* Timepicker 输入框样式 */
  :global(input[type="time"]) {
    background: rgb(243 244 246) !important;
    border: 1px solid rgb(229 231 235) !important;
    border-radius: 0.5rem !important;
  }

  :global(input[type="time"]:focus) {
    border-color: rgb(129 140 248) !important;
    --tw-ring-color: rgb(224 231 255) !important;
  }
</style>
