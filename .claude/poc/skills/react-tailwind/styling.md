# Tailwind CSS Styling Conventions

Advanced Tailwind patterns, utilities, and conventions.

## Configuration

### tailwind.config.js

```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class', // Enable class-based dark mode
  content: [
    './src/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // Custom color palette
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554',
        },
      },
      fontFamily: {
        sans: ['Inter var', 'system-ui', 'sans-serif'],
        mono: ['Fira Code', 'monospace'],
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      borderRadius: {
        '4xl': '2rem',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
```

### globals.css

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  /* Custom base styles */
  body {
    @apply bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100;
  }
}

@layer components {
  /* Custom component classes */
  .btn-primary {
    @apply px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2;
  }

  .card {
    @apply bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg shadow-sm;
  }
}

@layer utilities {
  /* Custom utility classes */
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
}
```

## Responsive Design Patterns

### Mobile-First Breakpoints

```tsx
// Stack on mobile, row on desktop
<div className="flex flex-col md:flex-row gap-4">
  <div className="w-full md:w-1/2">Left</div>
  <div className="w-full md:w-1/2">Right</div>
</div>

// Hidden on mobile, shown on desktop
<nav className="hidden lg:flex lg:gap-6">
  <a href="#">Link 1</a>
  <a href="#">Link 2</a>
</nav>

// Different padding per breakpoint
<div className="p-4 sm:p-6 lg:p-8">
  Content
</div>

// Different grid columns per breakpoint
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  {items.map(item => <Card key={item.id} />)}
</div>
```

### Container Queries (Tailwind 3.4+)

```tsx
// Enable in tailwind.config.js
module.exports = {
  plugins: [require('@tailwindcss/container-queries')],
}

// Usage:
<div className="@container">
  <div className="@lg:flex @lg:gap-4">
    {/* Flex layout when container is large, stack when small */}
  </div>
</div>
```

## Dark Mode Patterns

### Component-Level Dark Mode

```tsx
// Simple dark mode toggle
<div className="bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100">
  Content
</div>

// Dark mode with borders
<div className="border border-slate-200 dark:border-slate-700">
  Content
</div>

// Dark mode with hover states
<button className="bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600">
  Click me
</button>
```

### Theme Provider Pattern

```tsx
"use client"

import { createContext, useContext, useEffect, useState } from "react"

type Theme = "light" | "dark" | "system"

const ThemeContext = createContext<{
  theme: Theme
  setTheme: (theme: Theme) => void
}>({
  theme: "system",
  setTheme: () => {},
})

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>("system")

  useEffect(() => {
    const stored = localStorage.getItem("theme") as Theme
    if (stored) setTheme(stored)
  }, [])

  useEffect(() => {
    const root = document.documentElement

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      root.classList.toggle("dark", systemTheme === "dark")
    } else {
      root.classList.toggle("dark", theme === "dark")
    }

    localStorage.setItem("theme", theme)
  }, [theme])

  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

export const useTheme = () => useContext(ThemeContext)
```

## Animation & Transitions

### Built-in Animations

```tsx
// Spin
<div className="animate-spin">⚙️</div>

// Ping
<div className="animate-ping">📡</div>

// Pulse
<div className="animate-pulse">💓</div>

// Bounce
<div className="animate-bounce">⬇️</div>
```

### Custom Transitions

```tsx
// Hover transitions
<button className="transition-colors hover:bg-blue-600">
  Hover me
</button>

<button className="transition-all hover:scale-105 hover:shadow-lg">
  Scale on hover
</button>

// Complex transitions
<div className="transition-all duration-300 ease-in-out transform hover:translate-x-2 hover:shadow-xl">
  Slide and shadow
</div>
```

### Custom Animations

```js
// In tailwind.config.js
module.exports = {
  theme: {
    extend: {
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(0)' },
        },
      },
      animation: {
        fadeIn: 'fadeIn 0.3s ease-in',
        slideIn: 'slideIn 0.3s ease-out',
      },
    },
  },
}
```

```tsx
// Usage
<div className="animate-fadeIn">Fades in</div>
<div className="animate-slideIn">Slides in</div>
```

## Advanced Layout Patterns

### Sticky Header

```tsx
<header className="sticky top-0 z-50 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
  <nav className="container mx-auto px-4 py-4">
    Navigation
  </nav>
</header>
```

### Centered Modal

```tsx
<div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl max-w-md w-full mx-4">
    <div className="p-6">
      Modal content
    </div>
  </div>
</div>
```

### Sidebar Layout

```tsx
<div className="flex h-screen">
  {/* Sidebar */}
  <aside className="w-64 bg-slate-100 dark:bg-slate-800 border-r border-slate-200 dark:border-slate-700 overflow-y-auto">
    Sidebar
  </aside>

  {/* Main content */}
  <main className="flex-1 overflow-y-auto">
    <div className="container mx-auto p-6">
      Content
    </div>
  </main>
</div>
```

### Grid with Gap

```tsx
// Auto-fit: Creates as many columns as will fit
<div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
  {items.map(item => <Card key={item.id} />)}
</div>

// Auto-fill: Creates columns even if empty
<div className="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-6">
  {items.map(item => <Card key={item.id} />)}
</div>
```

## Typography Patterns

### Responsive Font Sizes

```tsx
<h1 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold">
  Responsive Heading
</h1>

<p className="text-sm sm:text-base lg:text-lg">
  Responsive paragraph
</p>
```

### Typography Plugin

```tsx
// Install: npm install @tailwindcss/typography
// Enable in tailwind.config.js plugins: [require('@tailwindcss/typography')]

// Usage: Automatically style markdown/HTML content
<article className="prose dark:prose-invert lg:prose-xl">
  {/* All child elements get beautiful typography */}
  <h1>Heading</h1>
  <p>Paragraph with <a href="#">links</a></p>
  <ul>
    <li>List items</li>
  </ul>
</article>
```

## Form Styling

### Forms Plugin

```js
// Install: npm install @tailwindcss/forms
// Enable in tailwind.config.js plugins: [require('@tailwindcss/forms')]
```

```tsx
// Styled inputs out of the box
<input
  type="text"
  className="rounded-md border-slate-300 dark:border-slate-600 focus:border-blue-500 focus:ring-blue-500"
/>

// Custom radio buttons
<input
  type="radio"
  className="h-4 w-4 text-blue-600 border-slate-300 focus:ring-blue-500"
/>

// Custom checkboxes
<input
  type="checkbox"
  className="rounded text-blue-600 border-slate-300 focus:ring-blue-500"
/>
```

## Utility Patterns

### Aspect Ratio

```tsx
// Fixed aspect ratio (16:9, 4:3, etc.)
<div className="aspect-video">
  <img src="..." className="w-full h-full object-cover" />
</div>

<div className="aspect-square">
  <div className="w-full h-full flex items-center justify-center">
    Square content
  </div>
</div>
```

### Truncate Text

```tsx
// Single line truncate
<p className="truncate">
  This is a very long text that will be truncated with an ellipsis
</p>

// Multi-line truncate (line-clamp)
<p className="line-clamp-3">
  This is a very long text that will be truncated after 3 lines with an ellipsis
</p>
```

### Scroll Behavior

```tsx
// Smooth scrolling
<div className="scroll-smooth">
  <a href="#section1">Jump to section</a>
</div>

// Hide scrollbar
<div className="overflow-auto scrollbar-hide">
  Content
</div>

// Custom scrollbar (WebKit only)
<style jsx>{`
  .custom-scrollbar::-webkit-scrollbar {
    width: 8px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #94a3b8;
    border-radius: 4px;
  }
`}</style>
```

### Focus States

```tsx
// Custom focus rings
<button className="focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-slate-900">
  Button
</button>

// Focus within (parent styles when child has focus)
<div className="p-4 border border-transparent focus-within:border-blue-500">
  <input type="text" className="focus:outline-none" />
</div>
```

## shadcn/ui Integration

### Installation

```bash
npx shadcn-ui@latest init
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add dialog
```

### Customizing Components

shadcn/ui components are added to `components/ui/` and can be customized:

```tsx
// components/ui/button.tsx
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2",
  {
    variants: {
      variant: {
        default: "bg-blue-600 text-white hover:bg-blue-700",
        outline: "border border-slate-300 hover:bg-slate-100",
        ghost: "hover:bg-slate-100",
        destructive: "bg-red-600 text-white hover:bg-red-700",
      },
      size: {
        sm: "h-9 px-3 text-sm",
        md: "h-10 px-4",
        lg: "h-11 px-8 text-lg",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "md",
    },
  }
)

export function Button({ className, variant, size, ...props }) {
  return (
    <button
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}
```

## Performance Optimization

### Purge Unused Styles

Tailwind automatically purges unused styles in production based on `content` config:

```js
// tailwind.config.js
module.exports = {
  content: [
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  // Only includes classes used in these files
}
```

### JIT Mode

Just-In-Time mode generates styles on-demand (enabled by default in Tailwind 3+):

```tsx
// Arbitrary values work out of the box
<div className="w-[347px] h-[97px] bg-[#1da1f2]">
  Custom sizes and colors
</div>
```

## Common Pitfalls

### 1. Dynamic Class Names

**Don't do this**:
```tsx
// Tailwind can't detect these classes for purging
<div className={`text-${color}-500`}>Text</div>
```

**Do this instead**:
```tsx
const colorClasses = {
  red: 'text-red-500',
  blue: 'text-blue-500',
  green: 'text-green-500',
}

<div className={colorClasses[color]}>Text</div>
```

### 2. Specificity Issues

**Problem**: Custom CSS overriding Tailwind
```css
/* This will override Tailwind's bg-blue-500 */
div {
  background-color: white !important;
}
```

**Solution**: Use Tailwind utilities or @layer
```css
@layer components {
  .my-component {
    @apply bg-white;
  }
}
```

### 3. Missing Dark Mode Classes

**Don't forget dark mode variants**:
```tsx
// Light mode only - text will be invisible in dark mode
<div className="bg-slate-900 text-slate-900">Oops</div>

// Correct
<div className="bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900">
  Visible in both modes
</div>
```
