import { Box, Container } from '@mui/material'
import type { Artwork } from '../../types/artwork'
import { ArtworkCard } from '../ArtworkCard/ArtworkCard'
import { SectionHeading } from '../SectionHeading/SectionHeading'

export function Gallery({ artworks }: { artworks: Artwork[] }) {
  return (
    <Box component="section" id="works" sx={{ py: { xs: 10, md: 16 } }}>
      <Container>
        <SectionHeading eyebrow="Избранное" title="Работы" />
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }, gap: { xs: 6, md: 4 }, alignItems: 'start' }}>
          {artworks.map((artwork) => <ArtworkCard key={artwork.id} artwork={artwork} />)}
        </Box>
      </Container>
    </Box>
  )
}
