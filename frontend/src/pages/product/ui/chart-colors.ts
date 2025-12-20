export const SERIES_COLORS = [
  'hsl(221, 83%, 53%)', // Blue
  'hsl(142, 71%, 45%)', // Green
  'hsl(38, 92%, 50%)', // Orange
  'hsl(0, 84%, 60%)', // Red
  'hsl(262, 83%, 58%)', // Purple
  'hsl(174, 72%, 40%)', // Teal
  'hsl(326, 85%, 55%)', // Pink
  'hsl(43, 96%, 56%)', // Yellow
  'hsl(199, 89%, 48%)', // Cyan
  'hsl(24, 94%, 53%)', // Dark Orange
  'hsl(220, 9%, 46%)', // Gray (for "Other")
] as const

export function getSeriesColor(index: number, key: string): string {
  // Special case: "Other" always gets the gray color
  if (key === 'Other') {
    return SERIES_COLORS[SERIES_COLORS.length - 1]
  }
  return SERIES_COLORS[index % (SERIES_COLORS.length - 1)]
}
