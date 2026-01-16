<script lang="ts">
    import StatItem from './StatItem.svelte';

    // Svelte 5: Props 使用 $props() 解构
    let { payday = 15 } = $props();
    
    // Svelte 5: 使用 $state 声明响应式状态
    let days = $state(0);

    function calculate() {
        const now = new Date();
        const currentDay = now.getDate();

        if (currentDay <= payday) {
            days = payday - currentDay;
        } else {
            const nextMonth = new Date(now.getFullYear(), now.getMonth() + 1, payday);
            days = Math.ceil((nextMonth.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
        }
    }

    // Svelte 5: 使用 $effect 替代 onMount
    $effect(() => {
        calculate();
        const timer = window.setInterval(calculate, 1000 * 60 * 10);  // 10min更新一次就够了
        return () => window.clearInterval(timer);
    });
</script>

<StatItem label="发薪" value={days} unit="天" />