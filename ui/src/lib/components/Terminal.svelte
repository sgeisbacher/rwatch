<script lang="ts">
	interface Props {
		data: string[];
	}
	const { data }: Props = $props();
</script>

<style>
      pre#terminal {
        background-color: black;
        color: white;
        padding: 10px;
      }
</style>

<div>
	Terminal-Output:<br />
	<pre id="terminal">
    {#each data as line}
		{line}
    {/each}
	</pre>
</div>
