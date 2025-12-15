# MUI Theming and Customization

Material UI theme setup, customization patterns, and design system implementation.

## Theme Setup

### Basic Theme Configuration

```tsx
// src/theme/theme.ts
import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
      light: '#42a5f5',
      dark: '#1565c0',
      contrastText: '#fff',
    },
    secondary: {
      main: '#9c27b0',
      light: '#ba68c8',
      dark: '#7b1fa2',
      contrastText: '#fff',
    },
    error: {
      main: '#d32f2f',
    },
    warning: {
      main: '#ed6c02',
    },
    info: {
      main: '#0288d1',
    },
    success: {
      main: '#2e7d32',
    },
    background: {
      default: '#fafafa',
      paper: '#fff',
    },
    text: {
      primary: 'rgba(0, 0, 0, 0.87)',
      secondary: 'rgba(0, 0, 0, 0.6)',
      disabled: 'rgba(0, 0, 0, 0.38)',
    },
  },
  typography: {
    fontFamily: [
      'Roboto',
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontSize: '3rem',
      fontWeight: 300,
      lineHeight: 1.167,
    },
    h2: {
      fontSize: '2.5rem',
      fontWeight: 300,
      lineHeight: 1.2,
    },
    h3: {
      fontSize: '2rem',
      fontWeight: 400,
      lineHeight: 1.167,
    },
    h4: {
      fontSize: '1.75rem',
      fontWeight: 400,
      lineHeight: 1.235,
    },
    h5: {
      fontSize: '1.5rem',
      fontWeight: 400,
      lineHeight: 1.334,
    },
    h6: {
      fontSize: '1.25rem',
      fontWeight: 500,
      lineHeight: 1.6,
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.5,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.43,
    },
    button: {
      textTransform: 'none', // Disable uppercase buttons
      fontWeight: 500,
    },
  },
  spacing: 8, // Base spacing unit (8px)
  shape: {
    borderRadius: 4,
  },
  shadows: [
    'none',
    '0px 2px 4px rgba(0,0,0,0.1)',
    '0px 4px 8px rgba(0,0,0,0.12)',
    // ... 25 elevation levels total
  ],
});
```

### Applying Theme

```tsx
// src/App.tsx
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { theme } from './theme/theme';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline /> {/* Normalize styles */}
      <YourApp />
    </ThemeProvider>
  );
}
```

## Dark Mode

### Toggle Theme Mode

```tsx
// src/theme/ThemeProvider.tsx
import { createTheme, ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { createContext, useContext, useMemo, useState, ReactNode } from 'react';
import CssBaseline from '@mui/material/CssBaseline';

type ThemeMode = 'light' | 'dark';

interface ThemeContextType {
  mode: ThemeMode;
  toggleMode: () => void;
}

const ThemeContext = createContext<ThemeContextType | null>(null);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setMode] = useState<ThemeMode>('light');

  const toggleMode = () => {
    setMode((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode,
          ...(mode === 'light'
            ? {
                // Light mode palette
                primary: { main: '#1976d2' },
                background: { default: '#fafafa', paper: '#fff' },
              }
            : {
                // Dark mode palette
                primary: { main: '#90caf9' },
                background: { default: '#121212', paper: '#1e1e1e' },
              }),
        },
      }),
    [mode]
  );

  return (
    <ThemeContext.Provider value={{ mode, toggleMode }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
}

export const useThemeMode = () => {
  const context = useContext(ThemeContext);
  if (!context) throw new Error('useThemeMode must be within ThemeProvider');
  return context;
};

// Usage in component
function ThemeToggle() {
  const { mode, toggleMode } = useThemeMode();

  return (
    <IconButton onClick={toggleMode} color="inherit">
      {mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
    </IconButton>
  );
}
```

## Component Customization

### Global Component Overrides

```tsx
import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 500,
        },
        contained: {
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0px 2px 4px rgba(0,0,0,0.2)',
          },
        },
      },
      defaultProps: {
        disableElevation: true,
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0px 2px 8px rgba(0,0,0,0.08)',
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        variant: 'outlined',
      },
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 8,
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 6,
          fontWeight: 500,
        },
      },
    },
  },
});
```

### Custom Variants

```tsx
import { createTheme } from '@mui/material/styles';

declare module '@mui/material/Button' {
  interface ButtonPropsVariantOverrides {
    gradient: true;
  }
}

export const theme = createTheme({
  components: {
    MuiButton: {
      variants: [
        {
          props: { variant: 'gradient' },
          style: {
            background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
            color: 'white',
            boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
            '&:hover': {
              background: 'linear-gradient(45deg, #FE5B7B 30%, #FF7E43 90%)',
            },
          },
        },
      ],
    },
  },
});

// Usage
<Button variant="gradient">Custom Button</Button>
```

## Responsive Design

### Responsive Values

```tsx
import { Box, Typography } from '@mui/material';
import { useTheme } from '@mui/material/styles';

function ResponsiveComponent() {
  const theme = useTheme();

  return (
    <Box
      sx={{
        // Responsive spacing
        p: { xs: 2, sm: 3, md: 4 },
        // Responsive width
        width: { xs: '100%', sm: '80%', md: '60%', lg: '50%' },
        // Responsive display
        display: { xs: 'block', md: 'flex' },
        // Responsive gap
        gap: { xs: 1, sm: 2, md: 3 },
        // Conditional styles based on breakpoint
        [theme.breakpoints.down('sm')]: {
          backgroundColor: 'primary.light',
        },
        [theme.breakpoints.up('md')]: {
          backgroundColor: 'primary.main',
        },
      }}
    >
      <Typography
        variant="h1"
        sx={{
          fontSize: { xs: '2rem', sm: '2.5rem', md: '3rem' },
        }}
      >
        Responsive Text
      </Typography>
    </Box>
  );
}
```

### Breakpoint Hooks

```tsx
import { useMediaQuery, useTheme } from '@mui/material';

function AdaptiveLayout() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const isTablet = useMediaQuery(theme.breakpoints.between('sm', 'md'));
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));

  if (isMobile) return <MobileLayout />;
  if (isTablet) return <TabletLayout />;
  return <DesktopLayout />;
}
```

## Design Tokens

### Spacing System

```tsx
// MUI uses 8px base spacing unit
theme.spacing(1);  // 8px
theme.spacing(2);  // 16px
theme.spacing(3);  // 24px
theme.spacing(0.5); // 4px

// In sx prop
<Box sx={{
  m: 2,           // margin: 16px
  mt: 3,          // margin-top: 24px
  px: 1.5,        // padding-left/right: 12px
  gap: 4,         // gap: 32px
}} />
```

### Color System

```tsx
// Palette colors
<Box sx={{
  bgcolor: 'primary.main',
  color: 'primary.contrastText',
  borderColor: 'primary.dark',
}} />

// Shades (50-900)
<Box sx={{
  bgcolor: 'grey.100',
  borderColor: 'grey.300',
}} />

// Semantic colors
<Box sx={{
  bgcolor: 'error.main',
  bgcolor: 'warning.main',
  bgcolor: 'info.main',
  bgcolor: 'success.main',
}} />

// Text colors
<Typography color="text.primary" />
<Typography color="text.secondary" />
<Typography color="text.disabled" />
```

### Typography System

```tsx
// Variant-based
<Typography variant="h1">Heading 1</Typography>
<Typography variant="h2">Heading 2</Typography>
<Typography variant="body1">Body text</Typography>
<Typography variant="body2">Secondary text</Typography>
<Typography variant="caption">Caption text</Typography>
<Typography variant="button">Button text</Typography>

// Custom typography
<Typography sx={{
  fontSize: theme.typography.h4.fontSize,
  fontWeight: theme.typography.fontWeightBold,
  lineHeight: theme.typography.body1.lineHeight,
}} />
```

### Elevation System

```tsx
// 0-24 elevation levels
<Paper elevation={0} /> // No shadow
<Paper elevation={1} /> // Subtle shadow
<Paper elevation={3} /> // Medium shadow
<Paper elevation={8} /> // Heavy shadow

// Custom shadows
<Box sx={{
  boxShadow: theme.shadows[2],
  boxShadow: (theme) => theme.shadows[4],
}} />
```

## Advanced Patterns

### Theme Augmentation

```tsx
// Extend theme with custom values
declare module '@mui/material/styles' {
  interface Theme {
    status: {
      danger: string;
    };
  }
  interface ThemeOptions {
    status?: {
      danger?: string;
    };
  }
}

export const theme = createTheme({
  status: {
    danger: '#e53e3e',
  },
});

// Usage
<Box sx={{ bgcolor: (theme) => theme.status.danger }} />
```

### Nested Themes

```tsx
import { ThemeProvider, createTheme } from '@mui/material/styles';

function AdminSection() {
  const adminTheme = createTheme({
    palette: {
      primary: { main: '#9c27b0' },
    },
  });

  return (
    <ThemeProvider theme={adminTheme}>
      <AdminDashboard />
    </ThemeProvider>
  );
}
```

### Accessing Theme in Components

```tsx
import { useTheme } from '@mui/material/styles';

function ThemedComponent() {
  const theme = useTheme();

  return (
    <Box
      sx={{
        padding: theme.spacing(2),
        backgroundColor: theme.palette.primary.main,
        borderRadius: theme.shape.borderRadius,
        [theme.breakpoints.down('sm')]: {
          padding: theme.spacing(1),
        },
      }}
    />
  );
}
```

## Best Practices

1. **Use sx prop for one-off styles**
   ```tsx
   <Box sx={{ p: 2, bgcolor: 'primary.main' }} />
   ```

2. **Use styled() for reusable components**
   ```tsx
   const StyledCard = styled(Card)(({ theme }) => ({
     padding: theme.spacing(2),
   }));
   ```

3. **Use theme tokens, not hardcoded values**
   ```tsx
   // Good
   <Box sx={{ color: 'primary.main', p: 2 }} />

   // Avoid
   <Box sx={{ color: '#1976d2', padding: '16px' }} />
   ```

4. **Leverage responsive values**
   ```tsx
   <Box sx={{ width: { xs: '100%', md: '50%' } }} />
   ```

5. **Use theme breakpoints for media queries**
   ```tsx
   sx={{
     [theme.breakpoints.down('sm')]: { display: 'none' }
   }}
   ```
