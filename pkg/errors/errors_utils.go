package errors_utils

import "fmt"

var (
	ErrGenerateUUID     = fmt.Errorf("erro criar um uuid")
	ErrDatabaseInsert   = fmt.Errorf("erro ao inserir o dado na base")
	ErrSendMessageQueue = fmt.Errorf("erro ao enviar mensagem para fila")
	ErrMarshalEvent     = fmt.Errorf("erro ao realizar o marshal do evento sqs")
)
