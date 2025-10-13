package compress

import(
	"github.com/mholt/archiver/v3"
)
type Archive v3.Archive

type unsupported struct{
	source string
	err error
}

func (u unsupported)Unarchive(source, destination string)error{
	return u.err
}

func NewUnarchiver(source string)v3.Unarchiver{
	uaIface, err := v3.ByExtension(source)
	if err != nil {
		return unsupported{source:source,err:err}
	}
	u, ok := uaIface.(Unarchiver)
	if !ok {
		return unsupported{source:source,err: fmt.Errorf("format specified by source filename is not an archive format: %s (%T)", source, uaIface))
	}
	return u
}