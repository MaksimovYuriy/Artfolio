import { useState } from 'react'
import MenuIcon from '@mui/icons-material/Menu'
import CloseIcon from '@mui/icons-material/Close'
import { AppBar, Container, IconButton, Link, Stack, Toolbar, Typography } from '@mui/material'

const navigation = [
  { label: 'Работы', href: '#works' },
  { label: 'О художнице', href: '#about' },
  { label: 'Контакты', href: '#contacts' },
]

export function Header({ brand }: { brand: string }) {
  const [open, setOpen] = useState(false)

  return (
    <AppBar position="sticky" elevation={0} color="transparent" sx={{ bgcolor: 'rgba(243,240,234,.92)', backdropFilter: 'blur(12px)', borderBottom: 1, borderColor: 'divider' }}>
      <Container>
        <Toolbar disableGutters sx={{ minHeight: { xs: 68, md: 80 }, justifyContent: 'space-between' }}>
          <Typography component="a" href="#top" variant="h6" color="text.primary" sx={{ fontFamily: 'Georgia, serif', textDecoration: 'none', letterSpacing: '.02em' }}>
            {brand}
          </Typography>
          <Stack component="nav" direction="row" spacing={4} sx={{ display: { xs: 'none', sm: 'flex' } }} aria-label="Основная навигация">
            {navigation.map((item) => <Link key={item.href} href={item.href} color="text.primary" sx={{ fontSize: 14 }}>{item.label}</Link>)}
          </Stack>
          <IconButton onClick={() => setOpen((value) => !value)} sx={{ display: { sm: 'none' } }} aria-label={open ? 'Закрыть меню' : 'Открыть меню'}>
            {open ? <CloseIcon /> : <MenuIcon />}
          </IconButton>
        </Toolbar>
        {open && (
          <Stack component="nav" spacing={2.5} sx={{ display: { sm: 'none' }, pb: 3 }}>
            {navigation.map((item) => <Link key={item.href} href={item.href} color="text.primary" onClick={() => setOpen(false)}>{item.label}</Link>)}
          </Stack>
        )}
      </Container>
    </AppBar>
  )
}
