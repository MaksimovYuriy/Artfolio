import { useEffect, useState } from 'react'
import { Box, CircularProgress } from '@mui/material'
import { About } from './components/About/About'
import { AdminPlaceholder } from './components/AdminPlaceholder/AdminPlaceholder'
import { Contacts } from './components/Contacts/Contacts'
import { Footer } from './components/Footer/Footer'
import { Gallery } from './components/Gallery/Gallery'
import { Header } from './components/Header/Header'
import { Hero } from './components/Hero/Hero'
import { getArtist } from './services/artistService'
import { getArtworks } from './services/artworksService'
import type { ArtistProfile } from './types/artist'
import type { Artwork } from './types/artwork'

function PublicSite() {
  const [artist, setArtist] = useState<ArtistProfile | null>(null)
  const [artworks, setArtworks] = useState<Artwork[]>([])

  useEffect(() => {
    Promise.all([getArtist(), getArtworks()]).then(([artistData, artworkData]) => {
      setArtist(artistData)
      setArtworks(artworkData)
    })
  }, [])

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
  return window.location.pathname.startsWith('/admin') ? <AdminPlaceholder /> : <PublicSite />
}

export default App
