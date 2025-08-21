package userdb

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type TokenStore struct {
	tokens map[string]int
	mu     sync.RWMutex
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]int),
	}
}

func (ts *TokenStore) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open token file: %w", err)
	}
	defer file.Close()

	ts.mu.Lock()
	defer ts.mu.Unlock()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			token := parts[0]
			userID, err := strconv.Atoi(parts[1])
			if err == nil {
				ts.tokens[token] = userID
			}
		}
	}

	return scanner.Err()
}

func (ts *TokenStore) AddToken(token string, userID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tokens[token] = userID
}

func (ts *TokenStore) ValidateToken(token string) (int, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	
	userID, exists := ts.tokens[token]
	return userID, exists
}

func (ts *TokenStore) GetTokenCount() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.tokens)
}
