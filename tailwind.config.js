/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {
            colors: {
                'nexus-bg': '#0a0a0a',
                'nexus-card': '#1a1a1a',
                'nexus-accent': '#00ffcc',
            },
        },
    },
    plugins: [],
}
