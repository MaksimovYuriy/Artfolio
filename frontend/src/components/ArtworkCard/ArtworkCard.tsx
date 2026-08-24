import { Box, Stack, Typography } from '@mui/material'
import type { Artwork } from '../../types/artwork'

export function ArtworkCard({ artwork }: { artwork: Artwork }) {
  return (
    <Box component="article">
      <Box sx={{ aspectRatio: '4 / 5', bgcolor: 'background.paper', overflow: 'hidden', position: 'relative', '& img': { width: '100%', height: '100%', objectFit: 'contain', transition: 'transform .6s ease' }, '&:hover img': { transform: 'scale(1.02)' } }}>
        {artwork.imageUrl ? (
          <img src={artwork.imageUrl} alt={artwork.imageAlt ?? artwork.title} />
        ) : (
          <Stack sx={{ position: 'absolute', inset: 0, alignItems: 'center', justifyContent: 'center', color: 'text.secondary' }}>
            <Typography variant="overline" sx={{ letterSpacing: '.15em' }}>Изображение работы</Typography>
          </Stack>
        )}
      </Box>
      <Stack direction="row" spacing={2} sx={{ pt: 2, justifyContent: 'space-between', alignItems: 'baseline' }}>
        <Typography variant="h6" sx={{ fontFamily: 'Georgia, serif', fontWeight: 400 }}>{artwork.title}</Typography>
        {artwork.year !== undefined && <Typography variant="body2" color="text.secondary">{artwork.year}</Typography>}
      </Stack>
      {artwork.technique && <Typography variant="body2" color="text.secondary" sx={{ mt: .5 }}>{artwork.technique}</Typography>}
    </Box>
  )
}
