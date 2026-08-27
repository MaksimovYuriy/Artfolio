import { Box, Container, Stack, Typography } from '@mui/material'
import type { ArtistProfile } from '../../types/artist'
import { SectionHeading } from '../SectionHeading/SectionHeading'

export function About({ artist }: { artist: ArtistProfile }) {
  return (
    <Box component="section" id="about" sx={{ py: { xs: 10, md: 16 }, bgcolor: 'background.paper' }}>
      <Container>
        <SectionHeading eyebrow="Знакомство" title="О художнице" />
        <Box sx={{ maxWidth: 760 }}>
          <Stack spacing={4}>
            <Typography variant="h3" sx={{ fontSize: { xs: '2rem', md: '3rem' } }}>{artist.name}</Typography>
            <Typography sx={{ fontSize: { xs: '1.1rem', md: '1.35rem' }, lineHeight: 1.65 }}>{artist.biography}</Typography>
            {artist.artistStatement && <Typography color="text.secondary" sx={{ maxWidth: 680, lineHeight: 1.8 }}>{artist.artistStatement}</Typography>}
          </Stack>
        </Box>
      </Container>
    </Box>
  )
}
