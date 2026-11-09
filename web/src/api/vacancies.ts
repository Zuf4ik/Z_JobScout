
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

interface SearchParams {
  query: string
  maxExperience: number
  signal?: AbortSignal
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

  const response = await fetch(
    `/api/v1/vacancies?${params.toString()}`,
    { signal }
  )

  if (!response.ok) {
    let message = `Ошибка HTTP: ${response.status}`

    try {
      const body: { error?: string } = await response.json()
      if (body.error) {
        message = body.error
      }
    } catch {
      // Если ответ не JSON, сохраняем сообщение по статусу.
    }

    throw new Error(message)
  }

  const data: VacanciesResponse = await response.json()

  return data
}
