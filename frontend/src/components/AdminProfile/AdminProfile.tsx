import { type FormEvent, useEffect, useState } from 'react'
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined'
import CollectionsOutlinedIcon from '@mui/icons-material/CollectionsOutlined'
import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined'
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Divider,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import {
  AdminProfileError,
  emptyArtistProfile,
  getEditableArtistProfile,
  updateArtistProfile,
} from '../../services/adminService'
import type { EditableArtistProfile } from '../../types/artist'
import { AdminSocialLinks } from '../AdminSocialLinks/AdminSocialLinks'

interface AdminProfileProps {
  loggingOut: boolean
  onLogout: () => void
  onOpenArtworks: () => void
  onSessionExpired: (message: string) => void
}

export function AdminProfile({ loggingOut, onLogout, onOpenArtworks, onSessionExpired }: AdminProfileProps) {
  const [profile, setProfile] = useState<EditableArtistProfile>({ ...emptyArtistProfile })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    let active = true

    getEditableArtistProfile()
      .then((data) => {
        if (active) setProfile(data)
      })
      .catch((caughtError: unknown) => {
        if (!active) return
        setError(caughtError instanceof AdminProfileError ? caughtError.message : 'Произошла неизвестная ошибка.')
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [])

  function updateField(field: keyof EditableArtistProfile, value: string) {
    setProfile((current) => ({ ...current, [field]: value }))
    setError(null)
    setSaved(false)
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSaving(true)
    setError(null)
    setSaved(false)

    try {
      await updateArtistProfile(profile)
      setSaved(true)
    } catch (caughtError) {
      if (caughtError instanceof AdminProfileError && caughtError.unauthorized) {
        onSessionExpired(caughtError.message)
        return
      }
      setError(caughtError instanceof AdminProfileError ? caughtError.message : 'Произошла неизвестная ошибка.')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <Stack spacing={3} sx={{ minHeight: 360, alignItems: 'center', justifyContent: 'center' }} aria-live="polite">
        <CircularProgress size={28} />
        <Typography color="text.secondary">Загружаем профиль…</Typography>
      </Stack>
    )
  }

  return (
    <Box>
      <Stack direction={{ xs: 'column', sm: 'row' }} sx={{ justifyContent: 'space-between', gap: 3, mb: 6 }}>
        <Box>
          <Typography variant="overline" color="primary.main" sx={{ letterSpacing: '.18em' }}>Личный кабинет</Typography>
          <Typography variant="h2" sx={{ mt: 1.5, fontSize: { xs: '2.8rem', md: '4rem' } }}>Профиль</Typography>
          <Typography color="text.secondary" sx={{ mt: 2, maxWidth: 560, lineHeight: 1.7 }}>
            Здесь можно обновить тексты и контакты, которые видят посетители портфолио.
          </Typography>
        </Box>
		<Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5} sx={{ alignSelf: 'flex-start' }}>
		  <Button type="button" onClick={onOpenArtworks} startIcon={<CollectionsOutlinedIcon />} color="inherit">Работы</Button>
		  <Button type="button" onClick={onLogout} disabled={loggingOut} startIcon={loggingOut ? <CircularProgress size={18} /> : <LogoutOutlinedIcon />} color="inherit">Выйти</Button>
		</Stack>
      </Stack>

      <Divider sx={{ mb: 5 }} />

      <Box component="form" onSubmit={handleSubmit} noValidate>
        <Box sx={{ maxWidth: 760 }}>
          <Stack spacing={3}>
            <TextField required label="Имя" value={profile.name} onChange={(e) => updateField('name', e.target.value)} slotProps={{ htmlInput: { maxLength: 64 } }} helperText={`${profile.name.length}/64`} />
            <TextField label="Слоган" value={profile.tagline} onChange={(e) => updateField('tagline', e.target.value)} slotProps={{ htmlInput: { maxLength: 256 } }} helperText={`${profile.tagline.length}/256`} />
            <TextField label="Биография" value={profile.biography} onChange={(e) => updateField('biography', e.target.value)} multiline minRows={5} />
            <TextField label="Творческое высказывание" value={profile.artistStatement} onChange={(e) => updateField('artistStatement', e.target.value)} multiline minRows={5} />
            <TextField label="Email" type="email" value={profile.email} onChange={(e) => updateField('email', e.target.value)} />
          </Stack>
        </Box>

        {error && <Alert severity="error" variant="outlined" sx={{ mt: 4 }}>{error}</Alert>}
        {saved && <Alert severity="success" variant="outlined" sx={{ mt: 4 }}>Профиль сохранён.</Alert>}

        <Button type="submit" variant="contained" size="large" disabled={saving || !profile.name.trim()} startIcon={saving ? <CircularProgress size={19} color="inherit" /> : <SaveOutlinedIcon />} sx={{ mt: 5, minHeight: 52, px: 4, boxShadow: 'none' }}>
          {saving ? 'Сохраняем…' : 'Сохранить профиль'}
        </Button>
      </Box>

      <Divider sx={{ my: 7 }} />
      <AdminSocialLinks onSessionExpired={onSessionExpired} />
    </Box>
  )
}
