import { type DragEvent, useEffect, useState } from 'react'
import AddOutlinedIcon from '@mui/icons-material/AddOutlined'
import ArrowDownwardOutlinedIcon from '@mui/icons-material/ArrowDownwardOutlined'
import ArrowUpwardOutlinedIcon from '@mui/icons-material/ArrowUpwardOutlined'
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutlined'
import EditOutlinedIcon from '@mui/icons-material/EditOutlined'
import ImageOutlinedIcon from '@mui/icons-material/ImageOutlined'
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined'
import PersonOutlineIcon from '@mui/icons-material/PersonOutlined'
import UploadFileOutlinedIcon from '@mui/icons-material/UploadFileOutlined'
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  IconButton,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import {
  ArtworkServiceError,
  createArtwork,
  deleteArtwork,
  getAdminArtworks,
  reorderArtworks,
  updateArtwork,
} from '../../services/artworksService'
import type { AdminArtwork, ArtworkInput } from '../../types/artwork'

const maxImageSize = 12 * 1024 * 1024
const acceptedImageTypes = ['image/jpeg', 'image/png']

interface AdminArtworksProps {
  loggingOut: boolean
  onLogout: () => void
  onOpenProfile: () => void
  onSessionExpired: (message: string) => void
}

interface EditorState {
  artwork: AdminArtwork | null
  input: ArtworkInput
  image: File | null
}

const emptyInput: ArtworkInput = {
  title: '',
  description: '',
  technique: '',
  year: '',
  imageAlt: '',
}

export function AdminArtworks({ loggingOut, onLogout, onOpenProfile, onSessionExpired }: AdminArtworksProps) {
  const [artworks, setArtworks] = useState<AdminArtwork[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editor, setEditor] = useState<EditorState | null>(null)
  const [deletingID, setDeletingID] = useState<number | null>(null)
  const [savedOrder, setSavedOrder] = useState<number[]>([])
  const [savingOrder, setSavingOrder] = useState(false)

  const currentOrder = artworks.map((artwork) => artwork.id)
  const orderChanged = !sameOrder(currentOrder, savedOrder)
  const featuredArtworkID = artworks.find((artwork) => artwork.isPublished)?.id

  useEffect(() => {
    let active = true
    getAdminArtworks()
      .then((items) => {
        if (active) {
          setArtworks(items)
          setSavedOrder(items.map((item) => item.id))
        }
      })
      .catch((caughtError: unknown) => {
        if (!active) return
        handleServiceError(caughtError, onSessionExpired, setError)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => { active = false }
  }, [onSessionExpired])

  useEffect(() => {
    if (!orderChanged) return
    const warnBeforeUnload = (event: BeforeUnloadEvent) => event.preventDefault()
    window.addEventListener('beforeunload', warnBeforeUnload)
    return () => window.removeEventListener('beforeunload', warnBeforeUnload)
  }, [orderChanged])

  function navigateWithOrderCheck(action: () => void) {
    if (orderChanged && !window.confirm('Порядок работ не сохранён. Уйти без сохранения?')) return
    action()
  }

  function moveArtwork(index: number, offset: -1 | 1) {
    setArtworks((items) => {
      const target = index + offset
      if (target < 0 || target >= items.length) return items
      const reordered = [...items]
      ;[reordered[index], reordered[target]] = [reordered[target], reordered[index]]
      return reordered
    })
    setError(null)
  }

  function cancelOrder() {
    const positions = new Map(savedOrder.map((id, index) => [id, index]))
    setArtworks((items) => [...items].sort((a, b) => (positions.get(a.id) ?? 0) - (positions.get(b.id) ?? 0)))
    setError(null)
  }

  async function saveOrder() {
    setSavingOrder(true)
    setError(null)
    try {
      await reorderArtworks(currentOrder)
      setArtworks((items) => items.map((item, position) => ({ ...item, position })))
      setSavedOrder(currentOrder)
    } catch (caughtError) {
      handleServiceError(caughtError, onSessionExpired, setError)
    } finally {
      setSavingOrder(false)
    }
  }

  function openCreate() {
    setError(null)
    setEditor({ artwork: null, input: { ...emptyInput }, image: null })
  }

  function openEdit(artwork: AdminArtwork) {
    setError(null)
    setEditor({
      artwork,
      image: null,
      input: {
        title: artwork.title,
        description: artwork.description ?? '',
        technique: artwork.technique ?? '',
        year: artwork.year === undefined ? '' : String(artwork.year),
        imageAlt: artwork.imageAlt ?? '',
      },
    })
  }

  async function handleDelete(artwork: AdminArtwork) {
    if (!window.confirm(`Удалить работу «${artwork.title}»? Это действие нельзя отменить.`)) return
    setDeletingID(artwork.id)
    setError(null)
    try {
      await deleteArtwork(artwork.id)
      setArtworks((items) => items.filter((item) => item.id !== artwork.id))
      setSavedOrder((ids) => ids.filter((id) => id !== artwork.id))
    } catch (caughtError) {
      handleServiceError(caughtError, onSessionExpired, setError)
    } finally {
      setDeletingID(null)
    }
  }

  return (
    <Box>
      <Stack direction={{ xs: 'column', md: 'row' }} sx={{ justifyContent: 'space-between', gap: 3, mb: 6 }}>
        <Box>
          <Typography variant="overline" color="primary.main" sx={{ letterSpacing: '.18em' }}>Личный кабинет</Typography>
          <Typography variant="h2" sx={{ mt: 1.5, fontSize: { xs: '2.8rem', md: '4rem' } }}>Работы</Typography>
          <Typography color="text.secondary" sx={{ mt: 2, maxWidth: 560, lineHeight: 1.7 }}>
            Добавляйте работы, сохраняйте черновики и управляйте тем, что видно посетителям.
          </Typography>
        </Box>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5} sx={{ alignSelf: { md: 'flex-start' } }}>
          <Button onClick={() => navigateWithOrderCheck(onOpenProfile)} startIcon={<PersonOutlineIcon />} color="inherit">Профиль</Button>
          <Button onClick={() => navigateWithOrderCheck(onLogout)} disabled={loggingOut} startIcon={loggingOut ? <CircularProgress size={18} /> : <LogoutOutlinedIcon />} color="inherit">Выйти</Button>
        </Stack>
      </Stack>

      <Divider sx={{ mb: 5 }} />
      <Stack direction={{ xs: 'column', sm: 'row' }} sx={{ justifyContent: 'space-between', alignItems: { sm: 'center' }, gap: 2, mb: 4 }}>
        <Typography color="text.secondary">
          {loading ? 'Загружаем список…' : `${artworks.length} ${workCountLabel(artworks.length)}`}
        </Typography>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5}>
          <Button variant="outlined" onClick={cancelOrder} disabled={!orderChanged || savingOrder}>Отменить изменения</Button>
          <Button variant="contained" onClick={() => void saveOrder()} disabled={!orderChanged || savingOrder} sx={{ boxShadow: 'none' }}>
            {savingOrder ? <CircularProgress size={20} color="inherit" /> : 'Сохранить порядок'}
          </Button>
          <Button variant="contained" startIcon={<AddOutlinedIcon />} onClick={openCreate} disabled={orderChanged || savingOrder} sx={{ minHeight: 48, px: 3, boxShadow: 'none' }}>
            Добавить работу
          </Button>
        </Stack>
      </Stack>

      {error && <Alert severity="error" variant="outlined" sx={{ mb: 4 }}>{error}</Alert>}

      {loading ? (
        <Box sx={{ minHeight: 280, display: 'grid', placeItems: 'center' }}><CircularProgress size={28} /></Box>
      ) : artworks.length === 0 ? (
        <Stack sx={{ minHeight: 320, alignItems: 'center', justifyContent: 'center', border: '1px solid', borderColor: 'divider', textAlign: 'center', px: 3 }}>
          <ImageOutlinedIcon color="disabled" sx={{ fontSize: 44 }} />
          <Typography variant="h5" sx={{ mt: 2, fontFamily: 'Georgia, serif' }}>Здесь пока нет работ</Typography>
          <Typography color="text.secondary" sx={{ mt: 1.5 }}>Добавьте первую работу и сохраните её как черновик или сразу опубликуйте.</Typography>
        </Stack>
      ) : (
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }, gap: 3 }}>
          {artworks.map((artwork, index) => (
            <Box key={artwork.id} sx={{ border: '1px solid', borderColor: 'divider', bgcolor: 'rgba(255,255,255,.24)' }}>
              <Box sx={{ aspectRatio: '4 / 5', bgcolor: 'background.paper', display: 'grid', placeItems: 'center', overflow: 'hidden' }}>
                <Box component="img" src={artwork.imageUrl} alt={artwork.imageAlt || artwork.title} sx={{ width: '100%', height: '100%', objectFit: 'contain' }} />
              </Box>
              <Stack spacing={2} sx={{ p: 2.5 }}>
                <Stack direction="row" spacing={2} sx={{ justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Box>
                    <Typography variant="h6" sx={{ fontFamily: 'Georgia, serif', fontWeight: 400 }}>{artwork.title}</Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mt: .5 }}>
                      {[artwork.year, artwork.technique].filter(Boolean).join(' · ') || 'Без дополнительных сведений'}
                    </Typography>
                  </Box>
                  <Stack spacing={1} sx={{ alignItems: 'flex-end' }}>
                    {artwork.id === featuredArtworkID && <Chip size="small" color="primary" label="Главная" />}
                    <Chip size="small" variant="outlined" color={artwork.isPublished ? 'primary' : 'default'} label={artwork.isPublished ? 'Опубликовано' : 'Черновик'} />
                  </Stack>
                </Stack>
                <Stack direction="row" spacing={1} sx={{ justifyContent: 'flex-end' }}>
                  <IconButton aria-label={`Переместить «${artwork.title}» раньше`} disabled={index === 0 || savingOrder} onClick={() => moveArtwork(index, -1)}><ArrowUpwardOutlinedIcon /></IconButton>
                  <IconButton aria-label={`Переместить «${artwork.title}» позже`} disabled={index === artworks.length - 1 || savingOrder} onClick={() => moveArtwork(index, 1)}><ArrowDownwardOutlinedIcon /></IconButton>
                  <IconButton aria-label={`Редактировать «${artwork.title}»`} disabled={orderChanged || savingOrder} onClick={() => openEdit(artwork)}><EditOutlinedIcon /></IconButton>
                  <IconButton aria-label={`Удалить «${artwork.title}»`} color="error" disabled={orderChanged || savingOrder || deletingID === artwork.id} onClick={() => void handleDelete(artwork)}>
                    {deletingID === artwork.id ? <CircularProgress size={20} /> : <DeleteOutlineIcon />}
                  </IconButton>
                </Stack>
              </Stack>
            </Box>
          ))}
        </Box>
      )}

    {editor && (
    <ArtworkEditor
      key={editor.artwork?.id ?? 'new'}
      state={editor}
      onClose={() => setEditor(null)}
      onSaved={(saved) => {
      setArtworks((items) => {
        const exists = items.some((item) => item.id === saved.id)
        return exists ? items.map((item) => item.id === saved.id ? saved : item) : [...items, saved]
      })
      setSavedOrder((ids) => ids.includes(saved.id) ? ids : [...ids, saved.id])
      setEditor(null)
      }}
      onSessionExpired={onSessionExpired}
    />
    )}
    </Box>
  )
}

function sameOrder(left: number[], right: number[]): boolean {
  return left.length === right.length && left.every((id, index) => id === right[index])
}

function ArtworkEditor({ state, onClose, onSaved, onSessionExpired }: {
  state: EditorState
  onClose: () => void
  onSaved: (artwork: AdminArtwork) => void
  onSessionExpired: (message: string) => void
}) {
  const [input, setInput] = useState<ArtworkInput>(state.input)
  const [image, setImage] = useState<File | null>(state.image)
  const [preview, setPreview] = useState<string | null>(state.artwork?.imageUrl ?? null)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    return () => {
    if (preview?.startsWith('blob:')) URL.revokeObjectURL(preview)
    }
  }, [preview])

  function setField(field: keyof ArtworkInput, value: string) {
    setInput((current) => ({ ...current, [field]: value }))
    setError(null)
  }

  function chooseImage(file: File | undefined) {
    if (!file) return
    if (!acceptedImageTypes.includes(file.type)) {
      setError('Выберите изображение JPEG или PNG.')
      return
    }
    if (file.size > maxImageSize) {
      setError('Файл слишком большой. Максимальный размер — 12 МБ.')
      return
    }
    const objectURL = URL.createObjectURL(file)
    setImage(file)
    setPreview(objectURL)
    setError(null)
  }

  function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault()
    chooseImage(event.dataTransfer.files[0])
  }

  async function save(isPublished: boolean) {
    if (!input.title.trim()) {
      setError('Введите название работы.')
      return
    }
    if (!state?.artwork && !image) {
      setError('Добавьте изображение работы.')
      return
    }
    setSaving(true)
    setError(null)
    try {
      if (state?.artwork) {
        await updateArtwork(state.artwork.id, input, image, isPublished)
        const items = await getAdminArtworks()
        const updated = items.find((item) => item.id === state.artwork?.id)
        if (!updated) throw new ArtworkServiceError('Не удалось получить обновлённую работу.')
        onSaved(updated)
      } else if (image) {
        onSaved(await createArtwork(input, image, isPublished))
      }
    } catch (caughtError) {
      handleServiceError(caughtError, onSessionExpired, setError)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={Boolean(state)} onClose={saving ? undefined : onClose} fullWidth maxWidth="md">
      <DialogTitle sx={{ fontFamily: 'Georgia, serif', fontSize: '2rem', pt: 4, px: { xs: 3, sm: 5 } }}>
        {state?.artwork ? 'Редактировать работу' : 'Новая работа'}
      </DialogTitle>
      <DialogContent sx={{ px: { xs: 3, sm: 5 }, pb: 2 }}>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'minmax(240px, .8fr) 1.2fr' }, gap: 4, pt: 1 }}>
          <Stack spacing={2}>
            <Box
              onDragOver={(event) => event.preventDefault()}
              onDrop={handleDrop}
              sx={{ aspectRatio: '4 / 5', border: '1px dashed', borderColor: 'divider', bgcolor: 'background.default', display: 'grid', placeItems: 'center', overflow: 'hidden', position: 'relative' }}
            >
              {preview ? (
                <Box component="img" src={preview} alt="Предпросмотр работы" sx={{ width: '100%', height: '100%', objectFit: 'contain' }} />
              ) : (
                <Stack spacing={1.5} sx={{ alignItems: 'center', textAlign: 'center', px: 3 }}>
                  <UploadFileOutlinedIcon color="disabled" sx={{ fontSize: 42 }} />
                  <Typography>Перетащите изображение сюда</Typography>
                  <Typography variant="body2" color="text.secondary">JPEG или PNG, до 12 МБ</Typography>
                </Stack>
              )}
            </Box>
            <Button component="label" variant="outlined" startIcon={<UploadFileOutlinedIcon />} disabled={saving}>
              {preview ? 'Заменить изображение' : 'Выбрать изображение'}
              <input hidden type="file" accept="image/jpeg,image/png" onChange={(event) => chooseImage(event.target.files?.[0])} />
            </Button>
            {image && <Typography variant="caption" color="text.secondary">{image.name} · {formatFileSize(image.size)}</Typography>}
          </Stack>

          <Stack spacing={2.5}>
            <TextField required label="Название" value={input.title} onChange={(event) => setField('title', event.target.value)} slotProps={{ htmlInput: { maxLength: 256 } }} />
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <TextField label="Год" type="number" value={input.year} onChange={(event) => setField('year', event.target.value)} slotProps={{ htmlInput: { min: 0, max: 9999 } }} sx={{ width: { sm: 150 } }} />
              <TextField fullWidth label="Техника" value={input.technique} onChange={(event) => setField('technique', event.target.value)} slotProps={{ htmlInput: { maxLength: 256 } }} placeholder="Холст, масло" />
            </Stack>
            <TextField label="Описание" value={input.description} onChange={(event) => setField('description', event.target.value)} multiline minRows={4} />
            <TextField label="Альтернативный текст" value={input.imageAlt} onChange={(event) => setField('imageAlt', event.target.value)} slotProps={{ htmlInput: { maxLength: 256 } }} helperText="Если оставить пустым, посетителям будет доступно название работы." />
          </Stack>
        </Box>
        {error && <Alert severity="error" variant="outlined" sx={{ mt: 3 }}>{error}</Alert>}
      </DialogContent>
      <DialogActions sx={{ px: { xs: 3, sm: 5 }, py: 4, flexWrap: 'wrap', gap: 1.5 }}>
        <Button onClick={onClose} disabled={saving} color="inherit">Отмена</Button>
        <Button onClick={() => void save(false)} disabled={saving} variant="outlined">Сохранить как черновик</Button>
        <Button onClick={() => void save(true)} disabled={saving} variant="contained" sx={{ boxShadow: 'none' }}>
          {saving ? <CircularProgress size={20} color="inherit" /> : 'Опубликовать'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

function handleServiceError(error: unknown, onSessionExpired: (message: string) => void, setError: (message: string) => void) {
  if (error instanceof ArtworkServiceError && error.unauthorized) {
    onSessionExpired(error.message)
    return
  }
  setError(error instanceof ArtworkServiceError ? error.message : 'Произошла неизвестная ошибка.')
}

function workCountLabel(count: number): string {
  const mod100 = count % 100
  const mod10 = count % 10
  if (mod100 >= 11 && mod100 <= 14) return 'работ'
  if (mod10 === 1) return 'работа'
  if (mod10 >= 2 && mod10 <= 4) return 'работы'
  return 'работ'
}

function formatFileSize(bytes: number): string {
  return `${(bytes / 1024 / 1024).toFixed(1)} МБ`
}
