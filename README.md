# groq

![GitHub License](https://img.shields.io/github/license/grqphical/groq)
![GitHub Tag](https://img.shields.io/github/v/tag/grqphical/groq?label=version)
[![Go Tests](https://github.com/grqphical/groq/actions/workflows/go.yml/badge.svg)](https://github.com/grqphical/groq/actions/workflows/go.yml)

An **UNOFFICIAL** API wrapper for [Groq's](https://groq.com) API for Golang

> ### NOTE:
>
> This is an unofficial project and is not affiliated with Groq in any way.

Currently vision APIs and Streaming results are unsupported (I plan to have this complete soon)

## Installation

Run:

```bash
go get github.com/grqphical/groq
```

In order to use this library you need a Groq API key. You can get one for _free_ from [here](https://console.groq.com/keys)

## Example

```go
package main

import "github.com/grqphical/groq"

const apiKey string = "API_KEY_HERE"

func main() {
    client, err := groq.NewClient(apiKey)
    if err != nil {
        panic(err)
    }

    conversation := groq.NewConversation("answer using markdown only if necessary")

    conversation.AddMessages(groq.Message{
        Role: groq.MessageRoleUser,
        Content: "How tall is the Eifel Tower?"
    })

    response, err := conversation.Complete(client, "llama-3-8b-8192", nil)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\n", response.Choices[0].Message.Content)
}
```

## License

`groq` is licensed under the MIT license
