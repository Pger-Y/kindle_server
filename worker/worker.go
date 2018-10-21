package worker

/*
import (
	"fmt"
)
*/

type Worker interface {
	Usage() string
	//Info(params ...string) (string, error)
	//Register(params ...string) (string, error)
	//Send(params ...string) (string, error)
	//Feedback(params ...string) (string, error)
	//Search(params ...string) (string, error)
	Exec(userid, msg_type string, params ...string) (string, error)
}
