package common

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const DefaultSystemInstruction = "You are Antigravity, a powerful agentic AI coding assistant designed by the Google Deepmind team working on Advanced Agentic Coding.You are pair programming with a USER to solve their coding task. The task may require creating a new codebase, modifying or debugging an existing codebase, or simply answering a question.**Absolute paths only****Proactiveness**"

// AttachDefaultSystemInstruction ensures the default system instruction is appended when absent.
func AttachDefaultSystemInstruction(rawJSON []byte) []byte {
	return appendSystemInstruction(rawJSON, DefaultSystemInstruction)
}

func appendSystemInstruction(rawJSON []byte, prompt string) []byte {
	paths := []string{"system_instruction", "request.systemInstruction"}
	for _, path := range paths {
		if gjson.GetBytes(rawJSON, path).Exists() {
			return appendSystemInstructionAtPath(rawJSON, path, prompt)
		}
	}

	return appendSystemInstructionAtPath(rawJSON, "system_instruction", prompt)
}

func appendSystemInstructionAtPath(rawJSON []byte, path, prompt string) []byte {
	partsPath := path + ".parts"
	if parts := gjson.GetBytes(rawJSON, partsPath); parts.Exists() && parts.IsArray() {
		for _, part := range parts.Array() {
			if part.Get("text").String() == prompt {
				return rawJSON
			}
		}
		rawJSON, _ = sjson.SetBytes(rawJSON, partsPath+".-1.text", prompt)
	} else {
		rawJSON, _ = sjson.SetBytes(rawJSON, partsPath+".0.text", prompt)
	}

	if !gjson.GetBytes(rawJSON, path+".role").Exists() {
		rawJSON, _ = sjson.SetBytes(rawJSON, path+".role", "user")
	}

	return rawJSON
}
