# Artfolio frontend

Одностраничный адаптивный прототип портфолио художницы на React, TypeScript, Vite и Material UI.

## Запуск

```bash
npm install
npm run dev
```

Проверка проекта: `npm run lint` и `npm run build`.

## Где заменять временные данные

- `src/data/artist.ts` — имя, описание, биография, контакты и ссылки.
- `src/data/artworks.ts` — работы, подписи и пути к изображениям.
- `src/theme/designTokens.ts` — палитра, шрифты и основные параметры layout.

Все временные значения помечены комментарием `TEMPORARY`. Отсутствующие изображения работ показываются как нейтральные placeholders, поэтому для наполнения достаточно указать `imageUrl`.

## Подключение API

Компоненты не импортируют mock-данные напрямую. Источник меняется в `src/services/artworksService.ts` и `src/services/artistService.ts`. Оба сервиса уже возвращают `Promise`, поэтому переход на HTTP не потребует изменения интерфейса компонентов.
