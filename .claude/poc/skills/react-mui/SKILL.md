---
name: react-mui
description: React 19+ with Material UI (MUI) implementation patterns. Provides component structure, hooks patterns, MUI theming, sx prop conventions, and Material Design principles. Use when building React frontends with MUI styling.
---

# React + MUI Implementation Skill

This skill provides concrete implementation patterns for React 19+ with Material UI (MUI v6+). Use this when you need specific React/MUI syntax, component patterns, or Material Design conventions.

## When This Skill Activates

- Building React components with MUI styling
- Implementing Material Design patterns
- Setting up MUI theming and customization
- Using MUI component library
- Implementing responsive Material layouts
- Working with sx prop and styled components

## Quick Reference

### Modern React Patterns

```tsx
// Modern function component with TypeScript
import { FC, useState } from 'react';
import { Button, Stack } from '@mui/material';

interface Props {
  initialCount?: number;
}

export const Counter: FC<Props> = ({ initialCount = 0 }) => {
  const [count, setCount] = useState(initialCount);

  return (
    <Stack spacing={2} alignItems="center">
      <Button
        variant="contained"
        onClick={() => setCount(c => c + 1)}
        sx={{ minWidth: 120 }}
      >
        Count: {count}
      </Button>
    </Stack>
  );
};
```

### MUI Styling Approaches

```tsx
// 1. Inline sx prop (preferred for one-offs)
<Box sx={{
  p: 2,
  bgcolor: 'primary.main',
  borderRadius: 1,
  '&:hover': { bgcolor: 'primary.dark' }
}} />

// 2. Styled components (for reusable patterns)
import { styled } from '@mui/material/styles';
import { Paper } from '@mui/material';

const StyledCard = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
  backgroundColor: theme.palette.background.paper,
  borderRadius: theme.shape.borderRadius,
  '&:hover': {
    boxShadow: theme.shadows[4],
  },
}));

// 3. Theme-based variants
<Button variant="contained" color="primary" size="large" />
```

### Material Design Spacing

```tsx
// MUI uses 8px base spacing unit
<Box sx={{
  m: 2,      // margin: 16px (2 * 8px)
  mt: 3,     // margin-top: 24px
  px: 1.5,   // padding-left/right: 12px
  py: 4,     // padding-top/bottom: 32px
}} />

// Stack for automatic spacing
<Stack spacing={2} direction="row" alignItems="center">
  <Item />
  <Item />
  <Item />
</Stack>
```

### Responsive Breakpoints

```tsx
// MUI breakpoints: xs, sm, md, lg, xl
<Box sx={{
  width: { xs: '100%', sm: '80%', md: '60%', lg: '50%' },
  p: { xs: 1, sm: 2, md: 3 },
  display: { xs: 'block', md: 'flex' }
}} />

// useMediaQuery hook
import { useMediaQuery, useTheme } from '@mui/material';

function ResponsiveComponent() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  return isMobile ? <MobileView /> : <DesktopView />;
}
```

### Common MUI Components

```tsx
import {
  AppBar, Toolbar, Typography, Button,
  Card, CardContent, CardActions,
  TextField, Select, MenuItem,
  Dialog, DialogTitle, DialogContent, DialogActions,
  Snackbar, Alert,
  Grid2, Stack, Box, Container,
  IconButton, Tooltip,
  CircularProgress, LinearProgress,
} from '@mui/material';

// Layout example
<Container maxWidth="lg">
  <Grid2 container spacing={3}>
    <Grid2 size={{ xs: 12, md: 6 }}>
      <Card>
        <CardContent>
          <Typography variant="h5" gutterBottom>
            Card Title
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Card content here
          </Typography>
        </CardContent>
        <CardActions>
          <Button size="small">Learn More</Button>
        </CardActions>
      </Card>
    </Grid2>
  </Grid2>
</Container>
```

### Form Handling

```tsx
import { TextField, Button, Stack } from '@mui/material';
import { useState } from 'react';

function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Validation logic
    if (!email) setErrors(prev => ({ ...prev, email: 'Required' }));
    // Submit logic
  };

  return (
    <form onSubmit={handleSubmit}>
      <Stack spacing={2}>
        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={!!errors.email}
          helperText={errors.email}
          fullWidth
          required
        />
        <TextField
          label="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={!!errors.password}
          helperText={errors.password}
          fullWidth
          required
        />
        <Button
          type="submit"
          variant="contained"
          fullWidth
          size="large"
        >
          Sign In
        </Button>
      </Stack>
    </form>
  );
}
```

### State Management Patterns

```tsx
// useState for local state
const [open, setOpen] = useState(false);

// useReducer for complex state
import { useReducer } from 'react';

type State = { count: number; step: number };
type Action =
  | { type: 'increment' }
  | { type: 'decrement' }
  | { type: 'setStep'; step: number };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'increment': return { ...state, count: state.count + state.step };
    case 'decrement': return { ...state, count: state.count - state.step };
    case 'setStep': return { ...state, step: action.step };
  }
}

const [state, dispatch] = useReducer(reducer, { count: 0, step: 1 });

// Context for shared state
import { createContext, useContext, ReactNode } from 'react';

const ThemeContext = createContext<{ mode: 'light' | 'dark' } | null>(null);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setMode] = useState<'light' | 'dark'>('light');
  return (
    <ThemeContext.Provider value={{ mode }}>
      {children}
    </ThemeContext.Provider>
  );
}

export const useThemeMode = () => {
  const context = useContext(ThemeContext);
  if (!context) throw new Error('useThemeMode must be within ThemeProvider');
  return context;
};
```

### Data Fetching

```tsx
import { useState, useEffect } from 'react';
import { CircularProgress, Alert } from '@mui/material';

function UserProfile({ userId }: { userId: string }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchUser() {
      try {
        setLoading(true);
        const response = await fetch(`/api/users/${userId}`);
        if (!response.ok) throw new Error('Failed to fetch');
        const data = await response.json();
        if (!cancelled) setUser(data);
      } catch (err) {
        if (!cancelled) setError(err.message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    fetchUser();
    return () => { cancelled = true; };
  }, [userId]);

  if (loading) return <CircularProgress />;
  if (error) return <Alert severity="error">{error}</Alert>;
  if (!user) return null;

  return <UserCard user={user} />;
}
```

## Project Structure

```
src/
├── components/
│   ├── common/           # Shared components
│   │   ├── Button/
│   │   ├── Card/
│   │   └── Layout/
│   ├── features/         # Feature-specific components
│   │   ├── auth/
│   │   └── dashboard/
│   └── index.ts          # Barrel exports
├── hooks/                # Custom hooks
│   ├── useAuth.ts
│   └── useApi.ts
├── theme/                # MUI theme configuration
│   ├── theme.ts
│   └── palette.ts
├── types/                # TypeScript types
│   └── index.ts
├── utils/                # Helper functions
│   └── format.ts
├── App.tsx
└── main.tsx
```

## TypeScript Integration

```tsx
// Component props
interface ButtonProps {
  variant?: 'text' | 'outlined' | 'contained';
  color?: 'primary' | 'secondary' | 'error';
  onClick?: () => void;
  disabled?: boolean;
  children: ReactNode;
}

// Event handlers
const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
  console.log(event.currentTarget);
};

const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
  setValue(event.target.value);
};

// Extending MUI types
import { ButtonProps as MuiButtonProps } from '@mui/material';

interface CustomButtonProps extends MuiButtonProps {
  isLoading?: boolean;
}
```

## Performance Patterns

```tsx
// React.memo for expensive components
import { memo } from 'react';

export const ExpensiveList = memo(({ items }: { items: Item[] }) => {
  return (
    <List>
      {items.map(item => (
        <ListItem key={item.id}>{item.name}</ListItem>
      ))}
    </List>
  );
});

// useCallback for stable references
const handleSubmit = useCallback((data: FormData) => {
  // Submit logic
}, [/* dependencies */]);

// useMemo for expensive computations
const sortedItems = useMemo(() => {
  return items.sort((a, b) => a.name.localeCompare(b.name));
}, [items]);

// Lazy loading
import { lazy, Suspense } from 'react';
import { CircularProgress } from '@mui/material';

const Dashboard = lazy(() => import('./Dashboard'));

function App() {
  return (
    <Suspense fallback={<CircularProgress />}>
      <Dashboard />
    </Suspense>
  );
}
```

## Additional Resources

- See `components.md` for common component patterns and examples
- See `theming.md` for MUI theme setup and customization
- [MUI Documentation](https://mui.com/material-ui/getting-started/)
- [React Documentation](https://react.dev/)

## Boundaries

**This skill provides:**
- React component patterns and syntax
- MUI component usage and styling
- Material Design conventions
- TypeScript integration
- Performance optimization patterns

**This skill does NOT provide:**
- Architecture decisions (agent handles this)
- Testing strategies (use QA agent)
- Backend integration (use backend agent)
- Build configuration (agent decides)
