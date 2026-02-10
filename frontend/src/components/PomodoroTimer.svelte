<script lang="ts">
  // 番茄钟总时长（秒）
  const TOTAL_TIME = 25 * 60;
  
  // 当前剩余时间（秒）
  let remainingTime = $state(TOTAL_TIME);
  let isRunning = $state(false);
  let timer: number | null = null;

  // 实时时钟时间
  let currentTime = $state(new Date());

  // 计算进度 (0-1)，倒计时进度
  let countdownProgress = $derived(remainingTime / TOTAL_TIME);
  
  // 秒表进度 (0-1)，每分钟循环一次
  let secondsProgress = $derived((60 - (remainingTime % 60)) / 60);
  
  // 圆周长（用于进度条）
  const OUTER_CIRCUMFERENCE = 2 * Math.PI * 98; // 倒计时圈 r=85

  // 格式化时间显示
  let displayMinutes = $derived(Math.floor(remainingTime / 60));
  let displaySeconds = $derived(remainingTime % 60);

  // 计算时针角度
  // 时针：每小时30度，每分钟0.5度，每12秒0.1度
  let hourAngle = $derived.by(() => {
    const hours = currentTime.getHours() % 12; // 12小时制
    const minutes = currentTime.getMinutes();
    const seconds = currentTime.getSeconds();
    return hours * 30 + minutes * 0.5 + seconds * 0.1 / 12;
  });

  // 计算分针角度
  // 分针：每分钟6度，每秒0.1度
  let minuteAngle = $derived.by(() => {
    const minutes = currentTime.getMinutes();
    const seconds = currentTime.getSeconds();
    return minutes * 6 + seconds * 0.1;
  });

  // 计算秒针角度
  // 秒针：每秒6度
  let secondAngle = $derived.by(() => {
    const seconds = currentTime.getSeconds();
    return seconds * 6;
  });

  // 根据分针角度计算进度条起点坐标
  let startX = $derived.by(() => {
    return 125 + 98 * Math.cos(minuteAngle * Math.PI / 180);
  });
  let startY = $derived.by(() => {
    return 125 + 98 * Math.sin(minuteAngle * Math.PI / 180);
  });

  // 实时更新时间
  let clockTimer: number | null = null;
  
  $effect(() => {
    // 立即更新一次
    currentTime = new Date();
    
    // 每秒更新一次
    clockTimer = window.setInterval(() => {
      currentTime = new Date();
    }, 1000);

    // 清理函数
    return () => {
      if (clockTimer) {
        clearInterval(clockTimer);
        clockTimer = null;
      }
    };
  });

</script>

<div class="pomodoro-container">
  <!-- 圆形时钟 -->
  <div class="clock-wrapper">
    <svg viewBox="-10 -10 270 270" class="clock">
      <defs>
        <!-- 外圈阴影滤镜 -->
        <filter id="outerShadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="rgba(0,0,0,0.2)" />
        </filter>
        <!-- 内凹阴影 -->
        <filter id="innerShadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="1" stdDeviation="2" flood-color="rgba(0,0,0,0.08)" />
        </filter>
        <!-- 倒计时渐变（蓝紫色） -->
        <linearGradient id="countdownGradient" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stop-color="#9cb0e5" />
          <stop offset="100%" stop-color="#77e7d7" />
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

        <!-- 指针四周阴影 -->
        <filter id="dropShadow" x="-50%" y="-50%" width="200%" height="200%">
          <feDropShadow dx="1" dy="1" stdDeviation="3" flood-color="rgba(0,0,0,0.6)" />
          <feDropShadow dx="-0.5" dy="-0.5" stdDeviation="0.5" flood-color="rgba(0,0,0,0.3)" />
        </filter>
        <filter id="second" x="-50%" y="-50%" width="200%" height="200%">
          <feDropShadow dx="1" dy="1" stdDeviation="3" flood-color="rgba(0,0,0,0.3)" />
        </filter>
      </defs>

      <!-- 最外层灰色圆环（带阴影） r=125，圆心坐标125, 125-->
      <circle 
        cx="125" cy="125" r="125" 
        fill="#FFFFFF"
        stroke="#E2E8F0" 
        stroke-width="0.8"
        filter="url(#outerShadow)"
      />

      <!-- 中间灰色圆环（带阴影） r=85，圆心坐标125, 125-->
      <circle 
        cx="125" cy="125" r="85" 
        fill="none" 
        stroke="#E2E8F0" 
        stroke-width="3"
        opacity="0.81"
      />
      <!-- 秒针表盘 r=30，圆心坐标125, 180-->
      <circle
        cx="125" cy="180" r="30" 
        fill="none" 
        stroke="#E2E8F0" 
        stroke-width="3"
        opacity="0.81"
      />

      <!-- 12个刻度 -->
      <!-- 12点 (长刻度) -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#333" stroke-width="2" stroke-linecap="round" />
      <!-- 1点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(30deg); transform-origin: 125px 125px;" />
      <!-- 2点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(60deg); transform-origin: 125px 125px;" />
      <!-- 3点 (长刻度) -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#333" stroke-width="2" stroke-linecap="round" style="transform: rotate(90deg); transform-origin: 125px 125px;" />
      <!-- 4点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(120deg); transform-origin: 125px 125px;" />
      <!-- 5点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(150deg); transform-origin: 125px 125px;" />
      <!-- 6点 (长刻度) -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#333" stroke-width="2" stroke-linecap="round" style="transform: rotate(180deg); transform-origin: 125px 125px;" />
      <!-- 7点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(210deg); transform-origin: 125px 125px;" />
      <!-- 8点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(240deg); transform-origin: 125px 125px;" />
      <!-- 9点 (长刻度) -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#333" stroke-width="2" stroke-linecap="round" style="transform: rotate(270deg); transform-origin: 125px 125px;" />
      <!-- 10点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(300deg); transform-origin: 125px 125px;" />
      <!-- 11点 -->
      <line x1="125" y1="0" x2="125" y2="3" stroke="#999" stroke-width="1.5" stroke-linecap="round" style="transform: rotate(330deg); transform-origin: 125px 125px;" />

      <!-- 倒计时进度条（蓝紫色，从顶部顺时针） -->
      <!-- 使用 path 代替 circle，可以自定义起始点 -->
      <!-- M 100,15 表示从12点钟方向开始（圆心100,100，半径85，所以顶部是100,15） -->
      <!-- M 125,27 A 98,98 0 0,1 223,125 表示从(125,27)点开始画一个98px的正圆弧到(223,125)
       0：x轴旋转角度 0：表示小弧线 ，1：表示顺时针方向-->
      <path 
        d="M 125,27 A 98,98 0 0,1 223,125"
        fill="none" 
        stroke="url(#countdownGradient)" 
        stroke-width="25"
        stroke-dasharray={OUTER_CIRCUMFERENCE}
        stroke-dashoffset={OUTER_CIRCUMFERENCE * (1 - countdownProgress)}
        class="progress-ring countdown"
      />
      <!-- 秒针 -->
      <rect 
        x="124" y="153" 
        width="2" height="27"
        rx="2" ry="2"
        stroke="none"
        fill="url(#countdownGradient)"
        filter="url(#second)"
        style="transform: rotate(0deg); transform-origin: 125px 180px;"
      />
      
      <!-- 时针 -->
      <rect 
        x="124" y="40" 
        width="2" height="95"
        rx="2" ry="2"
        stroke="#FFFFFF"
        stroke-width="0.1"
        fill="#F8F8F8"
        filter="url(#dropShadow)"
        style="transform: rotate({hourAngle}deg); transform-origin: 125px 125px; will-change: transform;"
        class="hour-hand"
      />
      
      <!-- 分针（整体） -->
      <g 
        filter="url(#dropShadow)"
        style="transform: rotate({minuteAngle}deg); transform-origin: 125px 125px; will-change: transform;"
        class="minute-hand"
      >
        <!-- 分针主体 -->
        <rect 
          x="124" y="14" 
          width="2" height="120"
          rx="2" ry="2"
          fill="#F8F8F8"
        />
        <!-- 分针末端深色部分 -->
        <rect 
          x="124" y="14" 
          width="2" height="25"
          rx="2" ry="2"
          fill="#111111"
        />
      </g>

      
      <!-- 中心圆点 -->
      <circle cx="125" cy="125" r="8" fill="#FFFFFF" filter="url(#centerShadow)" />
      <!-- <circle cx="125" cy="125" r="2" fill="url(#metal3D)" /> -->
      <circle cx="125" cy="180" r="5" fill="#FFFFFF" filter="url(#centerShadow)" />
    </svg>
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
    width: 200px;
    height: 200px;
  }

  .clock {
    width: 100%;
    height: 100%;
    /* transform: rotate(-90deg); */
  }

  .progress-ring {
    transition: stroke-dashoffset 0.5s ease;
  }

  .pointer {
    transition: transform 0.3s ease;
  }

  .minute-hand {
    /* 去掉 transition 避免与频繁更新冲突 */
    transform-style: preserve-3d;
    backface-visibility: hidden;
  }

  .hour-hand {
    /* 去掉 transition 避免与频繁更新冲突 */
    transform-style: preserve-3d;
    backface-visibility: hidden;
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
