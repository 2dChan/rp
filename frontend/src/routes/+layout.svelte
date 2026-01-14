<!--
  Copyright (C) 2026 Andrey Kriulin
  Licensed under the GNU Affero General Public License v3.0 or later.
  See the LICENSE file in the project root for the full license text.
-->

<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
