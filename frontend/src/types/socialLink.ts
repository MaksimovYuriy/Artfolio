export const socialPlatforms = [
  { id: 'telegram', label: 'Telegram', prefix: 't.me/', placeholder: 'demo_author' },
  { id: 'instagram', label: 'Instagram', prefix: 'instagram.com/', placeholder: 'demo.author' },
  { id: 'vk', label: 'VK', prefix: 'vk.com/', placeholder: 'demo_author' },
  { id: 'behance', label: 'Behance', prefix: 'behance.net/', placeholder: 'demo-author' },
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
