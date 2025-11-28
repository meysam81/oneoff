/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  darkMode: "class",
  theme: {
    screens: {
      'xs': '375px',
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
      '2xl': '1536px',
    },
    extend: {
      colors: {
        bg: {
          primary: "#0a0a0b",
          secondary: "#111113",
          tertiary: "#18181b",
          elevated: "#1f1f23",
        },
        fg: {
          primary: "#fafafa",
          secondary: "#a1a1aa",
          muted: "#52525b",
        },
        accent: {
          primary: "#22d3ee",
          secondary: "#06b6d4",
          glow: "rgba(34, 211, 238, 0.15)",
        },
        success: "#4ade80",
        warning: "#fbbf24",
        error: "#f87171",
        border: {
          subtle: "rgba(255, 255, 255, 0.06)",
          default: "rgba(255, 255, 255, 0.1)",
          accent: "rgba(34, 211, 238, 0.3)",
        },
      },
      fontFamily: {
        display: ["JetBrains Mono", "monospace"],
        body: ["Geist", "-apple-system", "BlinkMacSystemFont", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      fontSize: {
        display: ["4rem", { lineHeight: "1.1", letterSpacing: "-0.02em" }],
        headline: ["2.5rem", { lineHeight: "1.2", letterSpacing: "-0.01em" }],
        title: ["1.5rem", { lineHeight: "1.3" }],
        body: ["1rem", { lineHeight: "1.6" }],
        small: ["0.875rem", { lineHeight: "1.5" }],
        caption: ["0.75rem", { lineHeight: "1.4", letterSpacing: "0.02em" }],
      },
      animation: {
        "glow-pulse": "glow-pulse 2s ease-in-out infinite",
        typing: "typing 3.5s steps(40, end)",
        blink: "blink 1s step-end infinite",
        "fade-in": "fade-in 0.5s ease-out",
        "slide-up": "slide-up 0.5s ease-out",
      },
      keyframes: {
        "glow-pulse": {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0.5" },
        },
        typing: {
          from: { width: "0" },
          to: { width: "100%" },
        },
        blink: {
          "50%": { borderColor: "transparent" },
        },
        "fade-in": {
          from: { opacity: "0" },
          to: { opacity: "1" },
        },
        "slide-up": {
          from: { opacity: "0", transform: "translateY(20px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
      },
      boxShadow: {
        glow: "0 0 20px rgba(34, 211, 238, 0.15), 0 0 40px rgba(34, 211, 238, 0.15)",
        "glow-lg":
          "0 0 30px rgba(34, 211, 238, 0.2), 0 0 60px rgba(34, 211, 238, 0.1)",
      },
    },
  },
  plugins: [],
};
