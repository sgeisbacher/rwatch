<script lang="ts">
	interface Props {
		status: string;
	}
	const { status }: Props = $props();
</script>

<span id="status">{status}</span>

<style>
	#status {
		margin-left: auto;
		margin-right: 15px;
	}
</style>
