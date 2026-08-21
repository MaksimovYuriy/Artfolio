import { Box, Container, Link, Stack, Typography } from '@mui/material'
import type { ArtistProfile } from '../../types/artist'

export function Footer({ artist }: { artist: ArtistProfile }) {
  return (
    <Box component="footer" sx={{ borderTop: 1, borderColor: 'divider', py: 4 }}>
      <Container>
        <Stack direction={{ xs: 'column', sm: 'row' }} sx={{ gap: 1, justifyContent: 'space-between' }}>
          <Typography variant="body2">{artist.brandName} · {artist.name}</Typography>
          <Stack direction="row" spacing={3}>
            <Link href="/admin" variant="body2" color="text.secondary">Вход для автора</Link>
            <Typography variant="body2" color="text.secondary">© {new Date().getFullYear()}</Typography>
          </Stack>
        </Stack>
      </Container>
    </Box>
  )
}
