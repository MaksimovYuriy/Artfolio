import type { ArtistProfile } from '../types/artist'

// TEMPORARY CONTENT: this is the single place for all artist-facing copy and links.
export const mockArtist: ArtistProfile = {
  brandName: 'Artfolio',
  name: 'Капитано',
  tagline: 'Живопись, в которой пространство говорит тише слов.',
  biography:
    'Здесь появится краткая биография художницы: образование, творческий путь и важные этапы практики.',
  artistStatement:
    'Здесь будет текст о художественном методе, темах и материалах, с которыми работает автор.',
  email: 'artist@example.com',
  socialLinks: [
    { label: 'Telegram', url: 'https://t.me/example' },
    { label: 'Instagram', url: 'https://instagram.com/example' },
  ],
}
