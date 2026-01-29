<script lang="ts">
    export let label: string;
    export let value: string | number;
    export let unit: string;
</script>

<div class="stat-item">
    <div class="label">{label}</div>
    <div class="value">{value}<span class="unit">{unit}</span></div>
</div>

<style>
    .stat-item {
        text-align: center;
        padding: 0 10px;
    }
    .label {
        font-size: 12px;
        color: #434343;
        margin-bottom: 4px;
    }
    .value {
        font-size: 18px;
        font-weight: bold;
        color: #333;
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        /* min-width: 60px; */
        display: inline-block;
    }
    .unit {
        font-size: 12px;
        font-weight: normal;
        color: #666;
        margin-left: 4px;
    }
</style>