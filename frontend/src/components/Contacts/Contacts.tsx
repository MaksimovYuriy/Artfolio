import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined'
import NorthEastIcon from '@mui/icons-material/NorthEast'
import { Box, Container, Link, Stack, Typography } from '@mui/material'
import type { ArtistProfile } from '../../types/artist'
import { SectionHeading } from '../SectionHeading/SectionHeading'

export function Contacts({ artist }: { artist: ArtistProfile }) {
  const contacts = [
    artist.email && { label: 'Email', value: artist.email, href: `mailto:${artist.email}` },
    ...(artist.socialLinks ?? []).map((link) => ({ label: link.label, value: link.label, href: link.url })),
  ].filter(Boolean) as { label: string; value: string; href: string }[]

  return (
    <Box component="section" id="contacts" sx={{ py: { xs: 10, md: 16 } }}>
      <Container>
        <SectionHeading eyebrow="Связаться" title="Контакты" />
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 2fr' }, gap: { xs: 5, md: 12 } }}>
          <Typography color="text.secondary" sx={{ maxWidth: 340, lineHeight: 1.7 }}>По вопросам сотрудничества, выставок и приобретения работ.</Typography>
          <Stack sx={{ borderTop: 1, borderColor: 'divider' }}>
            {contacts.map((contact) => (
              <Link key={`${contact.label}-${contact.href}`} href={contact.href} target={contact.href.startsWith('http') ? '_blank' : undefined} rel="noreferrer" color="text.primary" sx={{ py: 2.5, borderBottom: 1, borderColor: 'divider', display: 'grid', gridTemplateColumns: { xs: '1fr auto', sm: '160px 1fr auto' }, gap: 2, alignItems: 'center', '&:hover': { color: 'primary.main' } }}>
                <Typography variant="overline" color="text.secondary">{contact.label}</Typography>
                <Typography sx={{ display: { xs: 'none', sm: 'block' } }}>{contact.value}</Typography>
                {contact.label === 'Email' ? <EmailOutlinedIcon fontSize="small" /> : <NorthEastIcon fontSize="small" />}
              </Link>
            ))}
          </Stack>
        </Box>
      </Container>
    </Box>
  )
}
