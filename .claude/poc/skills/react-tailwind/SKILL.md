---
name: react-tailwind
description: React 19+ with Tailwind CSS implementation patterns. Provides component structure, hooks patterns, Tailwind utilities, and shadcn/ui conventions. Use when building React frontends with Tailwind styling.
---

# React + Tailwind Implementation Skill

Concrete implementation patterns for React 19+ with Tailwind CSS. This skill provides the technical "HOW" that complements methodology-focused frontend agents.

## When This Skill Activates

- Building React components with Tailwind CSS
- Implementing responsive layouts and dark mode
- Using shadcn/ui components
- Setting up React hooks and state management
- Styling with utility classes

## Core Stack

**React 19+**:
- Server Components by default
- "use client" for interactivity
- Hooks: useState, useEffect, useCallback, useMemo, useRef
- Custom hooks for reusable logic

**Tailwind CSS 3.4+**:
- Utility-first styling
- Responsive design with breakpoints
- Dark mode with `class` strategy
- Custom configuration via `tailwind.config.js`

**shadcn/ui** (optional):
- Accessible components built on Radix UI
- Copy-paste component pattern
- Customizable via `components/ui/`

## Quick Reference

### File Structure

```
src/
├── app/                      # Next.js App Router (if applicable)
│   ├── layout.tsx
│   ├── page.tsx
│   └── api/
├── components/
│   ├── ui/                   # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   └── dialog.tsx
│   ├── layout/               # Layout components
│   │   ├── header.tsx
│   │   ├── footer.tsx
│   │   └── sidebar.tsx
│   └── features/             # Feature-specific components
│       ├── auth/
│       └── dashboard/
├── hooks/                    # Custom React hooks
│   ├── use-toast.ts
│   └── use-media-query.ts
├── lib/
│   └── utils.ts              # cn() utility for class merging
└── styles/
    └── globals.css           # Tailwind directives
```

### Tailwind Class Organization

**Order of utility classes** (for consistency):
1. Layout (flex, grid, block)
2. Positioning (relative, absolute, fixed)
3. Spacing (p-, m-, gap-)
4. Sizing (w-, h-, min-, max-)
5. Typography (text-, font-)
6. Colors (bg-, text-, border-)
7. Borders (border, rounded)
8. Effects (shadow, opacity)
9. States (hover:, focus:, active:)
10. Responsive (sm:, md:, lg:, xl:, 2xl:)
11. Dark mode (dark:)

**Example**:
```tsx
<div className="flex flex-col gap-4 p-6 w-full max-w-md bg-white dark:bg-slate-900 rounded-lg shadow-lg hover:shadow-xl transition-shadow">
```

### Component Patterns

**Client Component** (interactive):
```tsx
"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"

export function Counter() {
  const [count, setCount] = useState(0)

  return (
    <div className="flex flex-col items-center gap-4">
      <p className="text-2xl font-bold">{count}</p>
      <Button onClick={() => setCount(c => c + 1)}>
        Increment
      </Button>
    </div>
  )
}
```

**Server Component** (default in Next.js 13+):
```tsx
// No "use client" directive
// Can be async
// No useState, useEffect, or event handlers

export async function UserList() {
  const users = await fetchUsers() // Can fetch directly

  return (
    <ul className="space-y-2">
      {users.map(user => (
        <li key={user.id} className="p-4 bg-slate-100 dark:bg-slate-800 rounded">
          {user.name}
        </li>
      ))}
    </ul>
  )
}
```

### Responsive Design

**Tailwind Breakpoints**:
- `sm:` - 640px and up
- `md:` - 768px and up
- `lg:` - 1024px and up
- `xl:` - 1280px and up
- `2xl:` - 1536px and up

**Mobile-first approach**:
```tsx
<div className="w-full md:w-1/2 lg:w-1/3">
  {/* Full width on mobile, half on tablet, third on desktop */}
</div>

<nav className="hidden lg:flex lg:gap-6">
  {/* Hidden on mobile/tablet, shown as flex on desktop */}
</nav>
```

### Dark Mode

**Setup** (`tailwind.config.js`):
```js
module.exports = {
  darkMode: 'class', // Enable class-based dark mode
  // ...
}
```

**Usage**:
```tsx
<div className="bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100">
  {/* Light bg + dark text in light mode */}
  {/* Dark bg + light text in dark mode */}
</div>
```

**Theme Toggle Pattern**:
```tsx
"use client"

import { useEffect, useState } from "react"

export function ThemeToggle() {
  const [theme, setTheme] = useState<"light" | "dark">("light")

  useEffect(() => {
    // Check system preference or localStorage
    const stored = localStorage.getItem("theme")
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches
    setTheme(stored as "light" | "dark" || (prefersDark ? "dark" : "light"))
  }, [])

  useEffect(() => {
    // Apply theme to document
    document.documentElement.classList.toggle("dark", theme === "dark")
    localStorage.setItem("theme", theme)
  }, [theme])

  return (
    <button
      onClick={() => setTheme(t => t === "light" ? "dark" : "light")}
      className="p-2 rounded hover:bg-slate-200 dark:hover:bg-slate-700"
    >
      {theme === "light" ? "🌙" : "☀️"}
    </button>
  )
}
```

## Common Utilities

### Class Name Merger (cn)

**Setup** (`lib/utils.ts`):
```ts
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

**Usage**:
```tsx
import { cn } from "@/lib/utils"

<button
  className={cn(
    "px-4 py-2 rounded",
    isActive && "bg-blue-500 text-white",
    isDisabled && "opacity-50 cursor-not-allowed"
  )}
>
  Click me
</button>
```

### Custom Hooks

**useMediaQuery** (responsive hooks):
```ts
"use client"

import { useEffect, useState } from "react"

export function useMediaQuery(query: string) {
  const [matches, setMatches] = useState(false)

  useEffect(() => {
    const media = window.matchMedia(query)
    setMatches(media.matches)

    const listener = (e: MediaQueryListEvent) => setMatches(e.matches)
    media.addEventListener("change", listener)
    return () => media.removeEventListener("change", listener)
  }, [query])

  return matches
}

// Usage:
const isMobile = useMediaQuery("(max-width: 768px)")
```

## References

- See `components.md` for detailed component patterns and examples
- See `styling.md` for Tailwind conventions and advanced patterns

## Limitations

This skill provides implementation patterns for:
- React 19+ (Server Components, hooks)
- Tailwind CSS 3.4+
- shadcn/ui components

**Not covered**:
- State management libraries (Redux, Zustand, Jotai)
- Form libraries (React Hook Form, Formik)
- Animation libraries (Framer Motion)
- Testing patterns (see testing-focused skills)

For methodology and high-level decisions, defer to the agent loading this skill.
