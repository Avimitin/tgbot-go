package main

var unwrapTemplate = `
<b>Message Information</b>
=== <b>CHAT</b> ===
ID: <code>{{.Chat.ID}}</code>
TYPE: <code>{{.Chat.Type}}</code>
USERNAME: <code>{{.Chat.Username}}</code>
=== <b>USER</b> ===
ID: <code>{{.Sender.ID}}</code>
USERNAME: <code>{{.Sender.Username}}</code>
NICKNAME: <code>{{.Sender.FirstName}} {{if .Sender.LastName}}{{.Sender.LastName}}{{end}}</code>
LANGUAGE: <code>{{.Sender.LanguageCode}}</code>
=== <b>MSG</b> ===
ID: <code>{{.ID}}</code>
{{if .ReplyTo -}}
=== <b> Reply INFO </b> ===
=== <b> Forward INFO </b> ===
{{if .ReplyTo.OriginalChat -}}
FROM_CHAT: <code>{{- .ReplyTo.OriginalChat.ID -}}</code>
TYPE: <code>{{- .ReplyTo.OriginalChat.Type -}}</code>
	{{end}}
	{{- if .ReplyTo.OriginalSenderName -}}
Sender: {{- .ReplyTo.OriginalSenderName -}}
	{{end}}
	{{- if .ReplyTo.OriginalSender -}}
=== <b> Forward User INFO </b> ===
FROM_USER: <code>{{- .ReplyTo.OriginalSender.ID -}}
USERNAME: <code>{{- .ReplyTo.OriginalSender.Username -}}</code>
NICKNAME: <code>{{- .ReplyTo.OriginalSender.FirstName -}}
		{{- if .ReplyTo.OriginalSender.LastName -}}
			{{- .ReplyTo.OriginalSender.LastName -}}
		{{end}}</code>
	{{end}}
{{- end}}
=== <b> File Information </b> ===
{{if .ReplyTo.Document -}}
FILE_ID: <code>{{- .ReplyTo.Document.FileID -}}</code>
{{end -}}
{{if .ReplyTo.Video -}}
File_ID: <code>{{- .ReplyTo.Video.FileID -}}</code>
{{end -}}
{{if .ReplyTo.Photo -}}
File_ID: <code>{{- .ReplyTo.Photo.FileID -}}</code>
{{end -}}
{{if .ReplyTo.Audio -}}
File_ID: <code>{{- .ReplyTo.Audio.FileID -}}</code>
{{end -}}
`
