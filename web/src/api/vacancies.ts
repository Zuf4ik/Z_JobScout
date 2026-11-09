export interface Vacancy {
  id: string
  source: string
  title: string
  company_name: string
  url: string
  candidate_required_location: string
  description: string
  salary: string
  experience_years: number | null
}

export interface VacanciesResponse {
  items: Vacancy[]
  loaded: number
  text_matches: number
  total: number
}

export interface SearchParams {
  query: string
  maxExperience: number
  signal?: AbortSignal
}

export function vacancyKey(job: Pick<Vacancy, 'source' | 'id'>): string {
  return `${job.source}:${job.id}`
}

export async function searchVacancies({
  query,
  maxExperience,
  signal,
}: SearchParams): Promise<VacanciesResponse> {
  const params = new URLSearchParams({
    query: query.trim(),
    max_experience: String(maxExperience),
  })
  const response = await fetch(`/api/v1/vacancies?${params}`, { signal })
  if (!response.ok) {
    let message = `Ошибка сервера: HTTP ${response.status}`
    try {
      const body: unknown = await response.json()
      if (body && typeof body === 'object' && 'error' in body && typeof body.error === 'string') {
        message = body.error
      }
    } catch { /* HTTP response is not JSON */ }
    throw new Error(message)
  }
  const data: unknown = await response.json()
  if (!data || typeof data !== 'object' || !('items' in data) || !Array.isArray(data.items)) {
    throw new Error('Сервер вернул неожиданный формат данных')
  }
  const parsed = data as VacanciesResponse
  return { ...parsed, items: parsed.items ?? [] }
}
