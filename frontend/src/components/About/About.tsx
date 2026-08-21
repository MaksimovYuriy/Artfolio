import { Box, Container, Stack, Typography } from '@mui/material'
import type { ArtistProfile } from '../../types/artist'
import { SectionHeading } from '../SectionHeading/SectionHeading'

export function About({ artist }: { artist: ArtistProfile }) {
  return (
    <Box component="section" id="about" sx={{ py: { xs: 10, md: 16 }, bgcolor: 'background.paper' }}>
      <Container>
        <SectionHeading eyebrow="Знакомство" title="О художнице" />
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'minmax(280px, 4fr) minmax(0, 6fr)' }, gap: { xs: 5, md: 12 }, alignItems: 'start' }}>
          <Box sx={{ aspectRatio: '4 / 5', bgcolor: 'divider', display: 'grid', placeItems: 'center', overflow: 'hidden', '& img': { width: '100%', height: '100%', objectFit: 'cover' } }}>
            {artist.avatarUrl ? <img src={artist.avatarUrl} alt={artist.name} /> : <Typography variant="overline" color="text.secondary" sx={{ letterSpacing: '.15em' }}>Фотография художницы</Typography>}
          </Box>
          <Stack spacing={4} sx={{ pt: { md: 4 } }}>
            <Typography variant="h3" sx={{ fontSize: { xs: '2rem', md: '3rem' } }}>{artist.name}</Typography>
            <Typography sx={{ fontSize: { xs: '1.1rem', md: '1.35rem' }, lineHeight: 1.65 }}>{artist.biography}</Typography>
            {artist.artistStatement && <Typography color="text.secondary" sx={{ maxWidth: 680, lineHeight: 1.8 }}>{artist.artistStatement}</Typography>}
          </Stack>
        </Box>
      </Container>
    </Box>
  )
}
