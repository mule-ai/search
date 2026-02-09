package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/mule-ai/search/internal/searxng"
)

// JSONFormatter formats search results as JSON.
type JSONFormatter struct {
	Pretty bool // Enable pretty-printed output with indentation
}

// NewJSONFormatter creates a new JSON formatter with pretty-printing enabled.
//
// Example:
//
//	jf := formatter.NewJSONFormatter()
//	output, err := jf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		Pretty: true,
	}
}

// Format formats the search results as JSON.
//
// The output includes the query, result count, results array, and metadata.
// Answers, infoboxes, and suggestions are included if present.
// Returns a JSON string or an error if formatting fails.
//
// Example:
//
//	response := &searxng.SearchResponse{
//	    Query: "golang",
//	    Results: []searxng.SearchResult{...},
//	    NumberOfResults: 100,
//	}
//	jf := formatter.NewJSONFormatter()
//	output, err := jf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func (f *JSONFormatter) Format(result *searxng.SearchResponse) (string, error) {
	if result == nil {
		return "", fmt.Errorf("nil response provided")
	}

	// Create output structure matching SPEC
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", result.SearchTime),
		"instance":    result.Instance,
	}

	// Add pagination info if page > 1
	if result.Page > 1 {
		metadata["page"] = result.Page
	}

	output := map[string]interface{}{
		"query":         result.Query,
		"total_results": result.NumberOfResults,
		"results":       f.formatResults(result.Results),
		"metadata":      metadata,
	}

	if len(result.Answers) > 0 {
		output["answers"] = result.Answers
	}
	if len(result.Infoboxes) > 0 {
		output["infoboxes"] = result.Infoboxes
	}
	if len(result.Suggestions) > 0 {
		output["suggestions"] = result.Suggestions
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(data), nil
}

func (f *JSONFormatter) formatResults(results []searxng.SearchResult) []map[string]interface{} {
	var formatted []map[string]interface{}
	for _, result := range results {
		r := map[string]interface{}{
			"title":   result.Title,
			"url":     result.URL,
			"content": result.Content,
			"engine":  result.Engine,
			"category": result.Category,
			"score":   result.Score,
		}
		if result.ImgSrc != "" {
			r["img_src"] = result.ImgSrc
		}
		if len(result.ParsedURL) > 0 {
			r["parsed_url"] = result.ParsedURL
		}
		if result.Template != "" {
			r["template"] = result.Template
		}
		formatted = append(formatted, r)
	}
	return formatted
}

// FormatWithQuery formats results with a custom query string
func (f *JSONFormatter) FormatWithQuery(query string, results []searxng.SearchResult, searchTime float64, instance string) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
		"metadata": map[string]interface{}{
			"search_time": fmt.Sprintf("%.2fs", searchTime),
			"instance":    instance,
		},
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(data), nil
}

// DisablePretty disables pretty printing
func (f *JSONFormatter) DisablePretty() *JSONFormatter {
	f.Pretty = false
	return f
}

// EnablePretty enables pretty printing (default)
func (f *JSONFormatter) EnablePretty() *JSONFormatter {
	f.Pretty = true
	return f
}

// Is.pretty returns whether pretty printing is enabled
func (f *JSONFormatter) IsPretty() bool {
	return f.Pretty
}

// FormatAsArray formats results as a JSON array (useful for piping to jq)
func (f *JSONFormatter) FormatAsArray(results []searxng.SearchResult) (string, error) {
	var arr []map[string]interface{}
	for _, result := range results {
		r := map[string]interface{}{
			"title":   result.Title,
			"url":     result.URL,
			"content": result.Content,
			"engine":  result.Engine,
			"category": result.Category,
			"score":   result.Score,
		}
		if result.ImgSrc != "" {
			r["img_src"] = result.ImgSrc
		}
		arr = append(arr, r)
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(arr, "", "  ")
	} else {
		data, err = json.Marshal(arr)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON array: %w", err)
	}

	return string(data), nil
}

// FormatResult formats a single result as JSON object
func (f *JSONFormatter) FormatResult(result searxng.SearchResult) (string, error) {
	r := map[string]interface{}{
		"title":   result.Title,
		"url":     result.URL,
		"content": result.Content,
		"engine":  result.Engine,
		"category": result.Category,
		"score":   result.Score,
	}
	if result.ImgSrc != "" {
		r["img_src"] = result.ImgSrc
	}
	if len(result.ParsedURL) > 0 {
		r["parsed_url"] = result.ParsedURL
	}
	if result.Template != "" {
		r["template"] = result.Template
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(r, "", "  ")
	} else {
		data, err = json.Marshal(r)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal single result: %w", err)
	}

	return string(data), nil
}

// FormatWithMetadata formats results with custom metadata
func (f *JSONFormatter) FormatWithMetadata(query string, results []searxng.SearchResult, metadata map[string]interface{}) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
	}

	// Add custom metadata
	for k, v := range metadata {
		output[k] = v
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON with metadata: %w", err)
	}

	return string(data), nil
}

// FormatWithSearchTime formats results with a custom search time
func (f *JSONFormatter) FormatWithSearchTime(query string, results []searxng.SearchResult, searchTime float64) (string, error) {
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatWithInstance formats results with a custom instance URL
func (f *JSONFormatter) FormatWithInstance(query string, results []searxng.SearchResult, instance string) (string, error) {
	metadata := map[string]interface{}{
		"instance": instance,
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatCompact formats results in a compact JSON format
func (f *JSONFormatter) FormatCompact(query string, results []searxng.SearchResult) (string, error) {
	output := map[string]interface{}{
		"query": query,
	}

	var resultArr []map[string]interface{}
	for _, result := range results {
		r := map[string]interface{}{
			"title":    result.Title,
			"url":      result.URL,
			"content":  result.Content,
			"engine":   result.Engine,
			"category": result.Category,
			"score":    result.Score,
		}
		resultArr = append(resultArr, r)
	}
	output["results"] = resultArr

	data, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal compact JSON: %w", err)
	}

	return string(data), nil
}

// FormatWithTotal formats results with custom total count
func (f *JSONFormatter) FormatWithTotal(query string, results []searxng.SearchResult, total int) (string, error) {
	metadata := map[string]interface{}{
		"total_results": total,
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatWithAll formats results with all custom options
func (f *JSONFormatter) FormatWithAll(query string, results []searxng.SearchResult, searchTime float64, instance string, total int) (string, error) {
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
		"instance":    instance,
	}
	if total > 0 {
		metadata["total_results"] = total
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatResultsOnly formats only the results array
func (f *JSONFormatter) FormatResultsOnly(results []searxng.SearchResult) (string, error) {
	return f.FormatAsArray(results)
}

// FormatSimple formats a simple result list with just title and URL
func (f *JSONFormatter) FormatSimple(results []searxng.SearchResult) (string, error) {
	var arr []map[string]interface{}
	for _, result := range results {
		r := map[string]interface{}{
			"title": result.Title,
			"url":   result.URL,
		}
		arr = append(arr, r)
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(arr, "", "  ")
	} else {
		data, err = json.Marshal(arr)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal simple results: %w", err)
	}

	return string(data), nil
}

// FormatSearchResultWithMetadata formats a single result with metadata
func (f *JSONFormatter) FormatSearchResultWithMetadata(result searxng.SearchResult, index int, metadata map[string]interface{}) (string, error) {
	r := map[string]interface{}{
		"title":   result.Title,
		"url":     result.URL,
		"content": result.Content,
		"engine":  result.Engine,
		"category": result.Category,
		"score":   result.Score,
		"index":   index,
	}
	if result.ImgSrc != "" {
		r["img_src"] = result.ImgSrc
	}
	if len(result.ParsedURL) > 0 {
		r["parsed_url"] = result.ParsedURL
	}
	if result.Template != "" {
		r["template"] = result.Template
	}

	for k, v := range metadata {
		r[k] = v
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(r, "", "  ")
	} else {
		data, err = json.Marshal(r)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal result with metadata: %w", err)
	}

	return string(data), nil
}

// FormatError formats an error as JSON for scripting purposes
func (f *JSONFormatter) FormatError(message string, code int) string {
	output := map[string]interface{}{
		"error":   message,
		"code":    code,
	}
	data, _ := json.Marshal(output)
	return string(data)
}

// FormatSearchRequest formats a search request as JSON
func (f *JSONFormatter) FormatSearchRequest(query string, page int, results int) string {
	output := map[string]interface{}{
		"query":   query,
		"page":    page,
		"results": results,
	}
	data, _ := json.Marshal(output)
	return string(data)
}

// FormatSearchResponse formats the entire search response
func (f *JSONFormatter) FormatSearchResponse(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, numberOfResults int, searchTime float64) (string, error) {
	output := map[string]interface{}{
		"query":         query,
		"results":       f.formatResults(results),
		"answers":       answers,
		"infoboxes":     infoboxes,
		"suggestions":   suggestions,
		"total_results": numberOfResults,
	}
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal search response: %w", err)
	}

	return string(data), nil
}

// FormatWithInstanceAndTime formats with custom instance and search time
func (f *JSONFormatter) FormatWithInstanceAndTime(query string, results []searxng.SearchResult, instance string, searchTime float64) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
		"metadata": map[string]interface{}{
			"search_time": fmt.Sprintf("%.2fs", searchTime),
			"instance":    instance,
		},
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with instance and time: %w", err)
	}

	return string(data), nil
}

// FormatWithAnswers formats with answers section
func (f *JSONFormatter) FormatWithAnswers(query string, results []searxng.SearchResult, answers []searxng.Answer) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
		"answers":    answers,
	}
	metadata := map[string]interface{}{
		"search_time": "unknown",
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with answers: %w", err)
	}

	return string(data), nil
}

// FormatWithInfoboxes formats with infoboxes section
func (f *JSONFormatter) FormatWithInfoboxes(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
		"infoboxes":  infoboxes,
	}
	metadata := map[string]interface{}{
		"search_time": "unknown",
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with infoboxes: %w", err)
	}

	return string(data), nil
}

// FormatWithSuggestions formats with suggestions section
func (f *JSONFormatter) FormatWithSuggestions(query string, results []searxng.SearchResult, suggestions []string) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
		"suggestions": suggestions,
	}
	metadata := map[string]interface{}{
		"search_time": "unknown",
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with suggestions: %w", err)
	}

	return string(data), nil
}

// FormatWithAllMetadata formats with all metadata fields
func (f *JSONFormatter) FormatWithAllMetadata(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, numberOfResults int, searchTime float64) (string, error) {
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
	}
	if answers != nil {
		metadata["answers"] = answers
	}
	if infoboxes != nil {
		metadata["infoboxes"] = infoboxes
	}
	if suggestions != nil {
		metadata["suggestions"] = suggestions
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatWithAnswersAndInfoboxes formats with answers and infoboxes
func (f *JSONFormatter) FormatWithAnswersAndInfoboxes(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox) (string, error) {
	metadata := map[string]interface{}{
		"answers":   answers,
		"infoboxes": infoboxes,
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatWithAnswersAndSuggestions formats with answers and suggestions
func (f *JSONFormatter) FormatWithAnswersAndSuggestions(query string, results []searxng.SearchResult, answers []searxng.Answer, suggestions []string) (string, error) {
	metadata := map[string]interface{}{
		"answers":     answers,
		"suggestions": suggestions,
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatWithInfoboxesAndSuggestions formats with infoboxes and suggestions
func (f *JSONFormatter) FormatWithInfoboxesAndSuggestions(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox, suggestions []string) (string, error) {
	metadata := map[string]interface{}{
		"infoboxes":   infoboxes,
		"suggestions": suggestions,
	}
	return f.FormatWithMetadata(query, results, metadata)
}

// FormatFull formats the complete search response with all fields
func (f *JSONFormatter) FormatFull(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, numberOfResults int, searchTime float64) (string, error) {
	output := map[string]interface{}{
		"query":         query,
		"total_results": numberOfResults,
		"results":       f.formatResults(results),
		"answers":       answers,
		"infoboxes":     infoboxes,
		"suggestions":   suggestions,
	}
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal full response: %w", err)
	}

	return string(data), nil
}

// FormatWithFormattedResults formats with pre-formatted results
func (f *JSONFormatter) FormatWithFormattedResults(query string, formattedResults []map[string]interface{}, numberOfResults int, searchTime float64) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": numberOfResults,
		"results":    formattedResults,
		"metadata": map[string]interface{}{
			"search_time": fmt.Sprintf("%.2fs", searchTime),
		},
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with formatted results: %w", err)
	}

	return string(data), nil
}

// FormatWithCustomMetadata formats with custom metadata map
func (f *JSONFormatter) FormatWithCustomMetadata(query string, results []searxng.SearchResult, metadata map[string]interface{}) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
	}
	if metadata != nil {
		output["metadata"] = metadata
	}

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with custom metadata: %w", err)
	}

	return string(data), nil
}

// FormatWithQueryAndResults formats with custom query and results
func (f *JSONFormatter) FormatWithQueryAndResults(query string, results []searxng.SearchResult) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
	}
	metadata := map[string]interface{}{
		"search_time": "unknown",
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with query and results: %w", err)
	}

	return string(data), nil
}

// FormatWithQueryResultsAndTime formats with query, results, and search time
func (f *JSONFormatter) FormatWithQueryResultsAndTime(query string, results []searxng.SearchResult, searchTime float64) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
	}
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with query, results, and time: %w", err)
	}

	return string(data), nil
}

// FormatWithQueryResultsTimeAndInstance formats with all custom values
func (f *JSONFormatter) FormatWithQueryResultsTimeAndInstance(query string, results []searxng.SearchResult, searchTime float64, instance string) (string, error) {
	output := map[string]interface{}{
		"query":      query,
		"total_results": len(results),
		"results":    f.formatResults(results),
	}
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
		"instance":    instance,
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with query, results, time, and instance: %w", err)
	}

	return string(data), nil
}

// FormatWithAllOptions formats with all available options
func (f *JSONFormatter) FormatWithAllOptions(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, numberOfResults int, searchTime float64, instance string) (string, error) {
	output := map[string]interface{}{
		"query":         query,
		"total_results": numberOfResults,
		"results":       f.formatResults(results),
		"answers":       answers,
		"infoboxes":     infoboxes,
		"suggestions":   suggestions,
	}
	metadata := map[string]interface{}{
		"search_time": fmt.Sprintf("%.2fs", searchTime),
		"instance":    instance,
	}
	output["metadata"] = metadata

	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal with all options: %w", err)
	}

	return string(data), nil
}

// FormatResultsArray formats results as a simple array (for jq compatibility)
func (f *JSONFormatter) FormatResultsArray(results []searxng.SearchResult) string {
	var arr []map[string]interface{}
	for _, result := range results {
		r := map[string]interface{}{
			"title":   result.Title,
			"url":     result.URL,
			"content": result.Content,
			"engine":  result.Engine,
			"category": result.Category,
			"score":   result.Score,
		}
		arr = append(arr, r)
	}

	data, _ := json.MarshalIndent(arr, "", "  ")
	return string(data)
}

// FormatResponse formats a SearchResponse struct
func (f *JSONFormatter) FormatResponse(result *searxng.SearchResponse) string {
	s, _ := f.FormatWithAllMetadata(
		result.Query,
		result.Results,
		result.Answers,
		result.Infoboxes,
		result.Suggestions,
		result.NumberOfResults,
		result.SearchTime,
	)
	return s
}
