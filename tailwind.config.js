/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.html",
    "./internal/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // Platform primary colors - blue, white, black theme
        primary: {
          50: '#f8fafc',   // Near white
          100: '#f1f5f9',  // Light gray-blue
          200: '#e2e8f0',  // Light gray
          300: '#cbd5e1',  // Medium gray
          400: '#94a3b8',  // Dark gray
          500: '#64748b',  // Slate
          600: '#475569',  // Dark slate
          700: '#334155',  // Darker slate
          800: '#1e293b',  // Very dark slate
          900: '#0f172a',  // Near black
        },
        // Platform accent blue
        platform: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',  // Primary platform blue
          600: '#2563eb',  // Darker blue for buttons/links
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'pulse-slow': 'pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' }
        }
      }
    },
  },
  plugins: [
    // Add form plugin for better form styling
    require('@tailwindcss/forms')
  ],
  // Safelist important classes that might be generated dynamically
  safelist: [
    // Guild color variants
    {
      pattern: /bg-(platform|primary)-(50|100|200|300|400|500|600|700|800|900)/,
      variants: ['hover', 'focus', 'active'],
    },
    {
      pattern: /text-(platform|primary)-(50|100|200|300|400|500|600|700|800|900)/,
      variants: ['hover', 'focus', 'active'],
    },
    {
      pattern: /border-(platform|primary)-(50|100|200|300|400|500|600|700|800|900)/,
      variants: ['hover', 'focus', 'active'],
    },
    // HTMX-related classes
    'htmx-indicator',
    'htmx-request',
    'loading-spinner',
    'fade-in',
    'slide-up',
  ]
}