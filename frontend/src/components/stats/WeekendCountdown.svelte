<script lang="ts">
    import StatItem from './StatItem.svelte';

    // Svelte 5: 使用 $state 声明响应式状态
    let days = $state(0);
    
    function calculate() {
        const dayOfWeek = new Date().getDay();
        if (dayOfWeek === 0) days = 5;
        else if (dayOfWeek === 6) days = 6;
        else days = 5 - dayOfWeek;
    }

    // Svelte 5: 使用 $effect 替代 onMount
    $effect(() => {
        calculate();
        const timer = window.setInterval(calculate, 1000 * 60 * 10);
        return () => window.clearInterval(timer);
    });
</script>

<StatItem label="周五" value={days} unit="天" />
