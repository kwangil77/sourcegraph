<svelte:options immutable />

<script lang="ts">
    import { mdiAccount } from '@mdi/js'

    import Avatar from '$lib/Avatar.svelte'
    import Icon from '$lib/Icon.svelte'
    import { getOwnerDisplayName, getOwnerMatchURL, buildSearchURLQueryForOwner } from '$lib/search/results'
    import type { PersonMatch } from '$lib/shared'

    import SearchResult from './SearchResult.svelte'
    import { getSearchResultsContext } from './searchResultsContext'

    export let result: PersonMatch

    const queryState = getSearchResultsContext().queryState

    $: ownerURL = getOwnerMatchURL(result)
    $: displayName = getOwnerDisplayName(result)
    $: fileSearchQueryParams = buildSearchURLQueryForOwner($queryState, result)
</script>

<SearchResult>
    <Avatar
        slot="icon"
        avatar={{ displayName, username: result.user?.username ?? '', avatarURL: result.user?.avatarURL ?? null }}
        --avatar-size="1.5rem"
    />
    <div slot="title">
        &nbsp;
        {#if ownerURL}
            <a href={ownerURL}>{displayName}</a>
        {:else}
            {displayName}
        {/if}
        <span class="info">
            <Icon aria-label="Forked repository" svgPath={mdiAccount} inline />
            <small>Owner (person)</small>
        </span>
    </div>
    {#if fileSearchQueryParams}
        <p class="p-2 m-0">
            <a data-sveltekit-preload-data="tap" href="/search?{fileSearchQueryParams}">Show files</a>
        </p>
    {/if}
    {#if !result.user}
        <p class="p-2 m-0">
            <small class="font-italic"> This owner is not associated with any Sourcegraph user </small>
        </p>
    {/if}
</SearchResult>

<style lang="scss">
    .info {
        border-left: 1px solid var(--border-color);
        margin-left: 0.5rem;
        padding-left: 0.5rem;
    }
</style>
