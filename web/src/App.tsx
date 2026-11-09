
import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import {
  searchVacancies,
  type Vacancy,
  type VacanciesResponse,
} from './api/vacancies'
import './App.css'

function App() {
  const [query, setQuery] = useState('Golang')
  const [maxExperience, setMaxExperience] = useState(-1)
  const [vacancies, setVacancies] = useState<Vacancy[]>([])
  const [stats, setStats] = useState<VacanciesResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [searched, setSearched] = useState(false)

  const controllerRef = useRef<AbortController | null>(null)

  useEffect(() => {
    return () => controllerRef.current?.abort()
  }, [])

  async function handleSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    // Отменяем предыдущий поиск, если он ещё выполняется.
    controllerRef.current?.abort()

    const searchQuery = query.trim()

    if (!searchQuery) {
      setError('Введите название профессии или технологию')
      setVacancies([])
      setStats(null)
      setSearched(false)
      setLoading(false)
      return
    }

    const controller = new AbortController()
    controllerRef.current = controller

    setLoading(true)
    setError('')
    setVacancies([])
    setStats(null)
    setSearched(true)

    try {
      const result = await searchVacancies({
        query: searchQuery,
        maxExperience,
        signal: controller.signal,
      })

      if (controller.signal.aborted) return

      setVacancies(result.items)
      setStats(result)
    } catch (err) {
      if (controller.signal.aborted) return

      setError(
        err instanceof Error
          ? err.message
          : 'Не удалось получить вакансии'
      )
    } finally {
      if (controllerRef.current === controller) {
        setLoading(false)
        controllerRef.current = null
      }
    }
  }

  return (
    <div className="app">
      <header className="header">
        <div className="brand">
          <div className="brand-icon">J</div>
          <div>
            <strong>JobScout</strong>
            <span>Умный поиск работы</span>
          </div>
        </div>

        <div className="header-label">Go + React</div>
      </header>

      <main className="container">
        <section className="hero">
          <span className="eyebrow">ПОИСК ВАКАНСИЙ</span>
          <h1>
            Найди работу,
            <br />
            которая подходит <em>тебе.</em>
          </h1>
          <p>
            Ищи вакансии по технологиям, профессии и опыту.
            JobScout собирает предложения и помогает
            находить интересные возможности.
          </p>
        </section>

        <form className="search-panel" onSubmit={handleSearch}>
          <label htmlFor="query">Ключевые слова</label>
          <input
            id="query"
            type="search"
            placeholder="Например, Golang или Python"
            value={query}
            maxLength={200}
            onChange={(event) => setQuery(event.target.value)}
          />

          <div className="search-bottom">
            <div className="filter">
              <label htmlFor="experience">Требуемый опыт</label>
              <select
                id="experience"
                value={maxExperience}
                onChange={(event) =>
                  setMaxExperience(Number(event.target.value))
                }
              >
                <option value={-1}>Любой опыт</option>
                <option value={0}>Без опыта</option>
                <option value={1}>До 1 года</option>
                <option value={2}>До 2 лет</option>
                <option value={3}>До 3 лет</option>
              </select>
            </div>

            <button
              className="search-button"
              type="submit"
              disabled={loading}
            >
              {loading ? 'Поиск...' : 'Найти вакансии'}
            </button>
          </div>
        </form>

        <section className="results">
          <div className="results-heading">
            <h2>Результаты поиска</h2>
            {stats && <span>Найдено: {stats.total}</span>}
          </div>

          {stats && (
            <div className="stats">
              Загружено: {stats.loaded} · Совпадений по тексту:{' '}
              {stats.text_matches} · После фильтрации: {stats.total}
            </div>
          )}

          {loading && (
            <div className="message" role="status">
              Получаем вакансии. Это может занять некоторое время...
            </div>
          )}

          {error && (
            <div className="message error" role="alert">
              {error}
            </div>
          )}

          {!loading && !error && searched &&
            vacancies.length === 0 && (
              <div className="message">
                По вашему запросу вакансий не найдено.
                Попробуйте изменить ключевые слова или фильтр.
              </div>
            )}

          {!loading && !error && !searched && (
            <div className="message">
              Введите запрос и нажмите «Найти вакансии».
            </div>
          )}

          {!loading && !error && vacancies.map((job) => (
            <article
              className="vacancy-card"
              key={`${job.source}:${job.id}`}
            >
              <div className="vacancy-top">
                <div>
                  <span className="company">{job.company_name || 'Компания не указана'}</span>
                  <h3>{job.title}</h3>
                </div>
                <span className="source">{job.source}</span>
              </div>

              <div className="vacancy-details">
                <span>
                  📍 {job.candidate_required_location || 'Регион не указан'}
                </span>
                <span>
                  Опыт: {job.experience_years === null
                    ? 'не указан'
                    : `${job.experience_years} лет`}
                </span>
              </div>

              <div className="vacancy-footer">
                <strong>{job.salary || 'Зарплата не указана'}</strong>
                {job.url ? (
                  <a
                    href={job.url}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Открыть вакансию ↗
                  </a>
                ) : (
                  <span>Ссылка отсутствует</span>
                )}
              </div>
            </article>
          ))}
        </section>
      </main>
    </div>
  )
}

export default App
