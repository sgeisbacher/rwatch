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

<header>
	<!-- Fixed navbar -->
	<nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
		<a class="navbar-brand" href="#">rwatch</a>
	</nav>
</header>
<!-- Begin page content -->
<main role="main" class="container">
	<div class="cover-container d-flex h-100 p-3 flex-column">
		<Terminal {data} />
		<Logs {logs} />
	</div>
</main>

<footer class="footer">
	<div class="container">
		<span class="text-muted">Connection ...</span>
	</div>
</footer>

<style>
	.bg-dark {
		background-color: #1e1e1e !important;
	}
	.footer {
		background-color: #1e1e1e !important;
	}
</style>
