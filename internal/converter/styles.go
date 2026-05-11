package converter

const ResumeCSS = `
    * {
        text-align: left !important;
    }
    body {
        font-family: Arial, Helvetica, sans-serif;
        line-height: 1.5;
        color: #000;
        margin: 0;
        padding: 0.75in;
        background-color: #fff;
    }
    @page {
        margin: 0.001in;
    }
    h1 {
        font-size: 26px;
        font-weight: bold;
        margin-bottom: 8px;
        text-transform: uppercase;
        color: #000;
    }
    h2 {
        font-size: 16px;
        font-weight: bold;
        text-transform: uppercase;
        border-bottom: 1px solid #000;
        padding-bottom: 3px;
        margin-top: 25px;
        margin-bottom: 12px;
        color: #000;
    }
    h3 {
        font-size: 14px;
        font-weight: bold;
        margin-top: 15px;
        margin-bottom: 5px;
    }
    .contact-info {
        font-size: 13px;
        margin-bottom: 25px;
        color: #000;
    }
    .contact-info a {
        color: #000;
        text-decoration: underline;
    }
    .contact-info span:not(:last-child)::after {
        content: " | ";
        margin: 0 5px;
    }
    .section-item {
        margin-bottom: 15px;
    }
    .item-header {
        display: flex;
        justify-content: space-between;
        align-items: baseline;
        font-weight: bold;
        font-size: 14px;
    }
    .item-subheader {
        font-style: italic;
        font-size: 13px;
        margin-bottom: 5px;
    }
    .date {
        font-weight: normal;
        font-style: italic;
        font-size: 13px;
        text-align: right !important; /* Dates on the right are standard and ATS friendly */
    }
    ul {
        margin-top: 5px;
        margin-bottom: 10px;
        padding-left: 20px;
    }
    li {
        margin-bottom: 4px;
        font-size: 13px;
    }
    p {
        margin: 5px 0;
        font-size: 13px;
    }
    strong {
        font-weight: bold;
    }
    a {
        color: #0000EE;
        text-decoration: underline;
    }
`

func WrapInHTML(content string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Resume</title>
    <style>` + ResumeCSS + `</style>
</head>
<body>
    ` + content + `
</body>
</html>`
}
