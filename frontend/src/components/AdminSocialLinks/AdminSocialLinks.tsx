import { type FormEvent, useEffect, useState } from 'react'
import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined'
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import {
  getSocialHandles,
  replaceSocialHandles,
  SocialLinksServiceError,
} from '../../services/socialLinksService'
import {
  emptySocialHandles,
  socialPlatforms,
  type SocialHandles,
  type SocialPlatform,
} from '../../types/socialLink'

interface AdminSocialLinksProps {
  onSessionExpired: (message: string) => void
}

export function AdminSocialLinks({ onSessionExpired }: AdminSocialLinksProps) {
  const [handles, setHandles] = useState<SocialHandles>({ ...emptySocialHandles })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    let active = true

    getSocialHandles()
      .then((data) => {
        if (active) setHandles(data)
      })
      .catch((caughtError: unknown) => {
        if (!active) return
        if (caughtError instanceof SocialLinksServiceError && caughtError.unauthorized) {
          onSessionExpired(caughtError.message)
          return
        }
        setError(caughtError instanceof Error ? caughtError.message : 'Не удалось загрузить социальные сети.')
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [onSessionExpired])

  function updateHandle(platform: SocialPlatform, value: string) {
    setHandles((current) => ({ ...current, [platform]: value }))
    setError(null)
    setSaved(false)
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSaving(true)
    setError(null)
    setSaved(false)

    try {
      await replaceSocialHandles(handles)
      setSaved(true)
    } catch (caughtError) {
      if (caughtError instanceof SocialLinksServiceError && caughtError.unauthorized) {
        onSessionExpired(caughtError.message)
        return
      }
      setError(caughtError instanceof Error ? caughtError.message : 'Не удалось сохранить социальные сети.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Box component="form" onSubmit={handleSubmit} noValidate>
      <Typography variant="h3" sx={{ fontSize: { xs: '2rem', md: '2.5rem' } }}>Социальные сети</Typography>
      <Typography color="text.secondary" sx={{ mt: 2, mb: 4, maxWidth: 650, lineHeight: 1.7 }}>
        Укажите имя аккаунта или вставьте полную ссылку. Пустые площадки не будут показаны посетителям.
      </Typography>

      {loading ? (
        <Stack direction="row" spacing={2} sx={{ alignItems: 'center', minHeight: 72 }} aria-live="polite">
          <CircularProgress size={22} />
          <Typography color="text.secondary">Загружаем ссылки…</Typography>
        </Stack>
      ) : (
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, minmax(0, 1fr))' }, gap: 3, maxWidth: 900 }}>
          {socialPlatforms.map((platform) => (
            <TextField
              key={platform.id}
              label={platform.label}
              value={handles[platform.id]}
              onChange={(event) => updateHandle(platform.id, event.target.value)}
              placeholder={platform.placeholder}
              helperText={platform.prefix}
              disabled={saving}
              slotProps={{ htmlInput: { maxLength: 256, autoComplete: 'off' } }}
            />
          ))}
        </Box>
      )}

      {error && <Alert severity="error" variant="outlined" sx={{ mt: 4, maxWidth: 900 }}>{error}</Alert>}
      {saved && <Alert severity="success" variant="outlined" sx={{ mt: 4, maxWidth: 900 }}>Социальные сети сохранены.</Alert>}

      <Button
        type="submit"
        variant="outlined"
        size="large"
        disabled={loading || saving}
        startIcon={saving ? <CircularProgress size={19} color="inherit" /> : <SaveOutlinedIcon />}
        sx={{ mt: 4, minHeight: 52, px: 4 }}
      >
        {saving ? 'Сохраняем…' : 'Сохранить социальные сети'}
      </Button>
    </Box>
  )
}
