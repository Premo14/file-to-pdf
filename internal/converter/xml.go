package converter

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type ResumeXML struct {
	XMLName    xml.Name      `xml:"resume"`
	Header     HeaderXML     `xml:"header"`
	Skills     SkillsXML     `xml:"skills"`
	Experience ExperienceXML `xml:"experience"`
	Projects   ProjectsXML   `xml:"projects"`
	Education  EducationXML  `xml:"education"`
}

type HeaderXML struct {
	Name    string     `xml:"name"`
	Contact ContactXML `xml:"contact"`
}

type ContactXML struct {
	Email    LinkXML `xml:"email"`
	LinkedIn LinkXML `xml:"linkedin"`
	GitHub   LinkXML `xml:"github"`
	Website  LinkXML `xml:"website"`
}

type LinkXML struct {
	URL  string `xml:"url,attr"`
	Text string `xml:",chardata"`
}

type SkillsXML struct {
	Categories []SkillCategoryXML `xml:"category"`
}

type SkillCategoryXML struct {
	Name  string   `xml:"name,attr"`
	Items []string `xml:"item"`
}

type ExperienceXML struct {
	Jobs []JobXML `xml:"job"`
}

type JobXML struct {
	Title      string        `xml:"title"`
	Date       string        `xml:"date"`
	Role       string        `xml:"role"`
	Highlights HighlightsXML `xml:"highlights"`
}

type HighlightsXML struct {
	Bullets []string `xml:"bullet"`
}

type ProjectsXML struct {
	Projects []ProjectXML `xml:"project"`
}

type ProjectXML struct {
	Title       string `xml:"title"`
	Date        string `xml:"date"`
	Description string `xml:"description"`
	Tools       string `xml:"tools"`
}

type EducationXML struct {
	Institutions []InstitutionXML `xml:"institution"`
}

type InstitutionXML struct {
	Name   string `xml:"name,attr"`
	Degree string `xml:"degree"`
	Date   string `xml:"date"`
}

func XMLToHTML(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}

	var res ResumeXML
	if err := xml.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("failed to unmarshal xml: %v", err)
	}

	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("<h1>%s</h1>\n", res.Header.Name))
	sb.WriteString("<div class=\"contact-info\">\n")
	if res.Header.Contact.Email.Text != "" {
		sb.WriteString(fmt.Sprintf("<span><a href=\"%s\">%s</a></span>\n", res.Header.Contact.Email.URL, res.Header.Contact.Email.Text))
	}
	if res.Header.Contact.LinkedIn.Text != "" {
		sb.WriteString(fmt.Sprintf("<span><a href=\"%s\">%s</a></span>\n", res.Header.Contact.LinkedIn.URL, res.Header.Contact.LinkedIn.Text))
	}
	if res.Header.Contact.GitHub.Text != "" {
		sb.WriteString(fmt.Sprintf("<span><a href=\"%s\">%s</a></span>\n", res.Header.Contact.GitHub.URL, res.Header.Contact.GitHub.Text))
	}
	if res.Header.Contact.Website.Text != "" {
		sb.WriteString(fmt.Sprintf("<span><a href=\"%s\">%s</a></span>\n", res.Header.Contact.Website.URL, res.Header.Contact.Website.Text))
	}
	sb.WriteString("</div>\n")

	// Skills
	if len(res.Skills.Categories) > 0 {
		sb.WriteString("<h2>Skills</h2>\n")
		sb.WriteString("<div class=\"section-item\">\n")
		for _, cat := range res.Skills.Categories {
			sb.WriteString(fmt.Sprintf("<p><strong>%s:</strong> %s</p>\n", cat.Name, strings.Join(cat.Items, ", ")))
		}
		sb.WriteString("</div>\n")
	}

	// Experience
	if len(res.Experience.Jobs) > 0 {
		sb.WriteString("<h2>Professional Experience</h2>\n")
		for _, job := range res.Experience.Jobs {
			sb.WriteString("<div class=\"section-item\">\n")
			sb.WriteString("<div class=\"item-header\">\n")
			sb.WriteString(fmt.Sprintf("<span>%s</span>\n", job.Title))
			sb.WriteString(fmt.Sprintf("<span class=\"date\">%s</span>\n", job.Date))
			sb.WriteString("</div>\n")
			sb.WriteString(fmt.Sprintf("<div class=\"item-subheader\">%s</div>\n", job.Role))
			sb.WriteString("<ul>\n")
			for _, bullet := range job.Highlights.Bullets {
				sb.WriteString(fmt.Sprintf("<li>%s</li>\n", bullet))
			}
			sb.WriteString("</ul>\n")
			sb.WriteString("</div>\n")
		}
	}

	// Projects
	if len(res.Projects.Projects) > 0 {
		sb.WriteString("<h2>Projects</h2>\n")
		for _, proj := range res.Projects.Projects {
			sb.WriteString("<div class=\"section-item\">\n")
			sb.WriteString("<div class=\"item-header\">\n")
			sb.WriteString(fmt.Sprintf("<span>%s</span>\n", proj.Title))
			sb.WriteString(fmt.Sprintf("<span class=\"date\">%s</span>\n", proj.Date))
			sb.WriteString("</div>\n")
			sb.WriteString(fmt.Sprintf("<p>%s</p>\n", proj.Description))
			if proj.Tools != "" {
				sb.WriteString(fmt.Sprintf("<p><strong>Tools:</strong> %s</p>\n", proj.Tools))
			}
			sb.WriteString("</div>\n")
		}
	}

	// Education
	if len(res.Education.Institutions) > 0 {
		sb.WriteString("<h2>Education</h2>\n")
		for _, edu := range res.Education.Institutions {
			sb.WriteString("<div class=\"section-item\">\n")
			sb.WriteString("<div class=\"item-header\">\n")
			sb.WriteString(fmt.Sprintf("<span>%s</span>\n", edu.Name))
			sb.WriteString("</div>\n")
			sb.WriteString("<div class=\"item-header\" style=\"font-weight: normal;\">\n")
			sb.WriteString(fmt.Sprintf("<span>%s</span>\n", edu.Degree))
			sb.WriteString(fmt.Sprintf("<span class=\"date\">%s</span>\n", edu.Date))
			sb.WriteString("</div>\n")
			sb.WriteString("</div>\n")
		}
	}

	return WrapInHTML(sb.String()), nil
}
