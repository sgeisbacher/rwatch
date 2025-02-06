<script lang="ts">
	interface Props {
		data: string[];
	}
	const { data }: Props = $props();
</script>

<div>
	<div id="terminalContainer">
		{#if data.length > 0}
			<pre id="terminal">{#each data as line}{line}{/each}</pre>
		{/if}
		<br />
	</div>
</div>

<style>
	#terminalContainer {
		background-color: black;
		border-width: 3px;
		border-color: #8d5dff;
		border-style: solid;
		border-radius: 2px;
		min-height: 500px;
	}
	pre#terminal {
		background-color: black;
		color: white;
		padding: 15px 20px;
	}
</style>
