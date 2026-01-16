<script lang="ts">
    import StatItem from './StatItem.svelte';

    // Svelte 5: Props 使用 $props() 解构
    let { monthlySalary = 10000, workStartHour = 9, workEndHour = 18 } = $props();

    // Svelte 5: 使用 $state 声明响应式状态
    let earnings = $state("0.000");
    let timer = $state<number | undefined>(undefined);

    function calculate() {
        const now = new Date();
        const dailySalary = monthlySalary / 22;
        const workStart = new Date(now.getFullYear(), now.getMonth(), now.getDate(), workStartHour, 0, 0);
        const workEnd = new Date(now.getFullYear(), now.getMonth(), now.getDate(), workEndHour, 0, 0);

        if (now >= workStart && now < workEnd) {
            const worked = now.getTime() - workStart.getTime();
            const totalWorked = workEnd.getTime() - workStart.getTime();
            earnings = (dailySalary * worked / totalWorked).toFixed(3);
        } else if (now >= workEnd) {
            earnings = dailySalary.toFixed(3);
        } else {
            earnings = "0.000";
        }
    }

    // Svelte 5: 使用 $effect 替代 onMount 和 onDestroy
    $effect(() => {
        calculate();
        timer = window.setInterval(calculate, 100);
        
        return () => {
            if (timer) {
                window.clearInterval(timer);
            }
        };
    });
</script>


<StatItem label="今天赚了" value={earnings} unit="¥" />
