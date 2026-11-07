package httpapp

import "html/template"

var homeTemplate = template.Must(template.New("home").Funcs(template.FuncMap{
	"inc": func(value int) int { return value + 1 },
	"dec": func(value int) int { return value - 1 },
	"urlquery": func(value string) string { return template.URLQueryEscaper(value) },
}).Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="Browse quotes and jokes from five collections through a public API.">
    <meta name="theme-color" content="#101720">
    <title>Quotes API</title>
    <link rel="manifest" href="/manifest.webmanifest">
    <style>
      :root { color-scheme: light dark; --bg: #101720; --panel: #182430; --text: #f1f5f9; --muted: #b5c2cf; --accent: #f3a66b; }
      :root:has(body[data-theme="light"]) { color-scheme: light; --bg: #f4f7fa; --panel: #fff; --text: #17212b; --muted: #526171; --accent: #9a4814; }
      * { box-sizing: border-box; }
      body { margin: 0; background: var(--bg); color: var(--text); font: 1rem/1.6 system-ui, sans-serif; }
      a { color: var(--accent); }
      .wrap { width: min(70rem, 100% - 2rem); margin: 0 auto; }
      header { padding: 3rem 0 2rem; border-bottom: 1px solid #7890a033; }
      h1 { margin: 0; font-size: clamp(2.4rem, 7vw, 4.5rem); letter-spacing: -.05em; }
      header p, .muted { color: var(--muted); }
      main { padding: 2rem 0; }
      form, .quote, .panel { background: var(--panel); border: 1px solid #7890a033; border-radius: 1rem; padding: 1.25rem; }
      form { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: .8rem; align-items: end; }
      label { display: grid; gap: .25rem; font-weight: 650; }
      input, select, button { min-height: 2.75rem; padding: .5rem .7rem; border: 1px solid #7890a088; border-radius: .5rem; font: inherit; }
      button { background: var(--accent); border-color: var(--accent); color: #111; cursor: pointer; font-weight: 700; }
      .grid { display: grid; gap: 1rem; margin: 1.5rem 0; }
      blockquote { margin: 0; font-size: 1.2rem; }
      .quote footer { margin-top: .6rem; color: var(--muted); }
      nav { display: flex; flex-wrap: wrap; gap: 1rem; padding: 1rem 0; }
      .pager { display: flex; justify-content: space-between; }
      :focus-visible { outline: 3px solid var(--accent); outline-offset: 3px; }
      @media (prefers-color-scheme: light) { :root:has(body[data-theme="auto"]) { color-scheme: light; --bg: #f4f7fa; --panel: #fff; --text: #17212b; --muted: #526171; --accent: #9a4814; } }
    </style>
  </head>
  <body data-theme="{{.theme}}">
    <header><div class="wrap"><p>Five collections · Public read-only API</p><h1>Quotes</h1><p>Find a thought, a laugh, or a little inspiration.</p></div></header>
    <main class="wrap">
      <nav aria-label="Quote actions"><a href="/random?category={{urlquery .category}}&amp;theme={{.theme}}">Get a random quote{{if .category}} from {{.category}}{{end}}</a></nav>
      <form method="get" action="/">
        <label for="category">Collection
          <select id="category" name="category"><option value="">All collections</option>{{range $name, $count := .categories}}<option value="{{$name}}" {{if eq $.category $name}}selected{{end}}>{{$name}} ({{$count}})</option>{{end}}</select>
        </label>
        <label for="author">Author or character<input id="author" name="author" value="{{.author}}" maxlength="200"></label>
        <label for="theme">Theme<select id="theme" name="theme"><option value="auto" {{if eq .theme "auto"}}selected{{end}}>System</option><option value="light" {{if eq .theme "light"}}selected{{end}}>Light</option><option value="dark" {{if eq .theme "dark"}}selected{{end}}>Dark</option></select></label>
        <button type="submit">Filter quotes</button>
      </form>
      <p class="muted" aria-live="polite">{{.total}} results · Page {{.page}} of {{.pages}}</p>
      <section class="grid" aria-label="Quotes">
        {{range .quotes}}<article class="quote"><blockquote>“{{.Quote}}”</blockquote><footer>{{if .Author}}{{.Author}} · {{end}}{{.Category}} · <a href="/quotes/{{.ID}}?theme={{$.theme}}">Quote #{{.ID}}</a></footer></article>{{else}}<p class="panel">No quotes match those filters.</p>{{end}}
      </section>
      <nav class="pager" aria-label="Pagination">{{if gt .page 1}}<a href="/?page={{dec .page}}&amp;category={{urlquery .category}}&amp;author={{urlquery .author}}&amp;theme={{.theme}}">Previous page</a>{{else}}<span></span>{{end}}{{if lt .page .pages}}<a href="/?page={{inc .page}}&amp;category={{urlquery .category}}&amp;author={{urlquery .author}}&amp;theme={{.theme}}">Next page</a>{{end}}</nav>
      <section class="panel"><h2>Use the API</h2><p><code>GET /api/v1/quotes/random?category=programming</code></p><p><a href="/server/docs/swagger">OpenAPI reference</a> · <a href="/server/docs/graphql">GraphQL explorer</a> · <a href="/server/help">Help</a> · <a href="/server/about">About</a> · <a href="/server/privacy">Privacy</a> · <a href="/server/terms">Terms</a></p></section>
    </main>
    <script>if ('serviceWorker' in navigator) navigator.serviceWorker.register('/sw.js');</script>
  </body>
</html>`))

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
  <head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>{{.title}} · Quotes API</title><style>body{max-width:48rem;margin:3rem auto;padding:0 1rem;font:1rem/1.6 system-ui,sans-serif}a{color:#9a4814}</style></head>
  <body><main><p><a href="/">Quotes home</a></p><h1>{{.title}}</h1><p>{{.content}}</p></main></body>
</html>`))

var quoteTemplate = template.Must(template.New("quote").Funcs(template.FuncMap{
	"urlquery": func(value string) string { return template.URLQueryEscaper(value) },
}).Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="theme-color" content="#101720">
    <title>{{.title}} · Quotes</title>
    <style>
      :root { color-scheme: light dark; --bg: #101720; --panel: #182430; --text: #f1f5f9; --muted: #b5c2cf; --accent: #f3a66b; }
      :root:has(body[data-theme="light"]) { color-scheme: light; --bg: #f4f7fa; --panel: #fff; --text: #17212b; --muted: #526171; --accent: #9a4814; }
      body { margin: 0; background: var(--bg); color: var(--text); font: 1rem/1.6 system-ui, sans-serif; }
      main { width: min(48rem, 100% - 2rem); margin: 3rem auto; }
      form { display: flex; flex-wrap: wrap; align-items: center; gap: .75rem; margin: 1rem 0; }
      select, button { min-height: 2.75rem; padding: .5rem .7rem; border: 1px solid #7890a088; border-radius: .5rem; font: inherit; }
      button { background: var(--accent); border-color: var(--accent); color: #111; cursor: pointer; font-weight: 700; }
      a { color: var(--accent); }
      nav { display: flex; flex-wrap: wrap; gap: 1rem; margin: 1rem 0; }
      article { padding: clamp(1.5rem, 6vw, 3rem); border: 1px solid #7890a033; border-radius: 1rem; background: var(--panel); }
      blockquote { margin: 0; font-size: clamp(1.4rem, 4vw, 2rem); }
      footer { margin-top: 1.25rem; color: var(--muted); }
      :focus-visible { outline: 3px solid var(--accent); outline-offset: 3px; }
      @media (prefers-color-scheme: light) { :root:has(body[data-theme="auto"]) { color-scheme: light; --bg: #f4f7fa; --panel: #fff; --text: #17212b; --muted: #526171; --accent: #9a4814; } }
    </style>
  </head>
  <body data-theme="{{.theme}}">
    <main>
      <p><a href="/">Browse all quotes</a></p>
      <h1>{{.title}}</h1>
      <article>
        <blockquote>“{{.quote.Quote}}”</blockquote>
        <footer>{{if .quote.Author}}{{.quote.Author}} · {{end}}{{.quote.Category}} · Quote #{{.quote.ID}}</footer>
      </article>
      <form method="get">
        {{if eq .title "Random quote"}}<input type="hidden" name="category" value="{{.randomCategory}}">{{end}}
        <label for="theme">Theme</label>
        <select id="theme" name="theme"><option value="auto" {{if eq .theme "auto"}}selected{{end}}>System</option><option value="light" {{if eq .theme "light"}}selected{{end}}>Light</option><option value="dark" {{if eq .theme "dark"}}selected{{end}}>Dark</option></select>
        <button type="submit">Apply theme</button>
      </form>
      <nav aria-label="Quote actions">
        <a href="/random?category={{urlquery .quote.Category}}&amp;theme={{.theme}}">Another {{.quote.Category}} quote</a>
        <a href="/random?theme={{.theme}}">Random quote from any collection</a>
        <a href="/quotes/{{.quote.ID}}?theme={{.theme}}">Link to this quote</a>
      </nav>
    </main>
  </body>
</html>`))

var healthTemplate = template.Must(template.New("health").Parse(`<!doctype html>
<html lang="en">
  <head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Quotes - Health Status</title><style>body{font:1rem/1.6 system-ui,sans-serif;max-width:48rem;margin:3rem auto;padding:0 1rem}.healthy{color:#167343}</style></head>
  <body><main><p><a href="/">Quotes home</a></p><h1>Quotes health</h1><p class="healthy">Status: {{.Status}}</p><dl><dt>Version</dt><dd>{{.Version}}</dd><dt>Go</dt><dd>{{.GoVersion}}</dd><dt>Uptime</dt><dd>{{.Uptime}}</dd><dt>Mode</dt><dd>{{.Mode}}</dd><dt>Quote count</dt><dd>{{.Total}}</dd><dt>Checked</dt><dd><time datetime="{{.Timestamp}}">{{.Timestamp}}</time></dd></dl><h2>Collections</h2><ul>{{range $name, $count := .Features}}<li>{{$name}}: {{$count}}</li>{{end}}</ul></main></body>
</html>`))
