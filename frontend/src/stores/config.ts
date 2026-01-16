import { writable } from 'svelte/store';
import { ConfigService } from '../../bindings/fish-clock/app/index';
import type { Config } from '../../bindings/fish-clock/app/models';
import { Events } from '@wailsio/runtime';

// 创建配置 store
function createConfigStore() {
  const { subscribe, set, update } = writable<Config | null>(null);

  // 监听配置更新事件（从其他窗口发来的）
  Events.On('config:update', (ev: any) => {
    console.log('收到配置更新事件:', ev.data);
    set(ev.data as Config);
  });

  return {
    subscribe,
    set,
    update,
    // 从后端加载配置
    async load() {
      const config = await ConfigService.GetConfig();
      if (config) {
        set(config);
      }
      return config;
    },
    // 保存配置到后端
    async save() {
      let currentConfig: Config | null = null;
      subscribe(value => currentConfig = value)();
      
      if (currentConfig) {
        await ConfigService.UpdateConfig(currentConfig);
      }
    },
    // 广播配置更新（通知其他窗口）
    broadcast(config: Config) {
      Events.Emit('config:update', config);
    }
  };
}

export const configStore = createConfigStore();
