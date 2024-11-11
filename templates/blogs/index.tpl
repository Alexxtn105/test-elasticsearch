{{ define "blogs/index.tpl"}}
    {{ template "layouts/header.tpl"}}
        <div>
            <h1>Blogs</h1>
            <br/>
            <br/>

            <div class="search-form">
                <form action="/blogs" method="GET">
                    <input type="text" name="query" placeholder="">
                    <input type="submit" value="Search">
                 </form>
            </div>


            <ul>
                {{ with .blogs}}
                    {{ range .}}
                        <li>
                            <div>
                                <a href="/blogs/{{.ID}}">
                                    <h5>{{ .Title}}</h5>
                                </a>
                                <p>{{ .Content}}</p>
                            </div>
                            <br/>
                        </li>
                    {{end}}
                {{end}}
            </ul>
        </div>
    {{ template "layouts/footer.tpl" .}}
{{end}}