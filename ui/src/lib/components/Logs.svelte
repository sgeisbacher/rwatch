<script lang="ts">
	import type { LogEvent } from '$lib/core/webrtc-manager.svelte';

	let showLogs = $state(false);

	interface Props {
		logs: LogEvent[];
	}

	function toggleShowLogs() {
		showLogs = !showLogs;
	}
	const { logs }: Props = $props();
</script>

<div
	id="logsContainer"
	onclick={toggleShowLogs}
	onkeyup={toggleShowLogs}
	role="button"
	tabindex="0"
>
	Logs ({logs.length})
	{#if !showLogs}
		-
	{/if}
	<span class="defensiveText">
		{#if !showLogs}
			{logs.length > 0 ? logs[logs.length - 1].msg : ''}
		{/if}
	</span><br />
	{#if showLogs}
		<ul>
			{#each logs as log}
				<li>
					<span class="defensiveText">{log.timestamp.toLocaleTimeString()}</span>
					{log.level.toUpperCase()} - {log.msg}
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	#logsContainer {
		margin-top: 30px;
	}
	.defensiveText {
		color: gray;
	}
	ul {
		list-style-type: none;
	}
</style>
