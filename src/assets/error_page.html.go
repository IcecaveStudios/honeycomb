package assets

const Asset_error_page_html = `
<DOCTYPE html>
<html>
    <head>
        <title>{{ .Code }} {{ .Text }}</title>
        <link href="https://fonts.googleapis.com/css?family=Yantramanav:100" rel="stylesheet">

        <style>
            body {
                background: white;
                font-family: 'Yantramanav', sans-serif;
                font-weight: 100;
                font-size: 3.8vmin;
                border: 0;
                margin: 0;
            }

            #band {
                text-align: center;
                height: 36vh;
                position: absolute;
                top: 0;
                bottom: 0;
                left: 0;
                right: 0;
                margin: auto;
            }

            #code {
                display: inline;
                font-size: 20vmin;
                margin: 1.5vmin;
                padding: 0;
            }

            #text {
                display: inline;
                font-size: 6vmin;
                font-weight: inherit;
                margin: 1.5vmin;
                padding: 0;
                text-transform: uppercase;
                white-space: nowrap;
            }

            p {
                margin: 0 5vmin;
                padding: 0;
                color: gray;
            }
        </style>
    </head>
    <body>
        <div id="band">
            <header>
                <div id="code">{{ .Code }}</div>
                <h1 id="text">{{ .Text }}</h1>
            </header>
            <p>{{ .Message }}</p>
        </div>
    </body>
</html>
`
