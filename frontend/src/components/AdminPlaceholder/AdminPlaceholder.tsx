import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import { Box, Container, Link, Stack, Typography } from '@mui/material'

export function AdminPlaceholder() {
  return (
    <Box component="main" sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center', py: 6 }}>
      <Container>
        <Stack sx={{ maxWidth: 620, mx: 'auto', alignItems: 'flex-start' }}>
          <Typography variant="overline" color="text.secondary" sx={{ letterSpacing: '.18em' }}>
            Artfolio · Для автора
          </Typography>
          <Typography variant="h1" sx={{ mt: 3, fontSize: { xs: '3.5rem', md: '5.5rem' } }}>
            Управление портфолио
          </Typography>
          <Typography color="text.secondary" sx={{ mt: 4, maxWidth: 520, fontSize: '1.1rem', lineHeight: 1.7 }}>
            Здесь позднее появится защищённый вход для добавления и редактирования работ. На этапе прототипа административные функции отключены.
          </Typography>
          <Link href="/" color="text.primary" sx={{ mt: 6, display: 'inline-flex', alignItems: 'center', gap: 1 }}>
            <ArrowBackIcon fontSize="small" /> Вернуться к портфолио
          </Link>
        </Stack>
      </Container>
    </Box>
  )
}
