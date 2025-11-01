package main

import "strings"

func ParseCommentPoints(commentText string) []string {
	var points []string

	// Разделяем текст на части по блокам details
	detailsParts := strings.Split(commentText, "<details>")

	for _, part := range detailsParts {
		if strings.Contains(part, "</details>") {
			// Извлекаем содержимое между <summary> и </summary>
			summaryStart := strings.Index(part, "<summary>")
			summaryEnd := strings.Index(part, "</summary>")

			if summaryStart != -1 && summaryEnd != -1 && summaryEnd > summaryStart {
				summaryContent := part[summaryStart+9 : summaryEnd] // +9 чтобы пропустить "<summary>"

				// Очищаем HTML теги и лишние пробелы
				cleanedPoint := cleanSummaryText(summaryContent)
				if cleanedPoint != "" {
					points = append(points, cleanedPoint)
				}
			}
		}
	}

	return points
}

func cleanSummaryText(text string) string {
	// Удаляем HTML теги
	text = removeHTMLTags(text)

	// Удаляем лишние пробелы и переносы строк
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.TrimSpace(text)

	return text
}

func removeHTMLTags(text string) string {
	// Удаляем все HTML теги
	var result strings.Builder
	inTag := false

	for _, char := range text {
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(char)
		}
	}

	return result.String()
}
