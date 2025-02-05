<script lang="ts">
	interface Props {
		data: string[];
	}
	const { data }: Props = $props();
</script>

<div>
	Terminal-Output:<br />
	{#if data.length > 0}
		<pre id="terminal">{#each data as line}{line}{/each}</pre>
	{/if}
</div>

<style>
	pre#terminal {
		background-color: black;
		color: white;
		padding: 10px;
	}
</style>
