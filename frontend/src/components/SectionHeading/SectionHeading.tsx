import { Box, Typography } from '@mui/material'

export function SectionHeading({ eyebrow, title }: { eyebrow: string; title: string }) {
  return (
    <Box sx={{ mb: { xs: 5, md: 8 } }}>
      <Typography variant="overline" color="text.secondary" sx={{ letterSpacing: '.18em' }}>{eyebrow}</Typography>
      <Typography variant="h2" sx={{ mt: 1, fontSize: { xs: '2.65rem', md: '4.5rem' } }}>{title}</Typography>
    </Box>
  )
}
