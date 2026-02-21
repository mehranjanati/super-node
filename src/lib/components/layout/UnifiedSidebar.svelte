<script lang="ts">
    import { callManager } from '$lib/stores/callManager.svelte';
    import { Activity, Phone, User } from 'lucide-svelte';

    // Status logic
    // We consider connected if matrix client is initialized (in a real app we'd check sync state)
    let isConnected = $derived(callManager.isInitialized);
    let statusColor = $derived(isConnected ? 'bg-green-500' : 'bg-red-500');
    
    // Agent status
    let agentStatus = $state("Agent is ready");

    $effect(() => {
        if (callManager.callState === 'CONNECTED') {
            agentStatus = "Agent is on a call";
        } else if (callManager.callState === 'RINGING') {
            agentStatus = "Incoming call...";
        } else {
            agentStatus = "Agent is ready";
        }
    });
</script>

<aside class="w-64 border-r border-white/10 bg-black/40 backdrop-blur-md p-4 flex flex-col gap-6 h-full min-h-screen">
    <!-- Header -->
    <div class="flex items-center gap-3">
        <div class="h-8 w-8 rounded-lg bg-blue-600 flex items-center justify-center">
            <Activity size={18} class="text-white" />
        </div>
        <h1 class="font-bold text-white tracking-wider">SuperNode</h1>
    </div>

    <!-- Agent Status -->
    <div class="rounded-xl bg-white/5 p-4 border border-white/10">
        <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-400">AGENT STATUS</span>
            <div class={`h-2 w-2 rounded-full ${statusColor} animate-pulse`}></div>
        </div>
        <p class="text-sm text-gray-200">{agentStatus}</p>
        {#if callManager.callState === 'CONNECTED'}
            <p class="text-xs text-green-400 mt-1 flex items-center gap-1">
                <Phone size={10} /> On Call
            </p>
        {/if}
    </div>

    <!-- Active Calls -->
    <div class="flex-1">
        <h3 class="text-xs font-semibold text-gray-500 mb-3 uppercase tracking-wider">Active Calls</h3>
        
        {#if callManager.callState !== 'IDLE'}
            <div class="flex items-center gap-3 rounded-lg bg-white/5 p-3 border border-white/10">
                <div class="h-8 w-8 rounded-full bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center text-xs font-bold">
                    {callManager.callerId ? callManager.callerId[0] : '?'}
                </div>
                <div class="overflow-hidden">
                    <p class="truncate text-sm font-medium text-white">{callManager.callerId || "Unknown"}</p>
                    <p class="text-xs text-gray-400">{callManager.callState}</p>
                </div>
            </div>
        {:else}
            <p class="text-sm text-gray-600 italic">No active calls</p>
        {/if}
    </div>

    <!-- Footer / User -->
    <div class="mt-auto border-t border-white/10 pt-4 flex items-center gap-3">
        <div class="h-8 w-8 rounded-full bg-gray-700 flex items-center justify-center">
            <User size={16} />
        </div>
        <div>
            <p class="text-sm font-medium text-white">Admin User</p>
            <p class="text-xs text-gray-500">Online</p>
        </div>
    </div>
</aside>
