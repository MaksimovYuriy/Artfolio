import { createTheme } from '@mui/material/styles'
import { designTokens } from './designTokens'

export const theme = createTheme({
  palette: {
    mode: 'light',
    background: { default: designTokens.colors.background, paper: designTokens.colors.paper },
    text: { primary: designTokens.colors.text, secondary: designTokens.colors.textMuted },
    primary: { main: designTokens.colors.accent },
    divider: designTokens.colors.border,
  },
  typography: {
    fontFamily: designTokens.typography.body,
    h1: { fontFamily: designTokens.typography.display, fontWeight: 400, lineHeight: 0.94 },
    h2: { fontFamily: designTokens.typography.display, fontWeight: 400, lineHeight: 1 },
    h3: { fontFamily: designTokens.typography.display, fontWeight: 400 },
    button: { textTransform: 'none', letterSpacing: '0.04em' },
  },
  shape: { borderRadius: 0 },
  components: {
    MuiCssBaseline: { styleOverrides: { body: { overflowX: 'hidden' } } },
    MuiContainer: {
      defaultProps: { maxWidth: false },
      styleOverrides: { root: { maxWidth: `${designTokens.layout.maxWidth}px` } },
    },
    MuiLink: { defaultProps: { underline: 'none' } },
  },
})
