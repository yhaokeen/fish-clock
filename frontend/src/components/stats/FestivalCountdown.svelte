<script lang="ts">
    import {GetNextFestival} from '../../../bindings/fish-clock/app/festivalservice';
    import StatItem from './StatItem.svelte';

    // Svelte 5: 使用 $state 声明响应式状态
    let festivalName = $state("");
    let daysLeft = $state(0);

    async function loadFestival() {
        const info = await GetNextFestival();
        festivalName = info.name;
        daysLeft = info.days;
    }

    // Svelte 5: 使用 $effect 替代 onMount
    $effect(() => {
        loadFestival();
        const timer = window.setInterval(loadFestival, 1000 * 60 * 60);
        return () => window.clearInterval(timer);
    });
</script>

<StatItem label={festivalName} value={daysLeft} unit="天" />
