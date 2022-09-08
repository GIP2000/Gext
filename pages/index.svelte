<script lang="ts">
  import Header from "../src/components/Header.svelte"
  export let Name: string; 
</script>
<Header />
<h1>This is the Index Page loaded through my route</h1>
<h2>Hello {Name}</h2>
