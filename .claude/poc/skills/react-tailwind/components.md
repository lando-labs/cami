# Component Patterns

Concrete React component patterns with Tailwind styling.

## Component Categories

### 1. Layout Components

**Container**:
```tsx
export function Container({ children, className }: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn("mx-auto max-w-7xl px-4 sm:px-6 lg:px-8", className)}>
      {children}
    </div>
  )
}
```

**Grid Layout**:
```tsx
export function Grid({ children, cols = 3 }: {
  children: React.ReactNode
  cols?: 1 | 2 | 3 | 4
}) {
  return (
    <div className={cn(
      "grid gap-6",
      cols === 1 && "grid-cols-1",
      cols === 2 && "grid-cols-1 md:grid-cols-2",
      cols === 3 && "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
      cols === 4 && "grid-cols-1 md:grid-cols-2 lg:grid-cols-4"
    )}>
      {children}
    </div>
  )
}
```

**Stack** (vertical spacing):
```tsx
export function Stack({ children, gap = "md" }: {
  children: React.ReactNode
  gap?: "sm" | "md" | "lg"
}) {
  return (
    <div className={cn(
      "flex flex-col",
      gap === "sm" && "gap-2",
      gap === "md" && "gap-4",
      gap === "lg" && "gap-6"
    )}>
      {children}
    </div>
  )
}
```

### 2. Form Components

**Input with Label**:
```tsx
"use client"

import { forwardRef } from "react"

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, ...props }, ref) => {
    return (
      <div className="space-y-1">
        {label && (
          <label className="block text-sm font-medium text-slate-700 dark:text-slate-300">
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={cn(
            "w-full px-3 py-2 border rounded-md",
            "bg-white dark:bg-slate-800",
            "border-slate-300 dark:border-slate-600",
            "text-slate-900 dark:text-slate-100",
            "placeholder-slate-400 dark:placeholder-slate-500",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "disabled:opacity-50 disabled:cursor-not-allowed",
            error && "border-red-500 focus:ring-red-500",
            className
          )}
          {...props}
        />
        {error && (
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        )}
      </div>
    )
  }
)
Input.displayName = "Input"
```

**Form with Validation**:
```tsx
"use client"

import { useState } from "react"
import { Input } from "./input"
import { Button } from "@/components/ui/button"

export function LoginForm() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [errors, setErrors] = useState<Record<string, string>>({})

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    const newErrors: Record<string, string> = {}

    if (!email) newErrors.email = "Email is required"
    if (!password) newErrors.password = "Password is required"

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
      return
    }

    // Submit form
    console.log({ email, password })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <Input
        label="Email"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        error={errors.email}
      />
      <Input
        label="Password"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        error={errors.password}
      />
      <Button type="submit" className="w-full">
        Sign In
      </Button>
    </form>
  )
}
```

### 3. Data Display Components

**Card**:
```tsx
export function Card({ children, className }: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn(
      "bg-white dark:bg-slate-800",
      "border border-slate-200 dark:border-slate-700",
      "rounded-lg shadow-sm",
      className
    )}>
      {children}
    </div>
  )
}

export function CardHeader({ children, className }: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn("px-6 py-4 border-b border-slate-200 dark:border-slate-700", className)}>
      {children}
    </div>
  )
}

export function CardContent({ children, className }: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn("px-6 py-4", className)}>
      {children}
    </div>
  )
}

export function CardFooter({ children, className }: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn("px-6 py-4 border-t border-slate-200 dark:border-slate-700", className)}>
      {children}
    </div>
  )
}

// Usage:
<Card>
  <CardHeader>
    <h3 className="text-lg font-semibold">Card Title</h3>
  </CardHeader>
  <CardContent>
    <p>Card content goes here.</p>
  </CardContent>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>
```

**Table**:
```tsx
export function Table({ children }: { children: React.ReactNode }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse">
        {children}
      </table>
    </div>
  )
}

export function TableHeader({ children }: { children: React.ReactNode }) {
  return (
    <thead className="bg-slate-50 dark:bg-slate-800">
      {children}
    </thead>
  )
}

export function TableRow({ children }: { children: React.ReactNode }) {
  return (
    <tr className="border-b border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800/50">
      {children}
    </tr>
  )
}

export function TableHead({ children }: { children: React.ReactNode }) {
  return (
    <th className="px-4 py-3 text-left text-sm font-semibold text-slate-700 dark:text-slate-300">
      {children}
    </th>
  )
}

export function TableCell({ children }: { children: React.ReactNode }) {
  return (
    <td className="px-4 py-3 text-sm text-slate-900 dark:text-slate-100">
      {children}
    </td>
  )
}

// Usage:
<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Email</TableHead>
    </TableRow>
  </TableHeader>
  <tbody>
    {users.map(user => (
      <TableRow key={user.id}>
        <TableCell>{user.name}</TableCell>
        <TableCell>{user.email}</TableCell>
      </TableRow>
    ))}
  </tbody>
</Table>
```

### 4. Feedback Components

**Loading Spinner**:
```tsx
export function Spinner({ size = "md" }: {
  size?: "sm" | "md" | "lg"
}) {
  return (
    <div
      className={cn(
        "animate-spin rounded-full border-2 border-slate-300 border-t-blue-600",
        size === "sm" && "h-4 w-4",
        size === "md" && "h-8 w-8",
        size === "lg" && "h-12 w-12"
      )}
    />
  )
}
```

**Alert**:
```tsx
export function Alert({ children, variant = "info" }: {
  children: React.ReactNode
  variant?: "info" | "success" | "warning" | "error"
}) {
  return (
    <div className={cn(
      "px-4 py-3 rounded-lg border",
      variant === "info" && "bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800 text-blue-900 dark:text-blue-100",
      variant === "success" && "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800 text-green-900 dark:text-green-100",
      variant === "warning" && "bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800 text-yellow-900 dark:text-yellow-100",
      variant === "error" && "bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-900 dark:text-red-100"
    )}>
      {children}
    </div>
  )
}
```

**Toast Notification**:
```tsx
"use client"

import { useEffect, useState } from "react"

interface ToastProps {
  message: string
  variant?: "info" | "success" | "error"
  duration?: number
  onClose: () => void
}

export function Toast({ message, variant = "info", duration = 3000, onClose }: ToastProps) {
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false)
      setTimeout(onClose, 300) // Wait for fade out
    }, duration)

    return () => clearTimeout(timer)
  }, [duration, onClose])

  return (
    <div
      className={cn(
        "fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg border transition-opacity",
        isVisible ? "opacity-100" : "opacity-0",
        variant === "info" && "bg-blue-50 dark:bg-blue-900 border-blue-200 dark:border-blue-700 text-blue-900 dark:text-blue-100",
        variant === "success" && "bg-green-50 dark:bg-green-900 border-green-200 dark:border-green-700 text-green-900 dark:text-green-100",
        variant === "error" && "bg-red-50 dark:bg-red-900 border-red-200 dark:border-red-700 text-red-900 dark:text-red-100"
      )}
    >
      <div className="flex items-center gap-3">
        <p className="text-sm font-medium">{message}</p>
        <button
          onClick={() => {
            setIsVisible(false)
            setTimeout(onClose, 300)
          }}
          className="text-lg leading-none opacity-70 hover:opacity-100"
        >
          ×
        </button>
      </div>
    </div>
  )
}
```

### 5. Navigation Components

**Tabs**:
```tsx
"use client"

import { useState } from "react"

interface Tab {
  id: string
  label: string
  content: React.ReactNode
}

export function Tabs({ tabs }: { tabs: Tab[] }) {
  const [activeTab, setActiveTab] = useState(tabs[0].id)

  return (
    <div>
      <div className="border-b border-slate-200 dark:border-slate-700">
        <nav className="flex gap-4">
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                "px-4 py-2 -mb-px border-b-2 font-medium text-sm transition-colors",
                activeTab === tab.id
                  ? "border-blue-500 text-blue-600 dark:text-blue-400"
                  : "border-transparent text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100"
              )}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>
      <div className="py-4">
        {tabs.find(tab => tab.id === activeTab)?.content}
      </div>
    </div>
  )
}
```

**Dropdown Menu** (shadcn/ui pattern):
```tsx
"use client"

import { useState, useRef, useEffect } from "react"

export function DropdownMenu({ trigger, children }: {
  trigger: React.ReactNode
  children: React.ReactNode
}) {
  const [isOpen, setIsOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  return (
    <div ref={ref} className="relative">
      <div onClick={() => setIsOpen(!isOpen)}>
        {trigger}
      </div>
      {isOpen && (
        <div className="absolute right-0 mt-2 w-56 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg shadow-lg z-50">
          {children}
        </div>
      )}
    </div>
  )
}

export function DropdownMenuItem({ children, onClick }: {
  children: React.ReactNode
  onClick?: () => void
}) {
  return (
    <button
      onClick={onClick}
      className="w-full px-4 py-2 text-left text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 first:rounded-t-lg last:rounded-b-lg"
    >
      {children}
    </button>
  )
}
```

## Composition Patterns

### Compound Components

Pattern where components are designed to work together:

```tsx
// Parent manages state, children consume it
export function Accordion({ children }: { children: React.ReactNode }) {
  const [openItems, setOpenItems] = useState<Set<string>>(new Set())

  return (
    <div className="space-y-2">
      {React.Children.map(children, (child) => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child, { openItems, setOpenItems })
        }
        return child
      })}
    </div>
  )
}

export function AccordionItem({ id, title, children, openItems, setOpenItems }: {
  id: string
  title: string
  children: React.ReactNode
  openItems?: Set<string>
  setOpenItems?: (fn: (prev: Set<string>) => Set<string>) => void
}) {
  const isOpen = openItems?.has(id)

  return (
    <div className="border border-slate-200 dark:border-slate-700 rounded-lg">
      <button
        onClick={() => setOpenItems?.(prev => {
          const next = new Set(prev)
          isOpen ? next.delete(id) : next.add(id)
          return next
        })}
        className="w-full px-4 py-3 text-left font-medium flex justify-between items-center"
      >
        {title}
        <span>{isOpen ? "−" : "+"}</span>
      </button>
      {isOpen && (
        <div className="px-4 py-3 border-t border-slate-200 dark:border-slate-700">
          {children}
        </div>
      )}
    </div>
  )
}
```

### Render Props

Pattern for sharing logic between components:

```tsx
export function MouseTracker({ render }: {
  render: (position: { x: number; y: number }) => React.ReactNode
}) {
  const [position, setPosition] = useState({ x: 0, y: 0 })

  return (
    <div
      onMouseMove={(e) => setPosition({ x: e.clientX, y: e.clientY })}
      className="w-full h-64 bg-slate-100 dark:bg-slate-800 rounded-lg"
    >
      {render(position)}
    </div>
  )
}

// Usage:
<MouseTracker
  render={({ x, y }) => (
    <div className="p-4">
      Mouse at: {x}, {y}
    </div>
  )}
/>
```

## Performance Patterns

### Memoization

```tsx
import { memo, useMemo, useCallback } from "react"

// Memo: Prevent re-renders if props haven't changed
export const ExpensiveComponent = memo(function ExpensiveComponent({ data }: {
  data: Array<any>
}) {
  return (
    <ul>
      {data.map(item => <li key={item.id}>{item.name}</li>)}
    </ul>
  )
})

// useMemo: Cache expensive calculations
function DataTable({ items, filter }: { items: any[], filter: string }) {
  const filteredItems = useMemo(
    () => items.filter(item => item.name.includes(filter)),
    [items, filter]
  )

  return <div>{/* Render filteredItems */}</div>
}

// useCallback: Stable function references
function Parent() {
  const [count, setCount] = useState(0)

  const handleClick = useCallback(() => {
    setCount(c => c + 1)
  }, []) // Stable across re-renders

  return <Child onClick={handleClick} />
}
```
