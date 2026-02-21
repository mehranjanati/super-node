<script lang="ts">
    let userId = "";
    let amount = 0;
    let sourceWallet = "";
    let loading = false;
    let result: any = null;
    let error = "";

    async function submitTransaction() {
        loading = true;
        error = "";
        result = null;

        try {
            const response = await fetch('http://localhost:3000/api/transaction', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    user_id: userId,
                    amount: Number(amount),
                    source_wallet: sourceWallet
                })
            });

            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || 'Transaction failed');
            }

            result = await response.json();
        } catch (e: any) {
            error = e.message;
        } finally {
            loading = false;
        }
    }
</script>

<div class="container mx-auto p-8 max-w-md">
    <h1 class="text-2xl font-bold mb-6 text-gray-800">New Transaction</h1>
    
    <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="user_id">
                User ID
            </label>
            <input 
                class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                id="user_id" 
                type="text" 
                placeholder="User123"
                bind:value={userId}
            >
        </div>
        
        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="amount">
                Amount ($)
            </label>
            <input 
                class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                id="amount" 
                type="number" 
                placeholder="1000.00"
                bind:value={amount}
            >
        </div>
        
        <div class="mb-6">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="wallet">
                Source Wallet (Hex)
            </label>
            <input 
                class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                id="wallet" 
                type="text" 
                placeholder="0x123..."
                bind:value={sourceWallet}
            >
            <p class="text-xs text-gray-500 mt-1">Try starting with '0xDEAD' to trigger AML rejection.</p>
        </div>
        
        <div class="flex items-center justify-between">
            <button 
                class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50" 
                type="button"
                on:click={submitTransaction}
                disabled={loading}
            >
                {#if loading}
                    Processing...
                {:else}
                    Submit Transaction
                {/if}
            </button>
        </div>
    </div>

    {#if result}
        <div class={`border-l-4 p-4 ${result.status === 'COMPLETED' ? 'bg-green-100 border-green-500 text-green-700' : 'bg-red-100 border-red-500 text-red-700'}`}>
            <p class="font-bold">Status: {result.status}</p>
            <p>Transaction ID: {result.transaction_id}</p>
            <p>{result.message}</p>
        </div>
    {/if}

    {#if error}
        <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4" role="alert">
            <p class="font-bold">Error</p>
            <p>{error}</p>
        </div>
    {/if}
</div>
