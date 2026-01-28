<script lang="ts">
  // Svelte 5: Props 使用 $props() 解构
  let { offWorkHour = 18, offWorkMinute = 0, title = "下班还有" } = $props();
  
  // Svelte 5: 使用 $state 声明响应式状态
  let countdown = $state("00:00:00");
  let timer = $state<number | undefined>(undefined);

  function updateCountdown() {
    const now = new Date();
    const target = new Date();
    target.setHours(offWorkHour, offWorkMinute, 0, 0);
    let diff = target.getTime() - now.getTime();
    if (diff < 0) {
      countdown = "00:00:00";
      return;
    }

    const hours = Math.floor(diff / (1000 * 60 * 60));
    diff -= hours * (1000 * 60 * 60);
    const minutes = Math.floor(diff / (1000 * 60));
    diff -= minutes * (1000 * 60);
    const seconds = Math.floor(diff / 1000);
    
    countdown = [hours, minutes, seconds].map(num => num.toString().padStart(2, '0')).join(':')
  }

  // Svelte 5: 使用 $effect 替代 onMount 和 onDestroy
  $effect(() => {
    updateCountdown();
    timer = window.setInterval(updateCountdown, 1000);
    
    // 返回清理函数（相当于 onDestroy）
    return () => {
      if (timer) {
        window.clearInterval(timer);
      }
    };
  });
</script>

<div class="countdown-container">
  <div class="header">{title}</div>
  <div class="countdown">{countdown}</div>
</div>

<style>
    .countdown-container {
      text-align: center;
    }
    
    .header {
      font-size: 14px;
      color: #515151;
      margin-bottom: 8px;
    }
    
    .countdown {
      font-size: 48px;
      font-weight: bold;
      color: #333;
      font-family: "Consolas", monospace;
    }
  </style>