export const socialPlatforms = [
  { id: 'telegram', label: 'Telegram', prefix: 't.me/', placeholder: 'anna_art' },
  { id: 'instagram', label: 'Instagram', prefix: 'instagram.com/', placeholder: 'anna.art' },
  { id: 'vk', label: 'VK', prefix: 'vk.com/', placeholder: 'anna_art' },
  { id: 'behance', label: 'Behance', prefix: 'behance.net/', placeholder: 'anna-art' },
] as const

export type SocialPlatform = typeof socialPlatforms[number]['id']

export interface SocialLink {
  label: string
  url: string
}

export interface AdminSocialLink {
  platform: SocialPlatform
  handle: string
}

export type SocialHandles = Record<SocialPlatform, string>

export const emptySocialHandles: SocialHandles = {
  telegram: '',
  instagram: '',
  vk: '',
  behance: '',
}
