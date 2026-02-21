<script lang="ts">
    import { callManager } from '$lib/stores/callManager.svelte';
    import { fade, fly } from 'svelte/transition';
    import { Phone, PhoneOff, ArrowRight } from 'lucide-svelte';
</script>

{#if callManager.callState === 'RINGING'}
    <div 
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xl"
        in:fade={{ duration: 300 }}
        out:fade={{ duration: 200 }}
        role="dialog"
        aria-modal="true"
    >
        <div 
            class="relative w-full max-w-md overflow-hidden rounded-3xl bg-white/10 p-8 shadow-2xl border border-white/20 text-center"
            in:fly={{ y: 50, duration: 400 }}
        >
            <!-- Avatar / Caller Info -->
            <div class="mb-8 flex flex-col items-center">
                <div class="relative mb-4 flex h-24 w-24 items-center justify-center rounded-full bg-gradient-to-tr from-purple-500 to-blue-500 shadow-lg animate-pulse">
                    <span class="text-3xl font-bold text-white">
                        {callManager.callerId ? callManager.callerId[0].toUpperCase() : '?'}
                    </span>
                    <!-- Pulsing rings -->
                    <div class="absolute inset-0 -z-10 animate-ping rounded-full bg-blue-500 opacity-30"></div>
                </div>
                <h2 class="text-2xl font-semibold text-white tracking-wide">
                    {callManager.callerId || "Unknown Caller"}
                </h2>
                <p class="text-blue-200">Incoming Video Call...</p>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-between gap-6 px-4">
                <!-- Decline -->
                <button 
                    onclick={() => callManager.rejectCall()}
                    class="group flex flex-col items-center gap-2 transition-transform hover:scale-105"
                    aria-label="Decline Call"
                >
                    <div class="flex h-14 w-14 items-center justify-center rounded-full bg-red-500/80 text-white shadow-lg backdrop-blur-md transition-colors group-hover:bg-red-500">
                        <PhoneOff size={24} />
                    </div>
                    <span class="text-xs font-medium text-red-200">Decline</span>
                </button>

                <!-- Forward / Mobile -->
                <button 
                    class="group flex flex-col items-center gap-2 transition-transform hover:scale-105"
                    aria-label="Forward to Mobile"
                >
                    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-yellow-500/80 text-white shadow-lg backdrop-blur-md transition-colors group-hover:bg-yellow-500">
                        <ArrowRight size={20} />
                    </div>
                    <span class="text-xs font-medium text-yellow-200">Mobile</span>
                </button>

                <!-- Answer -->
                <button 
                    onclick={() => callManager.acceptCall()}
                    class="group flex flex-col items-center gap-2 transition-transform hover:scale-105"
                    aria-label="Answer Call"
                >
                    <div class="flex h-16 w-16 items-center justify-center rounded-full bg-green-500/80 text-white shadow-lg backdrop-blur-md transition-colors group-hover:bg-green-500 animate-bounce">
                        <Phone size={28} />
                    </div>
                    <span class="text-xs font-medium text-green-200">Answer</span>
                </button>
            </div>
        </div>
    </div>
{/if}
