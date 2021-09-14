package TemplateGlobal

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"html"
	"html/template"
	"net/url"
	"strings"
)

func GetUIFunctions() map[string]interface{} {
	return template.FuncMap{
		"round2Euro": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.2f €", amount)
		},
		"round0Euro": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.0f €", amount)
		},
		"round0": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.0f", amount)
		},
		"round2": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.2f", amount)
		},
		"round2Percent": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.2f %s", amount*100, "%")
		},
		"round0Percent": func(amount float64) string {
			if amount == 0 {
				return "-"
			}
			p := message.NewPrinter(language.French)
			return p.Sprintf("%.0f %s", amount*100, "%")
		},
		"plus1": func(amount int) int {
			return amount + 1
		},
		"inparenthesis": func(value interface{}) string {
			p := message.NewPrinter(language.French)
			temp := p.Sprintf("(%v)", value)
			if temp != "()" {
				return temp
			}
			return ""
		},
		"color": func(amount float64) string {
			if amount > 0 {
				return "price_green"
			}
			return "price"
		},
		"IsInInt": func(value int, array []int) bool {
			for _, v := range array {
				if v == value {
					return true
				}
			}
			return false
		},
		"Array": func(values ...interface{}) []interface{} {
			return values
		},
		"Split": func(value string, seperator string) []string {
			return strings.Split(value, seperator)
		},
		"decode": func(value string) string {
			decodedValue, err := url.QueryUnescape(value)
			if err != nil {
				return value
			}
			return decodedValue
		},
		"decode2": func(value string) string {
			decodedValue := html.UnescapeString(value)
			return decodedValue
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"Concat2": func(s1 string, s2 string) string {
			return s1 + s2
		},
		"ColumnFloat": func(tag string, position string, color string, class string, hidden bool, displayType string, style string, value float64) template.HTML {
			result := "<" + tag + " "
			if position != "" {
				result += "align=\"" + position + "\" "
			}
			result += " class=\""
			switch color {
			case "red":
				result += "price"
				break
			case "green":
				result += "price_green"
				break
			case "auto":
				if value >= 0 {
					result += "price_green"
				} else {
					result += "price"
				}
			}
			result += class + " \" "
			if hidden {
				result += " hidden "
			}
			result += " style=\"" + style + ";white-space: nowrap;\" "
			result += ">"
			p := message.NewPrinter(language.French)
			if value == 0 {
				result += "-"
			} else {
				switch displayType {
				case "round2Euro":
					result += p.Sprintf("%.2f €", value)
				case "round0Euro":
					result += p.Sprintf("%.0f €", value)
				case "round0":
					result += p.Sprintf("%.0f", value)
				case "round2":
					result += p.Sprintf("%.2f", value)
				case "round2Percent":
					result += p.Sprintf("%.2f %s", value*100, "%")
				case "round0Percent":
					result += p.Sprintf("%.0f %s", value*100, "%")
				}

			}
			result += "</" + tag + ">"
			return template.HTML(result)
		},
		"ColumnString": func(tag string, position string, color, class string, hidden bool, width string, value string, style string) template.HTML {
			result := "<" + tag + " "
			if position != "" {
				result += "align=\"" + position + "\" "
			}
			result += " class=\""
			switch color {
			case "red":
				result += "price"
				break
			case "green":
				result += "price_green"
				break
			}
			result += class + " \" "
			if hidden {
				result += " hidden "
			}
			result += " style=\"" + style + "\""
			if width != "" {
				result += " width=\"" + width + "\" "
			}
			result += "><span>" + value + "</span>"
			result += "</" + tag + ">"
			return template.HTML(result)
		},
	}
}
