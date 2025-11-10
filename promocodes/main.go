package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultCount = 3
	minWordLen   = 3
	maxWordLen   = 6
	dictPath     = "/usr/share/dict/words"
)

// model represents the application state
type model struct {
	codes []string
}

// initialModel returns the initial model
func initialModel(codes []string) model {
	return model{
		codes: codes,
	}
}

// Init is called when the program starts
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	var sb strings.Builder
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, code := range m.codes {
		// Generate a random color for each code
		color := randomColor(rng)
		style := lipgloss.NewStyle().Foreground(color)
		sb.WriteString(style.Render(code))
		if i < len(m.codes)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// randomColor generates a random color
func randomColor(rng *rand.Rand) lipgloss.Color {
	r := rng.Intn(256)
	g := rng.Intn(256)
	b := rng.Intn(256)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

// readWords reads and filters words from the dictionary file
func readWords() ([]string, error) {
	file, err := os.Open(dictPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		// Filter out proper nouns (capitalized), too short, or too long words
		if len(word) >= minWordLen && len(word) <= maxWordLen {
			// Check if first character is lowercase (not a proper noun)
			if len(word) > 0 && word[0] >= 'a' && word[0] <= 'z' {
				words = append(words, word)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dictionary file: %w", err)
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("no valid words found in dictionary")
	}

	return words, nil
}

// generatePromoCodes generates unique promo codes
func generatePromoCodes(words []string, count int) ([]string, error) {
	if len(words) < 3 {
		return nil, fmt.Errorf("insufficient words in dictionary (need at least 3)")
	}

	// Calculate maximum possible unique combinations
	maxCombinations := len(words) * len(words) * len(words)
	if count > maxCombinations {
		return nil, fmt.Errorf("requested count (%d) exceeds maximum possible combinations (%d)", count, maxCombinations)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	generated := make(map[string]bool)
	codes := make([]string, 0, count)

	for len(codes) < count {
		// Select 3 random words
		w1 := words[rng.Intn(len(words))]
		w2 := words[rng.Intn(len(words))]
		w3 := words[rng.Intn(len(words))]

		code := fmt.Sprintf("%s-%s-%s", w1, w2, w3)

		// Check for uniqueness
		if !generated[code] {
			generated[code] = true
			codes = append(codes, code)
		}
	}

	return codes, nil
}

func main() {
	// Parse command-line arguments
	count := defaultCount
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err != nil || parsed < 1 {
			fmt.Fprintf(os.Stderr, "Error: invalid count argument. Must be a positive integer.\n")
			os.Exit(1)
		}
		count = parsed
	}

	// Read words from dictionary
	words, err := readWords()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Generate promo codes
	codes, err := generatePromoCodes(words, count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create and run the TUI
	m := initialModel(codes)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
