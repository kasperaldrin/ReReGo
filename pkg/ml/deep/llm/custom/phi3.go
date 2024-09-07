package ml

import "kasperaldrin.com/rerego/pkg/core"

// Phi3Node is a node that wraps Phi3 functionality into a net call. You must host it yourself, which can be done by using the LLM server. I highly recommend you to, at the moment, configure your own genai functionality using net nodes.

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Phi3GenerationArgs struct {
	MaxNewTokens   int     `json:"max_new_tokens"`
	ReturnFullText bool    `json:"return_full_text"`
	Temperature    float64 `json:"temperature"`
	DoSample       bool    `json:"do_sample"`
}

type Phi3Request struct {
	Messages           []Message          `json:"messages"`
	Phi3GenerationArgs Phi3GenerationArgs `json:"generation_args"`
}

type Phi3Node struct {
	core.NodeBase[Phi3Request, string]
}

/*
func NewPhi3Node(label string) Phi3Node {
	return Phi3Node{
		NodeBase: core.NewNodeBase[Phi3Request, string](label).
			UseFunctionWhenAny(func(
				input core.DataObject[Phi3Request],
			) core.DataObject[string] {
				// Marshal the request into JSON
				// Marshal the request into JSON
				jsonData, err := json.Marshal(requestBody)
				if err != nil {
					log.Fatalf("Error marshalling JSON: %v", err)
				}
				// Send the POST request
				resp, err := http.Post("http://0.0.0.0:8000/generate/", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Fatalf("Error making POST request: %v", err)
				}
				defer resp.Body.Close()

				// Read the response
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("Error reading response: %v", err)
				}
			},
		},
	}






func (n WebGetNode) UseFunction() {

} // Cannot replace default function.



// Print the response
fmt.Println("Response from server:", string(body))

*/
