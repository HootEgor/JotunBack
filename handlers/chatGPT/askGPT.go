package chatGPT

import (
	"JotunBack/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func SendGPTRequest(userMessage string, acState *model.ACState) (model.GPTResponse, error) {
	prompt, err := ioutil.ReadFile("serviceAccount/prompt.txt")
	if err != nil {
		log.Fatal(err)
	}
	promptString := string(prompt) + time.Now().Format("2006-01-02 15:04")
	promptString += "Current temperature:" + acState.GetTargetTemp() + "\n"
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
		return model.GPTResponse{}, err
	}

	// Формируем HTTP-запрос
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.GPTResponse{}, err
	}

	gptKeyString := string(gptKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gptKeyString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.GPTResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.GPTResponse{}, err
	}

	response := model.GPTResponse{}
	for {
		response, err = parseGPTResponse(respBody)
		if err != nil {
			log.Println(err)
			continue
		}
		break
	}

	return response, nil
}

func parseGPTResponse(respBody []byte) (model.GPTResponse, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(respBody), &data)
	if err != nil {
		return model.GPTResponse{}, err
	}

	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return model.GPTResponse{}, errors.New("Choice not found or empty")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return model.GPTResponse{}, errors.New("Invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return model.GPTResponse{}, errors.New("Message not found or not a map")
	}

	content, ok := message["content"].(string)
	if !ok {
		return model.GPTResponse{}, errors.New("Content not found or not a string")
	}

	content = strings.ReplaceAll(content, "```json\n", "")
	content = strings.ReplaceAll(content, "```", "")

	var response model.GPTResponse
	err = json.Unmarshal([]byte(content), &response)
	if err != nil {
		return model.GPTResponse{}, err
	}

	return response, nil
}
