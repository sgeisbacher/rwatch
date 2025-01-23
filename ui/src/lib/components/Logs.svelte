<script lang="ts">
	import type { LogEvent } from '$lib/core/webrtc-manager.svelte';

	interface Props {
		logs: LogEvent[];
	}
	const { logs }: Props = $props();
</script>

<div>
	Logs:<br />
	<ul>
    {#each logs as log}
		<li>{log.level[0].toUpperCase()} - {log.msg}</li>
    {/each}
	</ul>
</div>
