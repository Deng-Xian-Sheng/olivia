package embeddings

import (
	"fmt"
	"sort"
	"strings"

	"github.com/olivia-ai/olivia/data"
	"github.com/olivia-ai/olivia/util"
	"github.com/tebeka/snowball"
)

type embeddingMapping struct {
	index int
	value float64
}

func tokenize(sentence string) (tokens []string) {
	tokens = strings.Fields(
		strings.ToLower(sentence),
	)

	for i, token := range tokens {
		tokens[i] = stem(token)
	}

	return
}

func stem(word string) string {
	stemmer, err := snowball.New("english")
	if err != nil {
		fmt.Printf("Unable to load stemmer. Returning un-stemmed word. %s\n", err)
		return word
	}

	return stemmer.Stem(word)
}

func appendToVocabulary(tokenBase *[]string, words ...string) {
	for _, word := range words {
		if !util.Contains(*tokenBase, word) {
			*tokenBase = append(*tokenBase, word)
		}
	}
}

func getEmbeddingWithIndex(vocabularySize, index int) []float64 {
	embedding := make([]float64, vocabularySize)
	embedding[index] = 1
	return embedding
}

func getEmbedding(vocabulary []string, word string) []float64 {
	embedding := []float64{
		// Initialize two values at the beginning for the start and end tokens
		0, 0,
	}

	for _, token := range vocabulary {
		var value float64 = 0
		if token == word {
			value = 1
		}

		embedding = append(embedding, value)
	}
	
	return embedding
}

// EstablishVocabulary takes a slice of Conversation structs and establish the vocabulary set.
func EstablishVocabulary(conversations []data.Conversation) (vocabulary []string) {
	for _, conversation := range conversations {
		// Iterate through the answer and question to avoid code duplication
		for _, sentence := range []string{conversation.Answer, conversation.Question} {
			appendToVocabulary(&vocabulary, tokenize(sentence)...)
		}
	}

	return
}

func ClosestEmbedding(embedding []float64) []float64 {
	// Initialize the embedding mapping the index with its value to prepare for the sort
	valuesMapping := []embeddingMapping{}
	for i, value := range embedding {
		valuesMapping = append(valuesMapping, embeddingMapping{i, value})
	}

	sort.Slice(valuesMapping, func(i, j int) bool {
		return valuesMapping[i].value > valuesMapping[j].value
	})

	result := make([]float64, len(embedding))
	result[valuesMapping[0].index] = 1
	return result
}

// ReconstructWord takes the vocabulary and an embedding and returns the natural language word
// matching the embedding.
func ReconstructWord(vocabulary []string, embedding []float64) string {
	for i, value := range embedding {
		if i == 0 || i == 1 {
			continue
		}

		if value == 1 {
			return vocabulary[i]
		}
	}

	return ""
}

// GetEOS returns the embedding for the end of sentence token.
func GetEOS(vocabularySize int) []float64 {
	return getEmbeddingWithIndex(vocabularySize, 1)
}

// GetEOS returns the embedding for the beginning of sentence token.
func GetBOS(vocabularySize int) []float64 {
	return getEmbeddingWithIndex(vocabularySize, 0)
}