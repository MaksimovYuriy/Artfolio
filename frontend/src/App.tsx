import { useEffect, useState } from 'react'
import { Alert, Box, Button, CircularProgress, Stack, Typography } from '@mui/material'
import { About } from './components/About/About'
import { AdminLogin } from './components/AdminLogin/AdminLogin'
import { Contacts } from './components/Contacts/Contacts'
import { Footer } from './components/Footer/Footer'
import { Gallery } from './components/Gallery/Gallery'
import { Header } from './components/Header/Header'
import { Hero } from './components/Hero/Hero'
import { ArtistServiceError, getArtist } from './services/artistService'
import { getArtworks } from './services/artworksService'
import type { ArtistProfile } from './types/artist'
import type { Artwork } from './types/artwork'

function PublicSite() {
  const [artist, setArtist] = useState<ArtistProfile | null>(null)
  const [artworks, setArtworks] = useState<Artwork[]>([])
  const [error, setError] = useState<string | null>(null)
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true

    Promise.all([getArtist(), getArtworks()])
      .then(([artistData, artworkData]) => {
        if (!active) return
        setArtist(artistData)
        setArtworks(artworkData)
      })
      .catch((caughtError: unknown) => {
        if (!active) return
        setError(caughtError instanceof ArtistServiceError ? caughtError.message : 'Не удалось загрузить сайт.')
      })

    return () => {
      active = false
    }
  }, [reloadKey])

  if (error) {
    return (
      <Box sx={{ minHeight: '100vh', px: 3, display: 'grid', placeItems: 'center' }}>
        <Stack spacing={3} sx={{ width: '100%', maxWidth: 520, alignItems: 'flex-start' }}>
          <Typography variant="h2" sx={{ fontSize: { xs: '2.8rem', md: '4rem' } }}>Портфолио временно недоступно</Typography>
          <Alert severity="error" variant="outlined" sx={{ width: '100%' }}>{error}</Alert>
          <Button
            variant="contained"
            onClick={() => {
              setError(null)
              setArtist(null)
              setReloadKey((key) => key + 1)
            }}
          >
            Попробовать снова
          </Button>
        </Stack>
      </Box>
    )
  }

  if (!artist) {
    return (
      <Box sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center' }}>
        <CircularProgress size={28} aria-label="Загрузка сайта" />
      </Box>
    )
  }

  return (
    <>
      <Header brand={artist.brandName} />
      <main>
        <Hero artist={artist} featuredArtwork={artworks[0]} />
        <Gallery artworks={artworks} />
        <About artist={artist} />
        <Contacts artist={artist} />
      </main>
      <Footer artist={artist} />
    </>
  )
}

function App() {
  return window.location.pathname.startsWith('/admin') ? <AdminLogin /> : <PublicSite />
}

export default App
