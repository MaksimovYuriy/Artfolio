import { type FormEvent, useEffect, useState } from 'react'
import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import ArrowForwardIcon from '@mui/icons-material/ArrowForward'
import KeyOutlinedIcon from '@mui/icons-material/KeyOutlined'
import VisibilityOffOutlinedIcon from '@mui/icons-material/VisibilityOffOutlined'
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined'
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Container,
  IconButton,
  InputAdornment,
  Link,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import {
  AdminAuthError,
  createAdminSession,
  revokeAdminSession,
  verifyAdminSession,
} from '../../services/adminService'
import { AdminProfile } from '../AdminProfile/AdminProfile'
import { AdminArtworks } from '../AdminArtworks/AdminArtworks'

type AuthStatus = 'checking' | 'anonymous' | 'authenticated'
type AdminSection = 'profile' | 'artworks'

export function AdminLogin() {
  const [accessKey, setAccessKey] = useState('')
  const [showKey, setShowKey] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [revoking, setRevoking] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [authStatus, setAuthStatus] = useState<AuthStatus>('checking')
  const [section, setSection] = useState<AdminSection>('profile')

  useEffect(() => {
    let active = true

    verifyAdminSession()
      .then((authenticated) => {
        if (active) setAuthStatus(authenticated ? 'authenticated' : 'anonymous')
      })
      .catch((caughtError: unknown) => {
        if (!active) return

        setError(caughtError instanceof AdminAuthError ? caughtError.message : 'Произошла неизвестная ошибка.')
        setAuthStatus('anonymous')
      })

    return () => {
      active = false
    }
  }, [])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const normalizedKey = accessKey.trim()
    if (!normalizedKey) {
      setError('Введите ключ доступа.')
      return
    }

    setSubmitting(true)
    setError(null)

    try {
      await createAdminSession(normalizedKey)
      setAccessKey('')
      setAuthStatus('authenticated')
    } catch (caughtError) {
      setError(caughtError instanceof AdminAuthError ? caughtError.message : 'Произошла неизвестная ошибка.')
    } finally {
      setSubmitting(false)
    }
  }

  async function handleLogout() {
    setRevoking(true)
    setError(null)

    try {
      await revokeAdminSession()
      setAuthStatus('anonymous')
    } catch (caughtError) {
      setError(caughtError instanceof AdminAuthError ? caughtError.message : 'Произошла неизвестная ошибка.')
    } finally {
      setRevoking(false)
    }
  }

  function handleSessionExpired(message: string) {
    setError(message)
    setAuthStatus('anonymous')
  }

  return (
    <Box component="main" sx={{ minHeight: '100dvh', display: 'flex' }}>
      <Container disableGutters sx={{ display: 'flex', flex: 1, maxWidth: 'none !important' }}>
        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: authStatus === 'authenticated'
              ? '1fr'
              : { xs: '1fr', md: 'minmax(360px, 47%) minmax(0, 53%)' },
            flex: 1,
          }}
        >
          <Box
            sx={{
              minHeight: { xs: 280, md: '100dvh' },
              px: { xs: 3, sm: 6, lg: 9 },
              py: { xs: 4, md: 6 },
              bgcolor: 'primary.main',
              color: '#f7f4ee',
              flexDirection: 'column',
              justifyContent: 'space-between',
              position: 'relative',
              overflow: 'hidden',
              display: authStatus === 'authenticated' ? 'none' : 'flex',
            }}
          >
            <Typography
              component="a"
              href="/"
              variant="h6"
              color="inherit"
              sx={{ fontFamily: 'Georgia, serif', textDecoration: 'none', position: 'relative', zIndex: 1 }}
            >
              Artfolio
            </Typography>

            <Box
              aria-hidden="true"
              sx={{
                position: 'absolute',
                width: { xs: 260, md: 430 },
                height: { xs: 260, md: 430 },
                border: '1px solid rgba(255,255,255,.15)',
                borderRadius: '50%',
                right: { xs: -100, md: -150 },
                bottom: { xs: -150, md: -80 },
                '&::after': {
                  content: '""',
                  position: 'absolute',
                  inset: '18%',
                  border: '1px solid rgba(255,255,255,.12)',
                  borderRadius: '50%',
                },
              }}
            />

            <Stack sx={{ maxWidth: 560, position: 'relative', zIndex: 1, my: { xs: 7, md: 10 } }}>
              <Typography variant="overline" sx={{ letterSpacing: '.2em', opacity: 0.72 }}>
                Личное пространство
              </Typography>
              <Typography
                variant="h1"
                sx={{ mt: 3, fontSize: { xs: '3.4rem', sm: '4.5rem', lg: '6.4rem' }, maxWidth: 560 }}
              >
                Всё важное — за одним ключом.
              </Typography>
              <Typography sx={{ mt: 4, maxWidth: 410, lineHeight: 1.75, opacity: 0.76 }}>
                Здесь автор управляет работами, текстами и тем, как портфолио выглядит для посетителей.
              </Typography>
            </Stack>

            <Typography variant="body2" sx={{ display: { xs: 'none', md: 'block' }, opacity: 0.58 }}>
              Закрытая зона · Только для автора
            </Typography>
          </Box>

          <Box
            sx={{
              px: { xs: 3, sm: 8, lg: 13 },
              py: { xs: 7, md: 8 },
              display: 'flex',
              flexDirection: 'column',
              justifyContent: authStatus === 'authenticated' ? 'flex-start' : 'center',
              bgcolor: 'background.default',
            }}
          >
            <Box sx={{ width: '100%', maxWidth: authStatus === 'authenticated' ? 1080 : 520, mx: 'auto' }}>
              <Link
                href="/"
                color="text.secondary"
                sx={{ display: 'inline-flex', alignItems: 'center', gap: 1, mb: { xs: 7, md: 10 }, fontSize: 14 }}
              >
                <ArrowBackIcon sx={{ fontSize: 18 }} /> Вернуться к портфолио
              </Link>

              {authStatus === 'checking' ? (
                <Stack sx={{ alignItems: 'flex-start' }} aria-live="polite">
                  <CircularProgress size={28} />
                  <Typography variant="h2" sx={{ mt: 4, fontSize: { xs: '3rem', md: '4rem' } }}>
                    Проверяем доступ
                  </Typography>
                  <Typography color="text.secondary" sx={{ mt: 3, lineHeight: 1.7 }}>
                    Это займёт всего несколько секунд.
                  </Typography>
                </Stack>
              ) : authStatus === 'authenticated' ? (
				section === 'profile' ? (
				  <AdminProfile
					loggingOut={revoking}
					onLogout={handleLogout}
					onOpenArtworks={() => setSection('artworks')}
					onSessionExpired={handleSessionExpired}
				  />
				) : (
				  <AdminArtworks
					loggingOut={revoking}
					onLogout={handleLogout}
					onOpenProfile={() => setSection('profile')}
					onSessionExpired={handleSessionExpired}
				  />
				)
              ) : (
                <>
                  <Typography variant="overline" color="primary.main" sx={{ letterSpacing: '.18em' }}>
                    Вход для автора
                  </Typography>
                  <Typography variant="h2" sx={{ mt: 2.5, fontSize: { xs: '3rem', md: '4.2rem' } }}>
                    Добро пожаловать
                  </Typography>
                  <Typography color="text.secondary" sx={{ mt: 3, mb: 6, lineHeight: 1.7, maxWidth: 430 }}>
                    Введите персональный ключ, созданный при настройке портфолио.
                  </Typography>

                  <Box component="form" onSubmit={handleSubmit} noValidate>
                    <Stack spacing={3}>
                      <TextField
                        label="Ключ доступа"
                        type={showKey ? 'text' : 'password'}
                        value={accessKey}
                        onChange={(event) => {
                          setAccessKey(event.target.value)
                          if (error) setError(null)
                        }}
                        placeholder="artfolio_••••••••••••"
                        autoComplete="current-password"
                        autoFocus
                        disabled={submitting}
                        error={Boolean(error)}
                        slotProps={{
                          input: {
                            endAdornment: (
                              <InputAdornment position="end">
                                <IconButton
                                  onClick={() => setShowKey((visible) => !visible)}
                                  edge="end"
                                  aria-label={showKey ? 'Скрыть ключ' : 'Показать ключ'}
                                >
                                  {showKey ? <VisibilityOffOutlinedIcon /> : <VisibilityOutlinedIcon />}
                                </IconButton>
                              </InputAdornment>
                            ),
                          },
                        }}
                        sx={{
                          '& .MuiOutlinedInput-root': { bgcolor: 'rgba(255,255,255,.38)', minHeight: 60 },
                        }}
                      />

                      {error && <Alert severity="error" variant="outlined">{error}</Alert>}

                      <Button
                        type="submit"
                        variant="contained"
                        size="large"
                        disabled={submitting}
                        endIcon={!submitting && <ArrowForwardIcon />}
                        sx={{ minHeight: 56, justifyContent: 'space-between', px: 3, boxShadow: 'none' }}
                      >
                        {submitting ? <CircularProgress size={22} color="inherit" /> : 'Продолжить'}
                      </Button>
                    </Stack>
                  </Box>

                  <Stack direction="row" spacing={1.5} sx={{ mt: 5, alignItems: 'flex-start' }}>
                    <KeyOutlinedIcon color="disabled" sx={{ fontSize: 19, mt: '2px' }} />
                    <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.65 }}>
                      Ключ не сохраняется в браузере и используется только для создания защищённой сессии.
                    </Typography>
                  </Stack>
                </>
              )}
            </Box>
          </Box>
        </Box>
      </Container>
    </Box>
  )
}
