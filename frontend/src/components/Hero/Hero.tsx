import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import { Box, Container, Link, Stack, Typography } from '@mui/material'
import type { ArtistProfile } from '../../types/artist'
import type { Artwork } from '../../types/artwork'

export function Hero({ artist, featuredArtwork }: { artist: ArtistProfile; featuredArtwork?: Artwork }) {
  return (
    <Box component="section" id="top" sx={{ minHeight: { md: 'calc(100vh - 80px)' }, display: 'flex', alignItems: 'center', py: { xs: 7, md: 8 } }}>
      <Container>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'minmax(0, 5fr) minmax(340px, 4fr)' }, gap: { xs: 6, md: 10 }, alignItems: 'center' }}>
          <Stack sx={{ py: { md: 5 }, minWidth: 0, alignItems: 'flex-start' }}>
            <Typography variant="overline" color="text.secondary" sx={{ letterSpacing: '.2em', mb: 3 }}>Художница · Портфолио</Typography>
            <Typography variant="h1" sx={{ maxWidth: '100%', overflowWrap: 'anywhere', fontSize: { xs: '4rem', sm: '5.7rem', lg: '8rem' } }}>{artist.name}</Typography>
            <Typography sx={{ mt: 4, maxWidth: 470, fontSize: { xs: '1.1rem', md: '1.3rem' }, lineHeight: 1.55 }} color="text.secondary">{artist.tagline}</Typography>
            <Link href="#works" color="text.primary" sx={{ mt: 6, display: 'inline-flex', gap: 1, alignItems: 'center', fontSize: 14 }}>Смотреть работы <ArrowDownwardIcon fontSize="small" /></Link>
          </Stack>
          <Box sx={{ aspectRatio: { xs: '4 / 5', sm: '4 / 3', md: '4 / 5' }, bgcolor: 'background.paper', display: 'grid', placeItems: 'center', overflow: 'hidden', '& img': { width: '100%', height: '100%', objectFit: 'cover' } }}>
            {featuredArtwork?.imageUrl ? <img src={featuredArtwork.imageUrl} alt={featuredArtwork.imageAlt ?? featuredArtwork.title} /> : <Typography variant="overline" color="text.secondary" sx={{ letterSpacing: '.15em' }}>Главная работа</Typography>}
          </Box>
        </Box>
      </Container>
    </Box>
  )
}
