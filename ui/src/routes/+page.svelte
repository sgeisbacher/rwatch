<script lang="ts">
	import { onMount } from 'svelte';
	import Logs from '$lib/components/Logs.svelte';
	import { initConnection, type LogEvent } from '$lib/core/webrtc-manager.svelte';
	import Terminal from '$lib/components/Terminal.svelte';

	const logs = $state<LogEvent[]>([]);
	let data = $state<string[]>([]);
	onMount(async () => {
		await initConnection(
			(dataLines) => {
				data = dataLines;
			},
			(logEvent) => logs.push(logEvent)
		);
	});
</script>

<Terminal {data} />
<Logs {logs} />
