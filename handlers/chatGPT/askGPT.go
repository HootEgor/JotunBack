package chatGPT

import (
	"JotunBack/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func SendGPTRequest(userMessage string) error {
	prompt, err := ioutil.ReadFile("serviceAccount/prompt.txt")
	if err != nil {
		log.Fatal(err)
	}
	promptString := string(prompt) + time.Now().Format("2006-01-02T15:04:05Z07:00")
	gptKey, err := ioutil.ReadFile("serviceAccount/gptKey.txt")
	if err != nil {
		log.Fatal(err)
	}

	requestData := map[string]interface{}{
		"model":      "gpt-3.5-turbo",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "system", "content": promptString},
			{"role": "user", "content": userMessage}},
	}

	// Кодируем данные запроса в JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	// Формируем HTTP-запрос
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	gptKeyString := string(gptKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gptKeyString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println("GPT response:", parseGPTResponse(respBody))

	return nil
}

func parseGPTResponse(respBody []byte) model.GPTResponse {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(respBody), &data)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем массив choices
	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Fatal("Choices not found or empty")
	}

	// Получаем первый элемент массива choices
	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Fatal("Invalid choice format")
	}

	// Получаем содержимое сообщения
	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		log.Fatal("Invalid message format")
	}

	// Получаем содержимое контента
	content, ok := message["content"].(string)
	if !ok {
		log.Fatal("Content not found or not a string")
	}

	content = strings.ReplaceAll(content, "```json\n", "")
	content = strings.ReplaceAll(content, "```", "")

	// Разбираем содержимое в структуру GPTResponse
	var response model.GPTResponse
	err = json.Unmarshal([]byte(content), &response)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
