export const SERIES_COLORS = [
  'hsl(172, 66%, 45%)', // Teal (primary)
  'hsl(265, 55%, 55%)', // Violet
  'hsl(195, 70%, 48%)', // Cyan
  'hsl(235, 50%, 55%)', // Indigo
  'hsl(160, 50%, 45%)', // Seafoam
  'hsl(280, 45%, 55%)', // Purple
  'hsl(205, 65%, 52%)', // Sky
  'hsl(250, 45%, 58%)', // Periwinkle
  'hsl(185, 55%, 42%)', // Dark Cyan
  'hsl(220, 55%, 50%)', // Blue
  'hsl(215, 15%, 50%)', // Slate (for "Other")
] as const

export function getSeriesColor(index: number, key: string): string {
  // Special case: "Other" always gets the gray color
  if (key === 'Other') {
    return SERIES_COLORS[SERIES_COLORS.length - 1]
  }
  return SERIES_COLORS[index % (SERIES_COLORS.length - 1)]
}
