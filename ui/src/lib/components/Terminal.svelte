<script lang="ts">
	import type { ExecutionInfo } from '$lib/utils/types';

	interface Props {
		execution: ExecutionInfo;
	}
	const { execution }: Props = $props();
</script>

<div>
	<div id="terminalContainer">
		{#if execution}
			<pre id="terminal">{execution.output}</pre>
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
