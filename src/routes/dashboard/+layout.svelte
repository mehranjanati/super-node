<script lang="ts">
    import { onMount } from 'svelte';
    import { callManager } from '$lib/stores/callManager.svelte';
    import CallOverlay from '$lib/components/cockpit/CallOverlay.svelte';
    import * as sdk from 'matrix-js-sdk';

    // TODO: Replace with real auth
    const MOCK_ACCESS_TOKEN = import.meta.env.VITE_MATRIX_ACCESS_TOKEN || "your_access_token";
    const MOCK_USER_ID = import.meta.env.VITE_MATRIX_USER_ID || "@user:supernode.local";
    const HOMESERVER_URL = import.meta.env.VITE_MATRIX_HOMESERVER_URL || "http://localhost:8008";

    let { children } = $props();

    onMount(async () => {
        try {
            const client = sdk.createClient({
                baseUrl: HOMESERVER_URL,
                accessToken: MOCK_ACCESS_TOKEN,
                userId: MOCK_USER_ID
            });

            // Start client
            await client.startClient({ initialSyncLimit: 10 });

            // Initialize Call Manager
            callManager.init(client);
            console.log("Matrix Client started and CallManager initialized");
        } catch (e) {
            console.error("Failed to initialize Matrix client", e);
        }
    });
</script>

<div class="relative min-h-screen bg-gray-900 text-white">
    <!-- Global Call Overlay -->
    <CallOverlay />
    
    <!-- Main Content -->
    {@render children()}
</div>
