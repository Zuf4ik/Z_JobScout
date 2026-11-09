import { useEffect, useMemo, useRef, useState } from 'react'
import type { FormEvent, ReactNode } from 'react'
import {
  ArrowDownAZ, ArrowRight, ArrowUpRight, Bell, Bookmark, BookmarkCheck,
  BriefcaseBusiness, Check, CheckCircle2, ChevronDown, CircleHelp, Clock3, Compass,
  ExternalLink, FileText, Filter, Github, Heart, LayoutDashboard, MapPin, Menu,
  MoreHorizontal, RotateCcw, Search, Settings2, ShieldCheck, SlidersHorizontal,
  Sparkles, TrendingUp, UserRound, Wallet, X, Zap,
} from 'lucide-react'
import { searchVacancies, vacancyKey, type Vacancy, type VacanciesResponse } from './api/vacancies'
import './App.css'

type Page = 'dashboard' | 'vacancies' | 'saved' | 'applications' | 'settings'
type ApplicationStage = 'planned' | 'sent' | 'interview' | 'rejected' | 'offer'
type Application = { job: Vacancy; stage: ApplicationStage; note: string; updatedAt: string }
type SortMode = 'relevant' | 'salary' | 'name'
const STAGES: Record<ApplicationStage, string> = {
  planned: 'Планирую отклик', sent: 'Отклик отправлен', interview: 'Собеседование',
  rejected: 'Отказ', offer: 'Предложение',
}
const PAGE_TITLES: Record<Page, string> = {
  dashboard: 'Обзор', vacancies: 'Поиск вакансий', saved: 'Избранное',
  applications: 'Мои отклики', settings: 'Настройки',
}

function loadStored<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key)
    return raw ? (JSON.parse(raw) as T) : fallback
  } catch { return fallback }
}
function saveStored(key: string, data: unknown) {
  try { localStorage.setItem(key, JSON.stringify(data)) } catch { /* storage unavailable */ }
}
function salaryNumber(text: string): number {
  const matches = text.match(/\d[\d\s\u00a0]*/g)
  if (!matches?.length) return 0
  return Number(matches[0].replace(/\s/g, '')) || 0
}
function shortenDescription(html: string): string {
  if (!html) return 'Описание отсутствует. Перейдите на сайт источника, чтобы узнать подробнее.'
  return html.replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/&amp;/g, '&').replace(/\s+/g, ' ').trim()
}
function meaningfulDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}
function JobMark({ size = 'regular' }: { size?: 'regular' | 'large' }) {
  return <div className={`job-mark ${size === 'large' ? 'job-mark-large' : ''}`} aria-hidden="true"><span>J</span><span className="mark-spark">✦</span></div>
}
function EmptyState({ icon, title, text, action }: { icon: ReactNode; title: string; text: string; action?: ReactNode }) {
  return <div className="empty-state"><div className="empty-icon">{icon}</div><h3>{title}</h3><p>{text}</p>{action}</div>
}
function DetailModal({ job, saved, onClose, onSave, onApply }: {
  job: Vacancy; saved: boolean; onClose: () => void; onSave: () => void; onApply: () => void
}) {
  useEffect(() => {
    function onEscape(event: KeyboardEvent) { if (event.key === 'Escape') onClose() }
    document.addEventListener('keydown', onEscape)
    return () => document.removeEventListener('keydown', onEscape)
  }, [onClose])
  return <div className="modal-backdrop" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose() }}>
    <section className="modal" role="dialog" aria-modal="true" aria-labelledby="job-modal-title">
      <div className="modal-header"><span className="eyebrow">ПОДРОБНОСТИ ВАКАНСИИ</span><button className="icon-button" onClick={onClose} aria-label="Закрыть"><X size={20}/></button></div>
      <div className="modal-content">
        <div className="job-company-avatar"><BriefcaseBusiness size={24}/></div>
        <h2 id="job-modal-title">{job.title}</h2>
        <p className="modal-company">{job.company_name || 'Компания не указана'} · {job.source}</p>
        <div className="detail-pills">
          <span><MapPin size={15}/>{job.candidate_required_location || 'Регион не указан'}</span>
          <span><Wallet size={15}/>{job.salary || 'Зарплата не указана'}</span>
          <span><Clock3 size={15}/>{job.experience_years === null ? 'Опыт не указан' : `Опыт: ${job.experience_years} г.`}</span>
        </div>
        <h3 className="detail-heading">Описание</h3>
        <p className="job-description">{shortenDescription(job.description)}</p>
      </div>
      <div className="modal-actions">
        <button className="button button-ghost" onClick={onSave}>{saved ? <BookmarkCheck size={17}/> : <Bookmark size={17}/>} {saved ? 'В избранном' : 'Сохранить'}</button>
        <button className="button button-outline" onClick={onApply}><CheckCircle2 size={17}/> В мои отклики</button>
        {job.url && <a className="button button-primary" href={job.url} target="_blank" rel="noopener noreferrer">Перейти к вакансии <ExternalLink size={16}/></a>}
      </div>
    </section>
  </div>
}
function JobCard({ job, saved, applied, onSave, onOpen, onApply }: {
  job: Vacancy; saved: boolean; applied: boolean; onSave: () => void; onOpen: () => void; onApply: () => void
}) {
  return <article className="job-card">
    <div className="job-card-main">
      <div className="job-company-avatar"><BriefcaseBusiness size={22}/></div>
      <div className="job-information">
        <div className="job-card-heading"><span className="job-company">{job.company_name || 'Компания не указана'}</span><span className="provider-tag">{job.source}</span></div>
        <button className="job-title-button" onClick={onOpen}><h3>{job.title}</h3></button>
        <div className="job-tags">
          <span><MapPin size={14}/>{job.candidate_required_location || 'Регион не указан'}</span>
          <span><Clock3 size={14}/>{job.experience_years === null ? 'Опыт не указан' : `${job.experience_years} г. опыта`}</span>
        </div>
      </div>
      <button className={`save-button ${saved ? 'is-saved' : ''}`} onClick={onSave} aria-label={saved ? 'Убрать из избранного' : 'Сохранить вакансию'} title={saved ? 'Убрать из избранного' : 'Сохранить'}>
        {saved ? <BookmarkCheck size={19}/> : <Bookmark size={19}/>}
      </button>
    </div>
    <div className="job-card-bottom">
      <strong className="job-salary">{job.salary || 'Зарплата не указана'}</strong>
      <div className="job-card-actions">
        <button className="text-button" onClick={onApply}>{applied ? 'Изменить отклик' : 'В мои отклики'} <ArrowRight size={15}/></button>
        <button className="button button-small button-outline" onClick={onOpen}>Подробнее <ArrowUpRight size={15}/></button>
      </div>
    </div>
  </article>
}

function App() {
  const [page, setPage] = useState<Page>('dashboard')
  const [mobileMenu, setMobileMenu] = useState(false)
  const [query, setQuery] = useState('Golang')
  const [maxExperience, setMaxExperience] = useState(-1)
  const [region, setRegion] = useState('')
  const [salaryOnly, setSalaryOnly] = useState(false)
  const [sortMode, setSortMode] = useState<SortMode>('relevant')
  const [vacancies, setVacancies] = useState<Vacancy[]>([])
  const [stats, setStats] = useState<VacanciesResponse | null>(null)
  const [searched, setSearched] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [saved, setSaved] = useState<Vacancy[]>(() => loadStored('jobscout.saved.v1', []))
  const [applications, setApplications] = useState<Application[]>(() => loadStored('jobscout.applications.v1', []))
  const [profileName, setProfileName] = useState(() => loadStored('jobscout.profileName.v1', ''))
  const [selectedJob, setSelectedJob] = useState<Vacancy | null>(null)
  const [editingApplication, setEditingApplication] = useState<string | null>(null)
  const [draftNote, setDraftNote] = useState('')
  const [draftStage, setDraftStage] = useState<ApplicationStage>('planned')
  const controllerRef = useRef<AbortController | null>(null)
  const searchCounter = useRef(0)

  useEffect(() => saveStored('jobscout.saved.v1', saved), [saved])
  useEffect(() => saveStored('jobscout.applications.v1', applications), [applications])
  useEffect(() => saveStored('jobscout.profileName.v1', profileName), [profileName])
  useEffect(() => () => { controllerRef.current?.abort() }, [])

  const savedKeys = useMemo(() => new Set(saved.map(vacancyKey)), [saved])
  const applicationKeys = useMemo(() => new Set(applications.map((a) => vacancyKey(a.job))), [applications])
  const matchingVacancies = useMemo(() => {
    const filtered = vacancies.filter((job) => {
      if (region.trim() && !job.candidate_required_location.toLocaleLowerCase('ru').includes(region.trim().toLocaleLowerCase('ru'))) return false
      if (salaryOnly && !job.salary.trim()) return false
      return true
    })
    if (sortMode === 'salary') filtered.sort((a, b) => salaryNumber(b.salary) - salaryNumber(a.salary))
    if (sortMode === 'name') filtered.sort((a, b) => a.title.localeCompare(b.title, 'ru'))
    return filtered
  }, [vacancies, region, salaryOnly, sortMode])

  const navigate = (next: Page) => { setPage(next); setMobileMenu(false); window.scrollTo({ top: 0, behavior: 'smooth' }) }
  const toggleSave = (job: Vacancy) => {
    const key = vacancyKey(job)
    setSaved((prev) => prev.some((item) => vacancyKey(item) === key) ? prev.filter((item) => vacancyKey(item) !== key) : [job, ...prev])
  }
  const openApplication = (job: Vacancy) => {
    const entry = applications.find((item) => vacancyKey(item.job) === vacancyKey(job))
    setDraftNote(entry?.note ?? '')
    setDraftStage(entry?.stage ?? 'planned')
    setEditingApplication(vacancyKey(job))
    setSelectedJob(null)
    setApplications((prev) => prev.some((item) => vacancyKey(item.job) === vacancyKey(job)) ? prev : [{ job, stage: 'planned', note: '', updatedAt: new Date().toISOString() }, ...prev])
    navigate('applications')
  }
  const updateApplication = (key: string) => {
    setApplications((prev) => prev.map((item) => vacancyKey(item.job) === key ? { ...item, stage: draftStage, note: draftNote, updatedAt: new Date().toISOString() } : item))
    setEditingApplication(null)
  }
  async function search(event?: FormEvent<HTMLFormElement>, overrideQuery?: string) {
    event?.preventDefault()
    if (!(overrideQuery ?? query).trim()) { setError('Введите название профессии или технологию'); setSearched(false); navigate('vacancies'); return }
    controllerRef.current?.abort()
    const controller = new AbortController()
    controllerRef.current = controller
    const thisSearch = ++searchCounter.current
    setLoading(true); setSearched(true); setError(''); setStats(null); setVacancies([])
    navigate('vacancies')
    try {
      const response = await searchVacancies({ query: overrideQuery ?? query, maxExperience, signal: controller.signal })
      if (thisSearch !== searchCounter.current || controller.signal.aborted) return
      setVacancies(response.items)
      setStats(response)
    } catch (err) {
      if (thisSearch !== searchCounter.current || controller.signal.aborted) return
      setError(err instanceof Error ? err.message : 'Неизвестная ошибка поиска')
    } finally {
      if (thisSearch === searchCounter.current) { setLoading(false); controllerRef.current = null }
    }
  }
  const runQuickSearch = (term: string) => {
    setQuery(term)
    setPage('vacancies')
    void search(undefined, term)
  }
  const cleanFilters = () => { setRegion(''); setSalaryOnly(false); setSortMode('relevant'); setMaxExperience(-1) }
  const renderCard = (job: Vacancy) => <JobCard key={vacancyKey(job)} job={job} saved={savedKeys.has(vacancyKey(job))} applied={applicationKeys.has(vacancyKey(job))} onSave={() => toggleSave(job)} onOpen={() => setSelectedJob(job)} onApply={() => openApplication(job)}/>

  return <div className="app-shell">
    {mobileMenu && <button className="mobile-overlay" aria-label="Закрыть меню" onClick={() => setMobileMenu(false)}/>}
    <aside className={`sidebar ${mobileMenu ? 'sidebar-open' : ''}`}>
      <div className="sidebar-brand"><JobMark/><div><strong>JobScout<span className="brand-period">.</span></strong><small>Workspace</small></div><button className="sidebar-close icon-button" onClick={() => setMobileMenu(false)} aria-label="Закрыть меню"><X size={20}/></button></div>
      <div className="nav-section-title">РАБОЧЕЕ ПРОСТРАНСТВО</div>
      <nav className="sidebar-nav" aria-label="Основная навигация">
        <button className={page === 'dashboard' ? 'nav-item active' : 'nav-item'} onClick={() => navigate('dashboard')}><LayoutDashboard size={18}/>Обзор</button>
        <button className={page === 'vacancies' ? 'nav-item active' : 'nav-item'} onClick={() => navigate('vacancies')}><Search size={18}/>Поиск вакансий</button>
        <button className={page === 'saved' ? 'nav-item active' : 'nav-item'} onClick={() => navigate('saved')}><Bookmark size={18}/>Избранное <span className="nav-count">{saved.length}</span></button>
        <button className={page === 'applications' ? 'nav-item active' : 'nav-item'} onClick={() => navigate('applications')}><BriefcaseBusiness size={18}/>Мои отклики <span className="nav-count">{applications.length}</span></button>
      </nav>
      <div className="nav-section-title second-title">СИСТЕМА</div>
      <nav className="sidebar-nav"><button className={page === 'settings' ? 'nav-item active' : 'nav-item'} onClick={() => navigate('settings')}><Settings2 size={18}/>Настройки</button></nav>
      <div className="sidebar-spacer"/>
      <div className="sidebar-ai"><div className="ai-icon"><Sparkles size={18}/></div><strong>ИИ-помощник</strong><p>Анализ резюме и умные рекомендации появятся в следующих версиях.</p><span>Скоро</span></div>
      <div className="sidebar-footer"><div className="avatar"><UserRound size={18}/></div><div><strong>{profileName || 'Мой профиль'}</strong><small>Локальная версия</small></div><MoreHorizontal size={18}/></div>
    </aside>

    <div className="main-region">
      <header className="topbar"><div className="topbar-title"><button className="mobile-menu-button icon-button" aria-label="Открыть меню" onClick={() => setMobileMenu(true)}><Menu size={21}/></button><span>Рабочее пространство</span><span className="crumb-sep">/</span><strong>{PAGE_TITLES[page]}</strong></div><div className="topbar-right"><span className="api-state"><span className="live-dot"/>Подключён к Go API</span><button className="icon-button top-icon" aria-label="Избранное" onClick={() => navigate('saved')}><Heart size={18}/></button><button className="top-avatar" aria-label="Настройки профиля" onClick={() => navigate('settings')}><UserRound size={17}/></button></div></header>
      <main className="page-content">
        {page === 'dashboard' && <>
          <section className="welcome-row"><div><span className="eyebrow"><Zap size={14}/> ТВОЙ КАРЬЕРНЫЙ НАВИГАТОР</span><h1>Добро пожаловать в <span>JobScout.</span></h1><p>Все вакансии, избранное и отклики — в одном месте. Начни с поиска подходящей работы.</p></div><div className="welcome-graphic"><Compass size={84}/><span className="orb orb-one"/><span className="orb orb-two"/></div></section>
          <div className="stats-grid"><div className="stat-card"><div className="stat-icon blue"><Search size={20}/></div><span>Последний поиск</span><strong>{stats?.loaded ?? '—'}</strong><small>загруженных вакансий</small></div><div className="stat-card"><div className="stat-icon aqua"><CheckCircle2 size={20}/></div><span>Подходящих по фильтрам</span><strong>{stats?.total ?? '—'}</strong><small>в последней выдаче</small></div><div className="stat-card"><div className="stat-icon violet"><Bookmark size={20}/></div><span>В избранном</span><strong>{saved.length}</strong><small>сохранённых вакансий</small></div><div className="stat-card"><div className="stat-icon amber"><BriefcaseBusiness size={20}/></div><span>Мои отклики</span><strong>{applications.length}</strong><small>записей в трекере</small></div></div>
          <div className="dashboard-grid"><section className="panel dashboard-search"><div className="section-heading"><div><span className="eyebrow">НАЧНИ ЗДЕСЬ</span><h2>Быстрый поиск</h2></div><Search size={20} className="muted-icon"/></div><form className="quick-search-form" onSubmit={search}><div className="search-input-wrap"><Search size={19}/><input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Должность или технология" maxLength={200}/></div><button className="button button-primary" disabled={loading} type="submit">Искать <ArrowRight size={17}/></button></form><div className="quick-terms"><span>Популярное:</span>{['Golang', 'Python', 'Backend', 'Java'].map((term) => <button key={term} onClick={() => runQuickSearch(term)}>{term}</button>)}</div></section><section className="panel progress-panel"><div className="section-heading"><div><span className="eyebrow">ТВОЙ ПРОГРЕСС</span><h2>Воронка поиска</h2></div><TrendingUp size={20} className="muted-icon"/></div><div className="progress-row"><span>В избранном</span><strong>{saved.length}</strong></div><div className="progress-track"><span style={{ width: `${Math.min(saved.length * 15, 100)}%` }}/></div><div className="progress-row"><span>Отклик отправлен</span><strong>{applications.filter((a) => a.stage !== 'planned').length}</strong></div><div className="progress-track teal"><span style={{ width: `${Math.min(applications.filter((a) => a.stage !== 'planned').length * 15, 100)}%` }}/></div><p>Отслеживай этапы поиска без таблиц и заметок.</p></section></div>
          <section className="panel saved-preview"><div className="section-heading"><div><span className="eyebrow">ВАЖНОЕ</span><h2>Последние сохранённые вакансии</h2></div><button className="text-button" onClick={() => navigate('saved')}>Все вакансии <ArrowRight size={16}/></button></div>{saved.length ? <div className="compact-jobs">{saved.slice(0, 3).map((job) => <button key={vacancyKey(job)} className="compact-job" onClick={() => setSelectedJob(job)}><span className="small-company-icon"><BriefcaseBusiness size={19}/></span><span><strong>{job.title}</strong><small>{job.company_name || 'Компания не указана'}</small></span><ArrowUpRight size={17}/></button>)}</div> : <EmptyState icon={<Bookmark size={25}/>} title="Пока пусто" text="Сохраняй интересные вакансии, чтобы вернуться к ним позже." action={<button className="button button-outline" onClick={() => navigate('vacancies')}>Найти вакансии</button>}/>}</section>
        </>}

        {page === 'vacancies' && <>
          <div className="page-heading"><div><span className="eyebrow">ВОЗМОЖНОСТИ</span><h1>Поиск вакансий</h1><p>Реальные предложения из подключённых источников</p></div><div className="source-note"><span className="live-dot"/> Источник: Работа России</div></div>
          <form className="main-search" onSubmit={search}><div className="search-input-wrap"><Search size={20}/><input aria-label="Поисковый запрос" value={query} onChange={(event) => setQuery(event.target.value)} maxLength={200} placeholder="Название должности, язык или технология"/></div><button className="button button-primary" type="submit" disabled={loading}>{loading ? <RotateCcw className="spin" size={18}/> : <Search size={18}/>} {loading ? 'Ищем...' : 'Найти вакансии'}</button></form>
          <div className="vacancies-layout"><aside className="filters-panel panel"><div className="filter-heading"><h3><SlidersHorizontal size={18}/> Фильтры</h3><button className="muted-text-button" onClick={cleanFilters}><RotateCcw size={13}/> Сбросить</button></div><label className="field-label" htmlFor="experience">Требуемый опыт</label><div className="select-wrap"><select id="experience" value={maxExperience} onChange={(e) => setMaxExperience(Number(e.target.value))}><option value={-1}>Любой опыт</option><option value={0}>Без опыта</option><option value={1}>До 1 года</option><option value={2}>До 2 лет</option><option value={3}>До 3 лет</option><option value={5}>До 5 лет</option></select><ChevronDown size={16}/></div><p className="filter-hint">Применяется при новом поиске. Вакансии с неизвестным опытом сохраняются.</p><label className="field-label" htmlFor="region">Регион</label><input className="input-field" id="region" placeholder="Например, Москва" value={region} onChange={(e) => setRegion(e.target.value)}/><p className="filter-hint">Фильтр по уже полученным вакансиям.</p><label className="check-field"><input type="checkbox" checked={salaryOnly} onChange={(e) => setSalaryOnly(e.target.checked)}/><span>Только с зарплатой</span></label><div className="filter-footer"><Filter size={17}/><span>Показаны данные только первой страницы источника</span></div></aside>
            <section className="vacancies-results"><div className="results-toolbar"><div><h2>{stats ? `${matchingVacancies.length} вакансий` : 'Результаты поиска'}</h2><p>{stats ? `Загружено ${stats.loaded} · По тексту ${stats.text_matches}` : 'Начни поиск, чтобы увидеть предложения'}</p></div><div className="sort-box"><ArrowDownAZ size={16}/><select value={sortMode} aria-label="Сортировка" onChange={(e) => setSortMode(e.target.value as SortMode)}><option value="relevant">По релевантности</option><option value="salary">По зарплате*</option><option value="name">По названию</option></select></div></div>
              {loading && <div className="loading-grid">{[1, 2, 3].map((i) => <div key={i} className="job-skeleton"><div/><div/><div/><div/></div>)}</div>}
              {error && !loading && <EmptyState icon={<CircleHelp size={27}/>} title="Не удалось загрузить вакансии" text={error} action={<button className="button button-outline" onClick={() => void search()}>Повторить поиск</button>}/>}
              {!loading && !error && !searched && <EmptyState icon={<Search size={28}/>} title="Начни свой поиск" text="Введи профессию или технологию, выбери опыт и нажми «Найти вакансии»."/>}
              {!loading && !error && searched && matchingVacancies.length === 0 && <EmptyState icon={<Compass size={28}/>} title="Совпадений пока нет" text="Попробуй другой запрос или сбрось фильтры. Источник может возвращать ограниченную выборку." action={<button className="button button-outline" onClick={cleanFilters}>Сбросить фильтры</button>}/>}
              {!loading && !error && matchingVacancies.map(renderCard)}
              {stats && !loading && <p className="under-results">* Сортировка по зарплате приблизительная: исходные значения могут быть диапазонами или текстом. Результаты ограничены данными, полученными от источника.</p>}
            </section></div>
        </>}

        {page === 'saved' && <><div className="page-heading"><div><span className="eyebrow">ТВОЯ ПОДБОРКА</span><h1>Избранные вакансии</h1><p>Предложения, к которым ты хочешь вернуться</p></div><span className="counter-pill">{saved.length} сохранено</span></div>{saved.length ? <div className="cards-narrow">{saved.map(renderCard)}</div> : <EmptyState icon={<Bookmark size={28}/>} title="Избранное пока пусто" text="На странице поиска нажми значок закладки рядом с вакансией — она появится здесь." action={<button className="button button-primary" onClick={() => navigate('vacancies')}>Перейти к поиску <ArrowRight size={16}/></button>}/>}</>}

        {page === 'applications' && <><div className="page-heading"><div><span className="eyebrow">КОНТРОЛЬ ПРОЦЕССА</span><h1>Мои отклики</h1><p>Отмечай этапы и веди заметки о вакансиях</p></div><span className="counter-pill">{applications.length} записей</span></div><div className="disclaimer"><ShieldCheck size={19}/><span>Это личный трекер. Изменение статуса здесь <strong>не отправляет резюме работодателю</strong> и не обновляет статус на сайте вакансии.</span></div>{applications.length ? <div className="application-list">{applications.map((item) => { const key = vacancyKey(item.job); const editing = editingApplication === key; return <div className="application-card" key={key}><div className="application-main"><div><span className="job-company">{item.job.company_name || 'Компания не указана'}</span><h3>{item.job.title}</h3><small>Обновлено: {meaningfulDate(item.updatedAt)}</small></div><span className={`stage stage-${item.stage}`}>{STAGES[item.stage]}</span></div>{editing ? <div className="application-editor"><label>Этап<select value={draftStage} onChange={(e) => setDraftStage(e.target.value as ApplicationStage)}>{Object.entries(STAGES).map(([value, title]) => <option key={value} value={value}>{title}</option>)}</select></label><label>Заметка<textarea rows={3} value={draftNote} onChange={(e) => setDraftNote(e.target.value)} placeholder="Например, написал рекрутеру во вторник"/></label><div className="editor-actions"><button className="button button-primary" onClick={() => updateApplication(key)}><Check size={16}/> Сохранить</button><button className="button button-ghost" onClick={() => setEditingApplication(null)}>Отмена</button></div></div> : <><p className="application-note">{item.note || 'Заметок пока нет'}</p><div className="application-actions"><button className="text-button" onClick={() => { setEditingApplication(key); setDraftStage(item.stage); setDraftNote(item.note) }}>Изменить статус / заметку <ArrowRight size={15}/></button><button className="muted-text-button" onClick={() => setApplications((p) => p.filter((a) => vacancyKey(a.job) !== key))}>Удалить</button>{item.job.url && <a className="muted-text-button" href={item.job.url} target="_blank" rel="noopener noreferrer">Оригинал <ExternalLink size={14}/></a>}</div></>}</div> })}</div> : <EmptyState icon={<FileText size={28}/>} title="Пока нет записей" text="Добавляй вакансии в трекер, чтобы не забывать, куда собирался откликнуться." action={<button className="button button-primary" onClick={() => navigate('vacancies')}>Найти вакансии</button>}/>}</>}

        {page === 'settings' && <><div className="page-heading"><div><span className="eyebrow">ПЕРСОНАЛИЗАЦИЯ</span><h1>Настройки</h1><p>Управление локальными данными и профилем</p></div></div><div className="settings-columns"><section className="panel settings-panel"><div className="section-heading"><div><span className="eyebrow">ПРОФИЛЬ</span><h2>Основная информация</h2></div><UserRound size={20}/></div><label className="field-label" htmlFor="user-name">Как к тебе обращаться</label><input className="input-field" id="user-name" placeholder="Например, Зуфар" value={profileName} maxLength={60} onChange={(e) => setProfileName(e.target.value)}/><p className="filter-hint">Имя сохраняется только в твоём браузере.</p><div className="settings-tip"><Bell size={19}/><span>Telegram-уведомления подключим на отдельном этапе разработки.</span></div></section><section className="panel settings-panel"><div className="section-heading"><div><span className="eyebrow">ДАННЫЕ</span><h2>Локальное хранилище</h2></div><ShieldCheck size={21}/></div><p className="settings-text">Избранное и статусы откликов сохраняются через localStorage на этом устройстве. Без серверного аккаунта они не синхронизируются между компьютерами.</p><button className="button button-outline danger-button" onClick={() => { if (window.confirm('Удалить сохранённые вакансии и историю откликов на этом устройстве?')) { setSaved([]); setApplications([]) } }}>Очистить избранное и отклики</button></section></div><div className="about-panel"><JobMark/><div><strong>JobScout MVP</strong><span>Go REST API · React · TypeScript</span></div><a href="https://github.com/Zuf4ik/Z_JobScout" target="_blank" rel="noopener noreferrer"><Github size={18}/> Репозиторий <ExternalLink size={14}/></a></div></>}
      </main>
      <footer className="site-footer"><span>© {new Date().getFullYear()} JobScout · Персональный поиск работы</span><span>Сделано на Go и React <span className="footer-heart">♥</span></span></footer>
    </div>
    {selectedJob && <DetailModal job={selectedJob} saved={savedKeys.has(vacancyKey(selectedJob))} onClose={() => setSelectedJob(null)} onSave={() => toggleSave(selectedJob)} onApply={() => openApplication(selectedJob)}/>}
  </div>
}

export default App
