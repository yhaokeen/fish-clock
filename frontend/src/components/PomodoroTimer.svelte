<script lang="ts">
  // 番茄钟总时长（秒）
  const TOTAL_TIME = 25 * 60;
  
  // 当前剩余时间（秒）
  let remainingTime = $state(TOTAL_TIME);
  let isRunning = $state(false);
  let timer: number | null = null;

  // 计算进度 (0-1)，倒计时进度
  let countdownProgress = $derived(remainingTime / TOTAL_TIME);
  
  // 秒表进度 (0-1)，每分钟循环一次
  let secondsProgress = $derived((60 - (remainingTime % 60)) / 60);
  
  // 计算指针角度 (0-360度)，跟随秒针
  // +90 是为了补偿 SVG 整体旋转的 -90deg
  let pointerAngle = $derived(((60 - (remainingTime % 60)) / 60) * 360 + 90);
  
  // 圆周长（用于进度条）
  const OUTER_CIRCUMFERENCE = 2 * Math.PI * 85; // 倒计时圈 r=85
  const INNER_CIRCUMFERENCE = 2 * Math.PI * 72; // 秒表圈 r=72

  // 格式化时间显示
  let displayMinutes = $derived(Math.floor(remainingTime / 60));
  let displaySeconds = $derived(remainingTime % 60);

  // 开始/暂停
  function toggleTimer() {
    if (isRunning) {
      if (timer) {
        clearInterval(timer);
        timer = null;
      }
    } else {
      timer = window.setInterval(() => {
        if (remainingTime > 0) {
          remainingTime--;
        } else {
          clearInterval(timer!);
          timer = null;
          isRunning = false;
        }
      }, 1000);
    }
    isRunning = !isRunning;
  }

  // 重置
  function resetTimer() {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    isRunning = false;
    remainingTime = TOTAL_TIME;
  }
</script>

<div class="pomodoro-container">
  <!-- 圆形时钟 -->
  <div class="clock-wrapper">
    <svg viewBox="0 0 200 200" class="clock">
      <defs>
        <!-- 外圈阴影滤镜 -->
        <filter id="outerShadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="4" flood-color="rgba(0,0,0,0.1)" />
        </filter>
        <!-- 内凹阴影 -->
        <filter id="innerShadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="1" stdDeviation="2" flood-color="rgba(0,0,0,0.08)" />
        </filter>
        <!-- 倒计时渐变（蓝紫色） -->
        <linearGradient id="countdownGradient" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stop-color="#a8b5d6" />
          <stop offset="100%" stop-color="#7b8fc4" />
        </linearGradient>
        <!-- 秒表渐变（青绿色） -->
        <linearGradient id="secondsGradient" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stop-color="#7bccc4" />
          <stop offset="100%" stop-color="#4ecdc4" />
        </linearGradient>
        <!-- 中心点阴影 -->
        <filter id="centerShadow">
          <feDropShadow dx="0" dy="1" stdDeviation="1" flood-color="rgba(0,0,0,0.15)" />
        </filter>
      </defs>

      <!-- 最外层灰色圆环（带阴影） -->
      <circle 
        cx="100" cy="100" r="95" 
        fill="#f5f5f5" 
        filter="url(#outerShadow)"
      />
      
      <!-- 外圈轨道（灰色背景） -->
      <circle 
        cx="100" cy="100" r="85" 
        fill="none" 
        stroke="#e8e8e8" 
        stroke-width="12"
      />
      
      <!-- 倒计时进度条（蓝紫色，从顶部顺时针） -->
      <circle 
        cx="100" cy="100" r="85" 
        fill="none" 
        stroke="url(#countdownGradient)" 
        stroke-width="12"
        stroke-linecap="round"
        stroke-dasharray={OUTER_CIRCUMFERENCE}
        stroke-dashoffset={OUTER_CIRCUMFERENCE * countdownProgress}
        class="progress-ring countdown"
      />

      <!-- 内层白色表盘 -->
      <circle 
        cx="100" cy="100" r="72" 
        fill="white" 
        filter="url(#innerShadow)"
      />
      
      <!-- 秒表轨道（浅灰色背景） -->
      <circle 
        cx="100" cy="100" r="65" 
        fill="none" 
        stroke="#f0f0f0" 
        stroke-width="8"
      />
      
      <!-- 秒表进度条（青绿色） -->
      <circle 
        cx="100" cy="100" r="65" 
        fill="none" 
        stroke="url(#secondsGradient)" 
        stroke-width="8"
        stroke-linecap="round"
        stroke-dasharray={INNER_CIRCUMFERENCE}
        stroke-dashoffset={INNER_CIRCUMFERENCE * (1 - secondsProgress)}
        class="progress-ring seconds"
      />

      <!-- 指针 -->
      <line 
        x1="100" y1="100" x2="100" y2="35"
        stroke="#2c3e50"
        stroke-width="2"
        stroke-linecap="round"
        class="pointer"
        style="transform: rotate({pointerAngle}deg); transform-origin: 100px 100px;"
      />
      
      <!-- 中心圆点 -->
      <circle 
        cx="100" cy="100" r="6" 
        fill="white" 
        stroke="#e0e0e0"
        stroke-width="1"
        filter="url(#centerShadow)"
      />
      <circle cx="100" cy="100" r="3" fill="#ddd" />
    </svg>
  </div>

  <!-- 时间显示 -->
  <div class="time-display">
    <span class="minutes">{String(displayMinutes).padStart(2, '0')}</span>
    <span class="separator">:</span>
    <span class="seconds">{String(displaySeconds).padStart(2, '0')}</span>
  </div>

  <!-- 控制按钮 -->
  <div class="controls" style="--wails-draggable: no-drag;">
    <button class="btn primary" onclick={toggleTimer}>
      {isRunning ? '暂停' : '开始'}
    </button>
    <button class="btn secondary" onclick={resetTimer}>
      重置
    </button>
  </div>
</div>

<style>
  .pomodoro-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 15px;
    gap: 12px;
  }

  .clock-wrapper {
    width: 160px;
    height: 160px;
  }

  .clock {
    width: 100%;
    height: 100%;
    transform: rotate(-90deg); /* 从顶部开始 */
  }

  .progress-ring {
    transition: stroke-dashoffset 0.5s ease;
  }

  .pointer {
    transition: transform 0.3s ease;
  }

  .time-display {
    font-size: 28px;
    font-weight: 300;
    color: #333;
    font-family: "Helvetica Neue", "Arial", sans-serif;
    display: flex;
    align-items: baseline;
    letter-spacing: 2px;
  }

  .minutes, .seconds {
    min-width: 40px;
    text-align: center;
  }

  .separator {
    margin: 0 2px;
    color: #999;
  }

  .controls {
    display: flex;
    gap: 10px;
    margin-top: 5px;
  }

  .btn {
    padding: 8px 20px;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn.primary {
    background: linear-gradient(135deg, #4ecdc4, #44a3aa);
    color: white;
  }

  .btn.primary:hover {
    transform: scale(1.05);
    box-shadow: 0 2px 8px rgba(78, 205, 196, 0.4);
  }

  .btn.secondary {
    background: #f0f0f0;
    color: #666;
  }

  .btn.secondary:hover {
    background: #e0e0e0;
  }
</style>
