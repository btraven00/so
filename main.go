package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const apiURL = "https://api.stackexchange.com/2.3"

type SearchResponse struct {
	Items []Question `json:"items"`
}

type Question struct {
	Title       string `json:"title"`
	QuestionID  int    `json:"question_id"`
	AnswerCount int    `json:"answer_count"`
	Score       int    `json:"score"`
	Link        string `json:"link"`
	IsAnswered  bool   `json:"is_answered"`
	Body        string `json:"body"`
}

type AnswerResponse struct {
	Items []Answer `json:"items"`
}

type Answer struct {
	Score      int    `json:"score"`
	IsAccepted bool   `json:"is_accepted"`
	Body       string `json:"body"`
}

var (
	verbosity   int
	listResults bool
	ansNum      int
	numResults  int
)

func main() {
	flag.IntVar(&verbosity, "v", 0, "verbosity level (0=snippet only, 1=with context, 2=full)")
	flag.BoolVar(&listResults, "l", false, "list search results with answer counts")
	flag.IntVar(&ansNum, "a", 1, "answer number to show")
	flag.IntVar(&numResults, "n", 10, "number of search results")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: so [-v level] [-l] [-a N] [-n N] <query>")
		os.Exit(1)
	}

	query := strings.Join(flag.Args(), " ")

	results, err := searchStackOverflow(query, numResults)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("no results")
		os.Exit(0)
	}

	if listResults {
		for i, r := range results {
			fmt.Printf("[%d] %s (%d answers)\n", i+1, r.Title, r.AnswerCount)
		}
		return
	}

	// Show answer from first result
	showAnswer(results[0], ansNum)
}

func showAnswer(q Question, answerNum int) {
	if q.AnswerCount == 0 {
		fmt.Println("no answers")
		return
	}

	answers, err := getAnswers(q.QuestionID)
	if err != nil || len(answers) == 0 {
		fmt.Println("no answers")
		return
	}

	if answerNum > len(answers) {
		answerNum = len(answers)
	}

	ans := answers[answerNum-1]

	// Verbosity level 2: full details
	if verbosity >= 2 {
		accepted := ""
		if ans.IsAccepted {
			accepted = " [ACCEPTED]"
		}
		fmt.Printf("Question: %s\n", q.Title)
		fmt.Printf("Score: %d | Answers: %d\n", q.Score, q.AnswerCount)
		fmt.Println(strings.Repeat("-", 70))
		
		qBody := stripHTML(q.Body)
		if len(qBody) > 300 {
			qBody = qBody[:300] + "..."
		}
		fmt.Printf("\n%s\n\n", qBody)
		fmt.Println(strings.Repeat("=", 70))
		fmt.Printf("Answer %d/%d | Score: %d%s\n", answerNum, len(answers), ans.Score, accepted)
		fmt.Println(strings.Repeat("-", 70))
	}

	// Parse answer body
	code, textBefore, textAfter := parseAnswer(ans.Body)

	// Verbosity level 1: show context around snippet
	if verbosity >= 1 && len(textBefore) > 0 {
		fmt.Println(textBefore)
		if len(code) > 0 {
			fmt.Println()
		}
	}

	// Show code snippet (or full text if no code)
	if len(code) > 0 {
		for i, snippet := range code {
			fmt.Println(snippet)
			if i < len(code)-1 {
				fmt.Println()
			}
		}
	} else {
		text := stripHTML(ans.Body)
		if verbosity == 0 && len(text) > 200 {
			text = text[:200] + "..."
		}
		fmt.Println(text)
	}

	// Verbosity level 1: show text after snippet
	if verbosity >= 1 && len(textAfter) > 0 {
		if len(code) > 0 {
			fmt.Println()
		}
		fmt.Println(textAfter)
	}

	// Verbosity level 2: show URL
	if verbosity >= 2 {
		fmt.Printf("\n%s\n", q.Link)
	}
}

func parseAnswer(html string) (codes []string, textBefore string, textAfter string) {
	// Split by code blocks to get surrounding text
	preCodeRe := regexp.MustCompile(`(?s)<pre[^>]*><code[^>]*>(.*?)</code></pre>`)
	
	// Find all code blocks
	matches := preCodeRe.FindAllStringSubmatchIndex(html, -1)
	
	if len(matches) > 0 {
		// Text before first code block
		textBefore = stripHTML(html[:matches[0][0]])
		textBefore = strings.TrimSpace(textBefore)
		if len(textBefore) > 150 {
			// Find last sentence
			sentences := strings.Split(textBefore, ".")
			if len(sentences) > 1 {
				textBefore = sentences[len(sentences)-2] + "."
			} else {
				textBefore = textBefore[len(textBefore)-150:]
			}
		}
		
		// Extract code blocks
		for _, match := range matches {
			code := html[match[2]:match[3]]
			code = cleanCode(code)
			code = strings.TrimSpace(code)
			if code != "" {
				codes = append(codes, code)
			}
		}
		
		// Text after last code block
		lastMatch := matches[len(matches)-1]
		textAfter = stripHTML(html[lastMatch[1]:])
		textAfter = strings.TrimSpace(textAfter)
		if len(textAfter) > 150 {
			// Get first sentence or two
			textAfter = textAfter[:150]
			if idx := strings.LastIndex(textAfter, "."); idx > 0 {
				textAfter = textAfter[:idx+1]
			}
		}
	}
	
	return codes, textBefore, textAfter
}

func cleanCode(code string) string {
	// Decode HTML entities first
	code = strings.ReplaceAll(code, "&lt;", "<")
	code = strings.ReplaceAll(code, "&gt;", ">")
	code = strings.ReplaceAll(code, "&quot;", "\"")
	code = strings.ReplaceAll(code, "&amp;", "&")
	code = strings.ReplaceAll(code, "&#39;", "'")
	code = strings.ReplaceAll(code, "&#x27;", "'")
	
	// Remove language hints like :shell, :python, etc at the start of code blocks
	langHintRe := regexp.MustCompile(`^:\w+\s*\n`)
	code = langHintRe.ReplaceAllString(code, "")
	
	// Remove snippet markers
	code = strings.ReplaceAll(code, "<!-- language: lang-", "")
	code = strings.ReplaceAll(code, "<!-- language-all: lang-", "")
	
	// Remove any remaining HTML tags
	tagRe := regexp.MustCompile(`<[^>]*>`)
	code = tagRe.ReplaceAllString(code, "")
	
	return code
}

func apiGet(endpoint string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var reader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		reader = resp.Body
	}

	return io.ReadAll(reader)
}

func searchStackOverflow(query string, limit int) ([]Question, error) {
	endpoint := fmt.Sprintf("%s/search/advanced?order=desc&sort=relevance&q=%s&site=stackoverflow&filter=withbody",
		apiURL, url.QueryEscape(query))

	body, err := apiGet(endpoint)
	if err != nil {
		return nil, err
	}

	var response SearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Items) > limit {
		response.Items = response.Items[:limit]
	}

	return response.Items, nil
}

func getAnswers(questionID int) ([]Answer, error) {
	endpoint := fmt.Sprintf("%s/questions/%d/answers?order=desc&sort=votes&site=stackoverflow&filter=withbody",
		apiURL, questionID)

	body, err := apiGet(endpoint)
	if err != nil {
		return nil, err
	}

	var response AnswerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "<p>", "\n")
	s = strings.ReplaceAll(s, "</p>", "\n")
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&#x27;", "'")

	// Remove remaining tags
	tagRe := regexp.MustCompile(`<[^>]*>`)
	s = tagRe.ReplaceAllString(s, "")

	// Clean up whitespace
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}
