package template

import "html/template"

// list of embedded template files.
var files = []struct {
	name string
	data string
}{
	{
		name: "index.tmpl",
		data: index,
	}, {
		name: "logs.tmpl",
		data: logs,
	}, {
		name: "stage.tmpl",
		data: stage,
	},
}

// T exposes the embedded templates.
var T *template.Template

func init() {
	T = template.New("_").Funcs(funcMap)
	for _, file := range files {
		T = template.Must(
			T.New(file.name).Parse(file.data),
		)
	}
}

//
// embedded template files.
//

// files/index.tmpl
var index = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="10">
<title>Dashboard</title>
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/style.css">
<link rel="icon" type="image/png" id="favicon" href="/static/favicon.png">
<script src="/static/timeago.js" type="text/javascript"></script>
</head>
<body>

<header class="navbar">
    <div class="logo">
        <svg viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><defs><path d="M12.086 5.814l-.257.258 10.514 10.514C20.856 18.906 20 21.757 20 25c0 9.014 6.618 15 15 15 3.132 0 6.018-.836 8.404-2.353l10.568 10.568C48.497 55.447 39.796 60 30 60 13.434 60 0 46.978 0 30 0 19.903 4.751 11.206 12.086 5.814zm5.002-2.97C20.998 1.015 25.378 0 30 0c16.566 0 30 13.022 30 30 0 4.67-1.016 9.04-2.835 12.923l-9.508-9.509C49.144 31.094 50 28.243 50 25c0-9.014-6.618-15-15-15-3.132 0-6.018.836-8.404 2.353l-9.508-9.508zM35 34c-5.03 0-9-3.591-9-9s3.97-9 9-9c5.03 0 9 3.591 9 9s-3.97 9-9 9z" id="a"></path></defs><use fill="#FFF" xlink:href="#a" fill-rule="evenodd"></use></svg>
    </div>
    <nav class="inline-nav">
        <ul>
            <li><a href="/" class="active">Dashboard</a></li>
            <li><a href="/logs">Logging</a></li>
        </ul>
    </nav>
</header>

<main>
    <section>
        <header>
            <h1>Dashboard</h1>
        </header>
        <article class="cards stages">
            {{ if not .Items }}
            <div class="alert sleeping">
                <p>There is no recent activity to display.</p>
            </div>
            {{ else if .Idle }}
            <div class="alert sleeping">
                <p>This runner is currently idle.</p>
            </div>
            {{ end }}
            {{ range .Items }}
            <a href="/view?id={{ .Stage.ID }}" class="card stage">
                <h2>{{ .Repo.Slug }}</h2>
                <img src="{{ .Build.AuthorAvatar }}" />
                <span class="connector"></span>
                <span class="status {{ .Stage.Status }}"></span>
                <span class="desc">assigned stage <em>{{ .Stage.Name }}</em> for build <em>#{{ .Build.Number }}</em></span>
                <span class="time" datetime="{{ if .Stage.Started }}{{ timestamp .Stage.Started }}{{ else }}{{ timestamp .Stage.Created }}{{ end }}"></span>
            </a>
            {{ end }}
        </article>
    </section>
</main>

<footer></footer>

<script>
timeago.render(document.querySelectorAll('.time'));
</script>
</body>
</html>`

// files/logs.tmpl
var logs = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Dashboard</title>
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/style.css">
<link rel="icon" type="image/png" id="favicon" href="/static/favicon.png">
<script src="/static/timeago.js" type="text/javascript"></script>
</head>
<body>

<header class="navbar">
    <div class="logo">
        <svg viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><defs><path d="M12.086 5.814l-.257.258 10.514 10.514C20.856 18.906 20 21.757 20 25c0 9.014 6.618 15 15 15 3.132 0 6.018-.836 8.404-2.353l10.568 10.568C48.497 55.447 39.796 60 30 60 13.434 60 0 46.978 0 30 0 19.903 4.751 11.206 12.086 5.814zm5.002-2.97C20.998 1.015 25.378 0 30 0c16.566 0 30 13.022 30 30 0 4.67-1.016 9.04-2.835 12.923l-9.508-9.509C49.144 31.094 50 28.243 50 25c0-9.014-6.618-15-15-15-3.132 0-6.018.836-8.404 2.353l-9.508-9.508zM35 34c-5.03 0-9-3.591-9-9s3.97-9 9-9c5.03 0 9 3.591 9 9s-3.97 9-9 9z" id="a"></path></defs><use fill="#FFF" xlink:href="#a" fill-rule="evenodd"></use></svg>
    </div>
    <nav class="inline-nav">
        <ul>
            <li><a href="/">Dashboard</a></li>
            <li><a href="/logs" class="active">Logging</a></li>
        </ul>
    </nav>
</header>

<main>
    <section>
        <header>
            <h1>Recent Logs</h1>
        </header>
        {{ if .Entries }}
        <article class="logs">
            {{ range .Entries }}
            <div class="entry">
                <span class="level {{ .Level }}">{{ .Level }}</span>
                <span class="message">{{ .Message }}</span>
                <span class="fields">
                    {{ range $key, $val := .Data }}
                    <span><em>{{ $key }}</em>{{ $val }}</span>
                    {{ end }}
                </span>
                <span class="time" datetime="{{ timestamp .Unix }}"></span>
            </div>
            {{ end }}
        </article>
        {{ else }}
        <div class="alert sleeping">
            <p>There is no recent log activity to display.</p>
        </div>
        {{ end }}
    </section>
</main>

<footer></footer>

<script>
timeago.render(document.querySelectorAll('.time'));
</script>
</body>
</html>`

// files/stage.tmpl
var stage = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
{{- if not (done .Stage.Status) }}
<meta http-equiv="refresh" content="30">
{{- end }}
<title>Dashboard</title>
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/style.css">
<link rel="icon" type="image/png" id="favicon" href="/static/favicon.png">
<script src="/static/timeago.js" type="text/javascript"></script>
</head>
<body>

<header class="navbar">
    <div class="logo">
        <svg viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><defs><path d="M12.086 5.814l-.257.258 10.514 10.514C20.856 18.906 20 21.757 20 25c0 9.014 6.618 15 15 15 3.132 0 6.018-.836 8.404-2.353l10.568 10.568C48.497 55.447 39.796 60 30 60 13.434 60 0 46.978 0 30 0 19.903 4.751 11.206 12.086 5.814zm5.002-2.97C20.998 1.015 25.378 0 30 0c16.566 0 30 13.022 30 30 0 4.67-1.016 9.04-2.835 12.923l-9.508-9.509C49.144 31.094 50 28.243 50 25c0-9.014-6.618-15-15-15-3.132 0-6.018.836-8.404 2.353l-9.508-9.508zM35 34c-5.03 0-9-3.591-9-9s3.97-9 9-9c5.03 0 9 3.591 9 9s-3.97 9-9 9z" id="a"></path></defs><use fill="#FFF" xlink:href="#a" fill-rule="evenodd"></use></svg>
    </div>
    <nav class="inline-nav">
        <ul>
            <li><a href="/">Dashboard</a></li>
            <li><a href="/logs">Logging</a></li>
        </ul>
    </nav>
</header>

<main>
    <section>
        <nav class="breadcrumb">
            <ul>
                <li><a href="/">Dashboard</a></li>
                <li class="separator"></li>
                <li>{{ .Repo.Slug }}</li>
                <li class="separator"></li>
                <li>#{{ .Build.Number }}</li>
            </ul>
        </nav>
        <header>
            <h1>{{ .Repo.Name }}</h1>
        </header>
        
        <div class="card stage">
            <h2>{{ .Build.Message }}</h2>
            <img src="{{ .Build.AuthorAvatar }}" />
            <span class="connector"></span>
            <span class="status {{ .Stage.Status }}"></span>
            {{ if eq .Build.Event "pull_request" }}
                {{ if eq .Build.Action "synchronized" }}
                <span class="desc">{{ .Build.Author }} synchronized pull request <em>#{{ pr .Build.Ref }}</em></span>
                {{ else }}
                <span class="desc">{{ .Build.Author }} opened pull request <em>#{{ pr .Build.Ref }}</em> to <em>{{ .Build.Target }}</em></span>
                {{ end }}
            {{ else if eq .Build.Event "tag" }}
            <span class="desc">{{ .Build.Author }} created reference <em>{{ tag .Build.Ref }}</em></span>
            {{ else if eq .Build.Event "promote"}}
            <span class="desc">{{ .Build.Author }} promoted build <em>#{{ .Build.Parent }}</em> to <em>{{ .Build.Deploy }}</em></span>
            {{ else if eq .Build.Event "rollback"}}
            <span class="desc">{{ .Build.Author }} reverted <em>{{ .Build.Deploy }}</em> to <em>#{{ .Build.Parent }}</em></span>
            {{ else if eq .Build.Event "cron"}}
            <span class="desc">Executed <em>{{ .Build.Cron }}</em> for branch <em>{{ .Build.Target }}</em></span>
            {{ else }}
            <span class="desc">{{ .Build.Author }} pushed <em>{{ sha .Build.After }}</em> to <em>{{ .Build.Target }}</em></span>
            {{ end }}
            <span class="time" datetime="{{ if .Stage.Started }}{{ timestamp .Stage.Started }}{{ else }}{{ timestamp .Stage.Created }}{{ end }}"></span>
        </div>

        {{ if .Stage.Steps }}
        <div class="card steps">
            <div class="body">
                {{ range .Stage.Steps }}
                <div class="step">
                    <span class="status {{ .Status }}"></span>
                    <span class="name">{{ .Name }}</span>
                    <span class="status-name">{{ .Status }}</span>
                </div>
                {{ end }}
            </div>
        </div>
        {{ end }}

        {{ if .Logs }}
        <div class="logs">
            {{ range .Logs }}
            <div class="entry">
                <span class="level {{ .Level }}">{{ .Level }}</span>
                <span class="message">{{ .Message }}</span>
                <span class="fields">
                    {{ range $key, $val := .Data }}
                    <span><em>{{ $key }}</em>{{ $val }}</span>
                    {{ end }}
                </span>
                <span class="time" datetime="{{ timestamp .Unix }}"></span>
            </div>
            {{ end }}
        </div>
        {{ end }}
    </section>
</main>

<footer></footer>

<script>
timeago.render(document.querySelectorAll('.time'));
</script>
</body>
</html>`
